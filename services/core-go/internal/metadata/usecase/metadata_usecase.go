package usecase

import (
	"context"
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
	Update(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) error
	Delete(ctx context.Context, user *permDomain.UserContext, id string) error
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

func (u *metadataUsecase) Update(ctx context.Context, user *permDomain.UserContext, meta *domain.Metadata) error {
	existing, err := u.repo.GetByID(ctx, meta.ID)
	if err != nil {
		return err
	}

	if meta.Title == "" {
		return domain.ErrTitleRequired
	}
	if meta.FolderID == "" {
		return domain.ErrFolderRequired
	}

	// Authorize metadata write permission (FR-32)
	err = u.permEngine.AuthorizeMetadata(ctx, user, meta.ID, permDomain.ActionWrite)
	if err != nil {
		return err
	}

	// If moving to another folder, authorize write permission on the destination folder
	if existing.FolderID != meta.FolderID {
		_, err = u.folderRepo.GetByID(ctx, meta.FolderID)
		if err != nil {
			return err
		}

		err = u.permEngine.AuthorizeFolder(ctx, user, meta.FolderID, permDomain.ActionWrite)
		if err != nil {
			return err
		}
	}

	return u.repo.Update(ctx, meta)
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
