package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type mockFolderRepo struct {
	create         func(ctx context.Context, folder *domain.Folder) error
	getByID        func(ctx context.Context, id string) (*domain.Folder, error)
	update         func(ctx context.Context, folder *domain.Folder) error
	delete         func(ctx context.Context, id string) error
	list           func(ctx context.Context) ([]*domain.Folder, error)
	isDescendantOf func(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error)
	isEmpty        func(ctx context.Context, id string) (bool, error)
}

func (m *mockFolderRepo) Create(ctx context.Context, folder *domain.Folder) error {
	return m.create(ctx, folder)
}

func (m *mockFolderRepo) GetByID(ctx context.Context, id string) (*domain.Folder, error) {
	return m.getByID(ctx, id)
}

func (m *mockFolderRepo) Update(ctx context.Context, folder *domain.Folder) error {
	return m.update(ctx, folder)
}

func (m *mockFolderRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func (m *mockFolderRepo) List(ctx context.Context) ([]*domain.Folder, error) {
	return m.list(ctx)
}

func (m *mockFolderRepo) IsDescendantOf(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error) {
	return m.isDescendantOf(ctx, potentialDescendantID, ancestorID)
}

func (m *mockFolderRepo) IsEmpty(ctx context.Context, id string) (bool, error) {
	return m.isEmpty(ctx, id)
}

type mockPermEvaluator struct {
	authFolder   func(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error
	authMetadata func(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error
	checkFolder  func(ctx context.Context, userID string, folderID string, action permDomain.Action) (bool, error)
}

func (m *mockPermEvaluator) AuthorizeFolder(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error {
	return m.authFolder(ctx, user, folderID, action)
}

func (m *mockPermEvaluator) AuthorizeMetadata(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error {
	return m.authMetadata(ctx, user, metadataID, action)
}

func (m *mockPermEvaluator) CheckFolderPermission(ctx context.Context, userID string, folderID string, action permDomain.Action) (bool, error) {
	return m.checkFolder(ctx, userID, folderID, action)
}

func (m *mockPermEvaluator) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action permDomain.Action) (bool, error) {
	return false, nil
}

func (m *mockPermEvaluator) AddPermission(ctx context.Context, user *permDomain.UserContext, perm *permDomain.ObjectPermission) error {
	return nil
}

func (m *mockPermEvaluator) DeletePermission(ctx context.Context, user *permDomain.UserContext, userID string, objectType string, objectID string, action string) error {
	return nil
}

func (m *mockPermEvaluator) GetEffectivePermissions(ctx context.Context, user *permDomain.UserContext, targetUserID string, objectType string, objectID string) ([]string, error) {
	return nil, nil
}

func TestMoveFolder(t *testing.T) {
	repo := &mockFolderRepo{}
	perm := &mockPermEvaluator{}
	uc := NewFolderUsecase(repo, perm)

	userCtx := &permDomain.UserContext{
		UserID: "editor-user",
		Role:   "editor",
	}

	repo.getByID = func(ctx context.Context, id string) (*domain.Folder, error) {
		return &domain.Folder{ID: id, Name: "folder"}, nil
	}

	perm.authFolder = func(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error {
		return nil
	}

	t.Run("fails when cycle is detected", func(t *testing.T) {
		repo.isDescendantOf = func(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error) {
			return true, nil
		}

		parentID := "child-folder"
		err := uc.MoveFolder(context.Background(), userCtx, "parent-folder", &parentID)
		if !errors.Is(err, domain.ErrCycleDetected) {
			t.Fatalf("expected ErrCycleDetected, got: %v", err)
		}
	})

	t.Run("succeeds when no cycle", func(t *testing.T) {
		repo.isDescendantOf = func(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error) {
			return false, nil
		}
		repo.update = func(ctx context.Context, folder *domain.Folder) error {
			return nil
		}

		parentID := "other-folder"
		err := uc.MoveFolder(context.Background(), userCtx, "folder-to-move", &parentID)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})
}
