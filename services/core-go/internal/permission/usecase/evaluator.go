package usecase

import (
	"context"

	"github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/permission/repository"
)

type PermissionEvaluator interface {
	AuthorizeFolder(ctx context.Context, user *domain.UserContext, folderID string, action domain.Action) error
	AuthorizeMetadata(ctx context.Context, user *domain.UserContext, metadataID string, action domain.Action) error
	CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error)
	CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error)
	AddPermission(ctx context.Context, user *domain.UserContext, perm *domain.ObjectPermission) error
	DeletePermission(ctx context.Context, user *domain.UserContext, userID string, objectType string, objectID string, action string) error
	GetEffectivePermissions(ctx context.Context, user *domain.UserContext, targetUserID string, objectType string, objectID string) ([]string, error)
}

type permissionEvaluator struct {
	repo repository.PermissionRepository
}

func NewPermissionEvaluator(repo repository.PermissionRepository) PermissionEvaluator {
	return &permissionEvaluator{repo: repo}
}

func (pe *permissionEvaluator) AuthorizeFolder(ctx context.Context, user *domain.UserContext, folderID string, action domain.Action) error {
	if user == nil {
		return domain.ErrUnauthorized
	}
	if user.Role == "trainer_admin" {
		return nil
	}
	allowed, err := pe.repo.CheckFolderPermission(ctx, user.UserID, folderID, action)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrForbidden
	}
	return nil
}

func (pe *permissionEvaluator) AuthorizeMetadata(ctx context.Context, user *domain.UserContext, metadataID string, action domain.Action) error {
	if user == nil {
		return domain.ErrUnauthorized
	}
	if user.Role == "trainer_admin" {
		return nil
	}
	allowed, err := pe.repo.CheckMetadataPermission(ctx, user.UserID, metadataID, action)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrForbidden
	}
	return nil
}

func (pe *permissionEvaluator) CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error) {
	return pe.repo.CheckFolderPermission(ctx, userID, folderID, action)
}

func (pe *permissionEvaluator) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error) {
	return pe.repo.CheckMetadataPermission(ctx, userID, metadataID, action)
}

func (pe *permissionEvaluator) AddPermission(ctx context.Context, user *domain.UserContext, perm *domain.ObjectPermission) error {
	if user == nil {
		return domain.ErrUnauthorized
	}
	objectType, err := normalizeObjectType(perm.ObjectType)
	if err != nil {
		return err
	}
	if err := validateAction(perm.Action); err != nil {
		return err
	}
	perm.ObjectType = objectType

	// Only trainer_admin or users with manage_permissions on the target folder/metadata can grant permissions
	if user.Role != "trainer_admin" {
		if perm.ObjectType == "folder" {
			err = pe.AuthorizeFolder(ctx, user, perm.ObjectID, domain.ActionManagePermissions)
		} else {
			err = pe.AuthorizeMetadata(ctx, user, perm.ObjectID, domain.ActionManagePermissions)
		}
		if err != nil {
			return domain.ErrForbidden
		}
	}
	return pe.repo.AddPermission(ctx, perm, user.UserID)
}

func (pe *permissionEvaluator) DeletePermission(ctx context.Context, user *domain.UserContext, userID string, objectType string, objectID string, action string) error {
	if user == nil {
		return domain.ErrUnauthorized
	}
	objectType, err := normalizeObjectType(objectType)
	if err != nil {
		return err
	}
	if err := validateAction(action); err != nil {
		return err
	}

	// Only trainer_admin or users with manage_permissions can revoke permissions
	if user.Role != "trainer_admin" {
		if objectType == "folder" {
			err = pe.AuthorizeFolder(ctx, user, objectID, domain.ActionManagePermissions)
		} else {
			err = pe.AuthorizeMetadata(ctx, user, objectID, domain.ActionManagePermissions)
		}
		if err != nil {
			return domain.ErrForbidden
		}
	}
	return pe.repo.DeletePermission(ctx, userID, objectType, objectID, action)
}

func (pe *permissionEvaluator) GetEffectivePermissions(ctx context.Context, user *domain.UserContext, targetUserID string, objectType string, objectID string) ([]string, error) {
	// For demo/debugging purposes, we allow anyone to query effective permissions.
	// We delegate the evaluation to the repository layer, which checks both explicit permission rules and trainer_admin status.
	objectType, err := normalizeObjectType(objectType)
	if err != nil {
		return nil, err
	}
	return pe.repo.GetEffectivePermissions(ctx, targetUserID, objectType, objectID)
}

func normalizeObjectType(objectType string) (string, error) {
	switch objectType {
	case "folder":
		return "folder", nil
	case "metadata", "metadata_item":
		return "metadata", nil
	default:
		return "", domain.ErrInvalidObjectType
	}
}

func validateAction(action string) error {
	switch domain.Action(action) {
	case domain.ActionRead, domain.ActionWrite, domain.ActionManagePermissions:
		return nil
	default:
		return domain.ErrInvalidAction
	}
}
