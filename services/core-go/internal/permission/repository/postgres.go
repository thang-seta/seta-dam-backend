package repository

import (
	"context"
	"database/sql"
	"github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type PermissionRepository interface {
	CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error)
	CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error)
	AddPermission(ctx context.Context, perm *domain.ObjectPermission) error
	DeletePermission(ctx context.Context, userID string, objectType string, objectID string, action string) error
	ListPermissions(ctx context.Context) ([]*domain.ObjectPermission, error)
	GetEffectivePermissions(ctx context.Context, userID string, objectType string, objectID string) ([]string, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PermissionRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT id, parent_id, 1 as depth
			FROM folders
			WHERE id = $1
			UNION ALL
			SELECT f.id, f.parent_id, ft.depth + 1
			FROM folders f
			JOIN folder_tree ft ON f.id = ft.parent_id
		)
		SELECT EXISTS (
			SELECT 1 
			FROM folder_tree ft
			JOIN object_permissions op ON op.object_id = ft.id AND op.object_type = 'folder'
			WHERE op.user_id = $2 AND op.action = $3
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, folderID, userID, string(action)).Scan(&exists)
	return exists, err
}

func (r *postgresRepository) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error) {
	queryDirect := `
		SELECT EXISTS (
			SELECT 1 
			FROM object_permissions 
			WHERE user_id = $1 AND object_type = 'metadata' AND object_id = $2 AND action = $3
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, queryDirect, userID, metadataID, string(action)).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	queryFolder := `
		SELECT folder_id FROM metadata WHERE id = $1
	`
	var folderID string
	err = r.db.QueryRowContext(ctx, queryFolder, metadataID).Scan(&folderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return r.CheckFolderPermission(ctx, userID, folderID, action)
}

func (r *postgresRepository) AddPermission(ctx context.Context, perm *domain.ObjectPermission) error {
	query := `
		INSERT INTO object_permissions (user_id, object_type, object_id, action)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT unique_user_object_action DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, perm.UserID, perm.ObjectType, perm.ObjectID, perm.Action)
	return err
}

func (r *postgresRepository) DeletePermission(ctx context.Context, userID string, objectType string, objectID string, action string) error {
	query := `
		DELETE FROM object_permissions
		WHERE user_id = $1 AND object_type = $2 AND object_id = $3 AND action = $4
	`
	_, err := r.db.ExecContext(ctx, query, userID, objectType, objectID, action)
	return err
}

func (r *postgresRepository) ListPermissions(ctx context.Context) ([]*domain.ObjectPermission, error) {
	query := `
		SELECT id, user_id, object_type, object_id, action, created_at
		FROM object_permissions
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*domain.ObjectPermission
	for rows.Next() {
		p := &domain.ObjectPermission{}
		err := rows.Scan(&p.ID, &p.UserID, &p.ObjectType, &p.ObjectID, &p.Action, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *postgresRepository) GetEffectivePermissions(ctx context.Context, userID string, objectType string, objectID string) ([]string, error) {
	var isAdmin bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM user_roles WHERE user_id = $1 AND role_name = 'trainer_admin'
		)
	`, userID).Scan(&isAdmin)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return []string{"read", "write", "manage_permissions"}, nil
	}

	actions := []string{}
	allActions := []domain.Action{domain.ActionRead, domain.ActionWrite, domain.ActionManagePermissions}
	for _, action := range allActions {
		var allowed bool
		var err error
		if objectType == "folder" {
			allowed, err = r.CheckFolderPermission(ctx, userID, objectID, action)
		} else {
			allowed, err = r.CheckMetadataPermission(ctx, userID, objectID, action)
		}
		if err != nil {
			return nil, err
		}
		if allowed {
			actions = append(actions, string(action))
		}
	}
	return actions, nil
}
