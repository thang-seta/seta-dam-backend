package usecase

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	folderRepo "github.com/user/seta-dam-backend/services/core-go/internal/folder/repository"
	"github.com/user/seta-dam-backend/services/core-go/internal/metadata/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/metadata/repository"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
	permUsecase "github.com/user/seta-dam-backend/services/core-go/internal/permission/usecase"
)

type MetadataUsecase interface {
	Create(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) (*domain.Metadata, error)
	GetByID(ctx context.Context, user *permDomain.UserContext, id string) (*domain.Metadata, error)
	Update(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) (*domain.Metadata, error)
	Delete(ctx context.Context, user *permDomain.UserContext, id string) error
	List(ctx context.Context, user *permDomain.UserContext) ([]*domain.Metadata, error)
	ListByFolder(ctx context.Context, user *permDomain.UserContext, folderID string) ([]*domain.Metadata, error)
}

type metadataUsecase struct {
	repo       repository.MetadataRepository
	folderRepo folderRepo.FolderRepository
	permEngine permUsecase.PermissionEvaluator
}

func NewMetadataUsecase(repo repository.MetadataRepository, folderRepo folderRepo.FolderRepository, permEngine permUsecase.PermissionEvaluator) MetadataUsecase {
	return &metadataUsecase{
		repo:       repo,
		folderRepo: folderRepo,
		permEngine: permEngine,
	}
}

func (u *metadataUsecase) Create(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) (*domain.Metadata, error) {
	if user == nil {
		return nil, permDomain.ErrUnauthorized
	}
	normalizeMetadata(meta)
	if err := validateMetadata(meta); err != nil {
		return nil, err
	}
	if err := validateMetadataJSON(meta.MetadataJSON); err != nil {
		return nil, err
	}

	if meta.Title == "" {
		return nil, domain.ErrTitleRequired
	}
	if meta.FolderID == "" {
		return nil, domain.ErrFolderRequired
	}

	// 1. Verify parent folder exists
	_, err := u.folderRepo.GetByID(ctx, meta.FolderID)
	if err != nil {
		return nil, err
	}

	// 2. Authorize folder write permission (FR-32)
	err = u.permEngine.AuthorizeFolder(ctx, user, meta.FolderID, permDomain.ActionWrite)
	if err != nil {
		return nil, err
	}

	meta.ID = uuid.New().String()
	meta.CreatedBy = user.UserID
	err = u.repo.Create(ctx, meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (u *metadataUsecase) GetByID(ctx context.Context, user *permDomain.UserContext, id string) (*domain.Metadata, error) {
	meta, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Authorize metadata read permission (FR-32)
	err = u.permEngine.AuthorizeMetadata(ctx, user, id, permDomain.ActionRead)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (u *metadataUsecase) Update(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) (*domain.Metadata, error) {
	existing, err := u.repo.GetByID(ctx, meta.ID)
	if err != nil {
		return nil, err
	}

	normalizeMetadata(meta)
	if err := validateMetadata(meta); err != nil {
		return nil, err
	}
	if err := validateMetadataJSON(meta.MetadataJSON); err != nil {
		return nil, err
	}

	// Authorize metadata write permission (FR-32)
	err = u.permEngine.AuthorizeMetadata(ctx, user, meta.ID, permDomain.ActionWrite)
	if err != nil {
		return nil, err
	}

	// If moving to another folder, authorize write permission on the destination folder
	if existing.FolderID != meta.FolderID {
		_, err = u.folderRepo.GetByID(ctx, meta.FolderID)
		if err != nil {
			return nil, err
		}

		err = u.permEngine.AuthorizeFolder(ctx, user, meta.FolderID, permDomain.ActionWrite)
		if err != nil {
			return nil, err
		}
	}

	meta.UpdatedBy = &user.UserID
	if err := u.repo.Update(ctx, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func (u *metadataUsecase) Delete(ctx context.Context, user *permDomain.UserContext, id string) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Authorize metadata write permission (FR-32)
	err = u.permEngine.AuthorizeMetadata(ctx, user, id, permDomain.ActionWrite)
	if err != nil {
		return err
	}

	return u.repo.Delete(ctx, id)
}

func (u *metadataUsecase) List(ctx context.Context, user *permDomain.UserContext) ([]*domain.Metadata, error) {
	if user == nil {
		return nil, permDomain.ErrUnauthorized
	}

	allMetadata, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if user.Role == "trainer_admin" {
		return allMetadata, nil
	}

	visibleMetadata := make([]*domain.Metadata, 0, len(allMetadata))
	for _, meta := range allMetadata {
		allowed, err := u.permEngine.CheckMetadataPermission(ctx, user.UserID, meta.ID, permDomain.ActionRead)
		if err == nil && allowed {
			visibleMetadata = append(visibleMetadata, meta)
		}
	}
	return visibleMetadata, nil
}

func (u *metadataUsecase) ListByFolder(ctx context.Context, user *permDomain.UserContext, folderID string) ([]*domain.Metadata, error) {
	// Validate folder exists
	_, err := u.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// Authorize folder read permission to view metadata within it (FR-32)
	err = u.permEngine.AuthorizeFolder(ctx, user, folderID, permDomain.ActionRead)
	if err != nil {
		return nil, err
	}

	return u.repo.ListByFolder(ctx, folderID)
}

func normalizeMetadata(meta *domain.Metadata) {
	meta.Title = strings.TrimSpace(meta.Title)
	meta.FolderID = strings.TrimSpace(meta.FolderID)
	meta.MetadataJSON = strings.TrimSpace(meta.MetadataJSON)
}

func validateMetadata(meta *domain.Metadata) error {
	if meta.Title == "" {
		return domain.ErrTitleRequired
	}
	if meta.FolderID == "" {
		return domain.ErrFolderRequired
	}
	return nil
}

func validateMetadataJSON(raw string) error {
	if raw == "" {
		return nil
	}
	if !json.Valid([]byte(raw)) {
		return domain.ErrInvalidMetadataJSON
	}
	return nil
}
