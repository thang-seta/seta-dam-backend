package usecase

import (
	"context"
	"errors"
	"testing"

	folderDomain "github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
	metadataDomain "github.com/user/seta-dam-backend/services/core-go/internal/metadata/domain"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type mockMetadataRepo struct {
	create       func(ctx context.Context, meta *metadataDomain.Metadata) error
	getByID      func(ctx context.Context, id string) (*metadataDomain.Metadata, error)
	update       func(ctx context.Context, meta *metadataDomain.Metadata) error
	delete       func(ctx context.Context, id string) error
	list         func(ctx context.Context) ([]*metadataDomain.Metadata, error)
	listByFolder func(ctx context.Context, folderID string) ([]*metadataDomain.Metadata, error)
}

func (m *mockMetadataRepo) Create(ctx context.Context, meta *metadataDomain.Metadata) error {
	return m.create(ctx, meta)
}

func (m *mockMetadataRepo) GetByID(ctx context.Context, id string) (*metadataDomain.Metadata, error) {
	return m.getByID(ctx, id)
}

func (m *mockMetadataRepo) Update(ctx context.Context, meta *metadataDomain.Metadata) error {
	return m.update(ctx, meta)
}

func (m *mockMetadataRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func (m *mockMetadataRepo) List(ctx context.Context) ([]*metadataDomain.Metadata, error) {
	return m.list(ctx)
}

func (m *mockMetadataRepo) ListByFolder(ctx context.Context, folderID string) ([]*metadataDomain.Metadata, error) {
	return m.listByFolder(ctx, folderID)
}

type mockFolderRepo struct {
	getByID func(ctx context.Context, id string) (*folderDomain.Folder, error)
}

func (m *mockFolderRepo) Create(ctx context.Context, folder *folderDomain.Folder) error {
	return nil
}

func (m *mockFolderRepo) GetByID(ctx context.Context, id string) (*folderDomain.Folder, error) {
	return m.getByID(ctx, id)
}

func (m *mockFolderRepo) Update(ctx context.Context, folder *folderDomain.Folder) error {
	return nil
}

func (m *mockFolderRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockFolderRepo) List(ctx context.Context) ([]*folderDomain.Folder, error) {
	return nil, nil
}

func (m *mockFolderRepo) IsDescendantOf(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error) {
	return false, nil
}

func (m *mockFolderRepo) IsEmpty(ctx context.Context, id string) (bool, error) {
	return true, nil
}

type mockPermEvaluator struct {
	authFolder    func(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error
	authMetadata  func(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error
	checkMetadata func(ctx context.Context, userID string, metadataID string, action permDomain.Action) (bool, error)
}

func (m *mockPermEvaluator) AuthorizeFolder(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error {
	return m.authFolder(ctx, user, folderID, action)
}

func (m *mockPermEvaluator) AuthorizeMetadata(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error {
	return m.authMetadata(ctx, user, metadataID, action)
}

func (m *mockPermEvaluator) CheckFolderPermission(ctx context.Context, userID string, folderID string, action permDomain.Action) (bool, error) {
	return false, nil
}

func (m *mockPermEvaluator) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action permDomain.Action) (bool, error) {
	return m.checkMetadata(ctx, userID, metadataID, action)
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

func TestMetadataCreate(t *testing.T) {
	userCtx := &permDomain.UserContext{UserID: "editor-user", Role: "editor"}
	repo := &mockMetadataRepo{}
	folders := &mockFolderRepo{}
	perm := &mockPermEvaluator{}
	uc := NewMetadataUsecase(repo, folders, perm)

	folders.getByID = func(ctx context.Context, id string) (*folderDomain.Folder, error) {
		return &folderDomain.Folder{ID: id, Name: "folder"}, nil
	}
	perm.authFolder = func(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error {
		if action != permDomain.ActionWrite {
			t.Fatalf("expected folder write authorization, got %s", action)
		}
		return nil
	}
	repo.create = func(ctx context.Context, meta *metadataDomain.Metadata) error {
		if meta.ID == "" {
			t.Fatal("expected generated metadata id")
		}
		if meta.CreatedBy != userCtx.UserID {
			t.Fatalf("expected created_by %s, got %s", userCtx.UserID, meta.CreatedBy)
		}
		return nil
	}

	meta, err := uc.Create(context.Background(), userCtx, &metadataDomain.Metadata{
		FolderID:     "folder-1",
		Title:        "  Dog image  ",
		MetadataJSON: `{"source":"test"}`,
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if meta.Title != "Dog image" {
		t.Fatalf("expected trimmed title, got %q", meta.Title)
	}
}

func TestMetadataCreateValidation(t *testing.T) {
	uc := NewMetadataUsecase(&mockMetadataRepo{}, &mockFolderRepo{}, &mockPermEvaluator{})
	userCtx := &permDomain.UserContext{UserID: "editor-user", Role: "editor"}

	tests := []struct {
		name string
		meta *metadataDomain.Metadata
		want error
	}{
		{name: "title required", meta: &metadataDomain.Metadata{FolderID: "folder-1"}, want: metadataDomain.ErrTitleRequired},
		{name: "folder required", meta: &metadataDomain.Metadata{Title: "Image"}, want: metadataDomain.ErrFolderRequired},
		{name: "metadata json invalid", meta: &metadataDomain.Metadata{FolderID: "folder-1", Title: "Image", MetadataJSON: "{"}, want: metadataDomain.ErrInvalidMetadataJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Create(context.Background(), userCtx, tt.meta)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got: %v", tt.want, err)
			}
		})
	}
}

func TestMetadataUpdate(t *testing.T) {
	userCtx := &permDomain.UserContext{UserID: "editor-user", Role: "editor"}
	repo := &mockMetadataRepo{}
	folders := &mockFolderRepo{}
	perm := &mockPermEvaluator{}
	uc := NewMetadataUsecase(repo, folders, perm)

	repo.getByID = func(ctx context.Context, id string) (*metadataDomain.Metadata, error) {
		return &metadataDomain.Metadata{ID: id, FolderID: "source-folder", Title: "Before"}, nil
	}
	folders.getByID = func(ctx context.Context, id string) (*folderDomain.Folder, error) {
		return &folderDomain.Folder{ID: id, Name: "folder"}, nil
	}
	perm.authMetadata = func(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error {
		if action != permDomain.ActionWrite {
			t.Fatalf("expected metadata write authorization, got %s", action)
		}
		return nil
	}
	perm.authFolder = func(ctx context.Context, user *permDomain.UserContext, folderID string, action permDomain.Action) error {
		if folderID != "destination-folder" || action != permDomain.ActionWrite {
			t.Fatalf("expected destination folder write authorization, got %s %s", folderID, action)
		}
		return nil
	}
	repo.update = func(ctx context.Context, meta *metadataDomain.Metadata) error {
		return nil
	}

	meta, err := uc.Update(context.Background(), userCtx, &metadataDomain.Metadata{
		ID:       "meta-1",
		FolderID: "destination-folder",
		Title:    "After",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if meta.UpdatedBy == nil || *meta.UpdatedBy != userCtx.UserID {
		t.Fatalf("expected updated_by %s, got %v", userCtx.UserID, meta.UpdatedBy)
	}
}

func TestMetadataDeleteRequiresWritePermission(t *testing.T) {
	repo := &mockMetadataRepo{}
	perm := &mockPermEvaluator{}
	uc := NewMetadataUsecase(repo, &mockFolderRepo{}, perm)
	userCtx := &permDomain.UserContext{UserID: "viewer-user", Role: "viewer"}

	repo.getByID = func(ctx context.Context, id string) (*metadataDomain.Metadata, error) {
		return &metadataDomain.Metadata{ID: id, FolderID: "folder-1", Title: "Image"}, nil
	}
	perm.authMetadata = func(ctx context.Context, user *permDomain.UserContext, metadataID string, action permDomain.Action) error {
		return permDomain.ErrForbidden
	}
	repo.delete = func(ctx context.Context, id string) error {
		t.Fatal("delete should not be called when permission is denied")
		return nil
	}

	err := uc.Delete(context.Background(), userCtx, "meta-1")
	if !errors.Is(err, permDomain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got: %v", err)
	}
}
