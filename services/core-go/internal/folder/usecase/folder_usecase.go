package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/folder/repository"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
	permUsecase "github.com/user/seta-dam-backend/services/core-go/internal/permission/usecase"
)

type FolderUsecase interface {
	Create(ctx context.Context, user *permDomain.UserContext, name string, description *string, parentID *string) (*domain.Folder, error)
	GetByID(ctx context.Context, user *permDomain.UserContext, id string) (*domain.Folder, error)
	Update(ctx context.Context, user *permDomain.UserContext, id string, name string, description *string) error
	MoveFolder(ctx context.Context, user *permDomain.UserContext, id string, parentID *string) error
	Delete(ctx context.Context, user *permDomain.UserContext, id string) error
	ListTree(ctx context.Context, user *permDomain.UserContext) ([]*domain.Folder, error)
}

type folderUsecase struct {
	repo       repository.FolderRepository
	permEngine permUsecase.PermissionEvaluator
}

func NewFolderUsecase(repo repository.FolderRepository, permEngine permUsecase.PermissionEvaluator) FolderUsecase {
	return &folderUsecase{
		repo:       repo,
		permEngine: permEngine,
	}
}

func (u *folderUsecase) Create(ctx context.Context, user *permDomain.UserContext, name string, description *string, parentID *string) (*domain.Folder, error) {
	if user == nil {
		return nil, permDomain.ErrUnauthorized
	}

	// 1. Authorize creation
	if parentID == nil {
		// Creating root folder: only trainer_admin or editor
		if user.Role != "trainer_admin" && user.Role != "editor" {
			return nil, permDomain.ErrForbidden
		}
	} else {
		// Validate parent exists
		_, err := u.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, err // Will return ErrFolderNotFound if parent doesn't exist (FR-03)
		}
		// Creating child folder: needs write permission on the parent folder (FR-31)
		err = u.permEngine.AuthorizeFolder(ctx, user, *parentID, permDomain.ActionWrite)
		if err != nil {
			return nil, err
		}
	}

	folder := &domain.Folder{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		ParentID:    parentID,
		CreatedBy:   user.UserID,
	}

	err := u.repo.Create(ctx, folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (u *folderUsecase) GetByID(ctx context.Context, user *permDomain.UserContext, id string) (*domain.Folder, error) {
	// 1. Fetch folder first
	folder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Authorize read (FR-31)
	err = u.permEngine.AuthorizeFolder(ctx, user, id, permDomain.ActionRead)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (u *folderUsecase) Update(ctx context.Context, user *permDomain.UserContext, id string, name string, description *string) error {
	folder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Authorize write (FR-31)
	err = u.permEngine.AuthorizeFolder(ctx, user, id, permDomain.ActionWrite)
	if err != nil {
		return err
	}

	folder.Name = name
	folder.Description = description
	folder.UpdatedBy = &user.UserID
	return u.repo.Update(ctx, folder)
}

func (u *folderUsecase) MoveFolder(ctx context.Context, user *permDomain.UserContext, id string, parentID *string) error {
	folder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Prevent moving folder into itself
	if parentID != nil && *parentID == id {
		return domain.ErrCycleDetected
	}

	// Check write permission on current folder (FR-31)
	err = u.permEngine.AuthorizeFolder(ctx, user, id, permDomain.ActionWrite)
	if err != nil {
		return err
	}

	// Check write permission on destination folder (FR-31)
	if parentID != nil {
		// Validate new parent exists
		_, err = u.repo.GetByID(ctx, *parentID)
		if err != nil {
			return err
		}

		err = u.permEngine.AuthorizeFolder(ctx, user, *parentID, permDomain.ActionWrite)
		if err != nil {
			return err
		}

		// Prevent cycle (moving folder into one of its descendants)
		isDescendant, err := u.repo.IsDescendantOf(ctx, *parentID, id)
		if err != nil {
			return err
		}
		if isDescendant {
			return domain.ErrCycleDetected
		}
	}

	folder.ParentID = parentID
	folder.UpdatedBy = &user.UserID
	return u.repo.Update(ctx, folder)
}

func (u *folderUsecase) Delete(ctx context.Context, user *permDomain.UserContext, id string) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Authorize write/delete (FR-31)
	err = u.permEngine.AuthorizeFolder(ctx, user, id, permDomain.ActionWrite)
	if err != nil {
		return err
	}

	// Check if folder is empty (FR-10 / FR-09)
	isEmpty, err := u.repo.IsEmpty(ctx, id)
	if err != nil {
		return err
	}
	if !isEmpty {
		return domain.ErrNotEmpty
	}

	return u.repo.Delete(ctx, id)
}

func (u *folderUsecase) ListTree(ctx context.Context, user *permDomain.UserContext) ([]*domain.Folder, error) {
	if user == nil {
		return nil, permDomain.ErrUnauthorized
	}

	allFolders, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Filter folders by read permission
	var visibleFolders []*domain.Folder
	for _, folder := range allFolders {
		allowed, err := u.permEngine.CheckFolderPermission(ctx, user.UserID, folder.ID, permDomain.ActionRead)
		if err == nil && (user.Role == "trainer_admin" || allowed) {
			visibleFolders = append(visibleFolders, folder)
		}
	}
	return visibleFolders, nil
}
