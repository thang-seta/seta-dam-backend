package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type mockPermRepo struct {
	checkFolder   func(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error)
	checkMetadata func(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error)
	addPerm       func(ctx context.Context, perm *domain.ObjectPermission, grantedByID string) error
	deletePerm    func(ctx context.Context, userID string, objectType string, objectID string, action string) error
	getEffective  func(ctx context.Context, userID string, objectType string, objectID string) ([]string, error)
}

func (m *mockPermRepo) CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error) {
	return m.checkFolder(ctx, userID, folderID, action)
}

func (m *mockPermRepo) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error) {
	return m.checkMetadata(ctx, userID, metadataID, action)
}

func (m *mockPermRepo) AddPermission(ctx context.Context, perm *domain.ObjectPermission, grantedByID string) error {
	return m.addPerm(ctx, perm, grantedByID)
}

func (m *mockPermRepo) DeletePermission(ctx context.Context, userID string, objectType string, objectID string, action string) error {
	return m.deletePerm(ctx, userID, objectType, objectID, action)
}

func (m *mockPermRepo) ListPermissions(ctx context.Context) ([]*domain.ObjectPermission, error) {
	return nil, nil
}

func (m *mockPermRepo) GetEffectivePermissions(ctx context.Context, userID string, objectType string, objectID string) ([]string, error) {
	return m.getEffective(ctx, userID, objectType, objectID)
}

func TestAuthorizeFolder(t *testing.T) {
	repo := &mockPermRepo{}
	evaluator := NewPermissionEvaluator(repo)

	t.Run("trainer_admin role bypasses everything", func(t *testing.T) {
		userCtx := &domain.UserContext{
			UserID: "user-1",
			Role:   "trainer_admin",
		}
		err := evaluator.AuthorizeFolder(context.Background(), userCtx, "folder-1", domain.ActionWrite)
		if err != nil {
			t.Fatalf("expected nil error for trainer_admin, got: %v", err)
		}
	})

	t.Run("user has permission", func(t *testing.T) {
		userCtx := &domain.UserContext{
			UserID: "user-2",
			Role:   "editor",
		}
		repo.checkFolder = func(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error) {
			if userID == "user-2" && folderID == "folder-1" && action == domain.ActionWrite {
				return true, nil
			}
			return false, nil
		}

		err := evaluator.AuthorizeFolder(context.Background(), userCtx, "folder-1", domain.ActionWrite)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("user does not have permission", func(t *testing.T) {
		userCtx := &domain.UserContext{
			UserID: "user-3",
			Role:   "viewer",
		}
		repo.checkFolder = func(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error) {
			return false, nil
		}

		err := evaluator.AuthorizeFolder(context.Background(), userCtx, "folder-1", domain.ActionWrite)
		if !errors.Is(err, domain.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})
}
