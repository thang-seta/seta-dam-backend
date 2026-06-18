package repository

import (
	"context"
	"database/sql"
	"github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type PermissionRepository interface {
	CheckFolderPermission(ctx context.Context, userID string, folderID string, action domain.Action) (bool, error)
	CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error)
	AddPermission(ctx context.Context, perm *domain.ObjectPermission, grantedByID string) error
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
	// First, check global role-based permission
	queryGlobal := `
		SELECT EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN role_permissions rp ON ur.role_id = rp.role_id
			JOIN permission_actions pa ON rp.action_id = pa.id
			WHERE ur.user_id = $1 AND rp.resource_type = 'folder' AND pa.code = $2
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, queryGlobal, userID, string(action)).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	// Second, check explicit/inherited folder permissions
	queryLocal := `
		WITH RECURSIVE folder_tree AS (
			SELECT id, parent_id, 1 as depth
			FROM folders
			WHERE id = $1 AND deleted_at IS NULL
			UNION ALL
			SELECT f.id, f.parent_id, ft.depth + 1
			FROM folders f
			JOIN folder_tree ft ON f.id = ft.parent_id
			WHERE f.deleted_at IS NULL
		)
		SELECT EXISTS (
			SELECT 1 FROM folder_tree ft
			JOIN folder_permissions fp ON fp.folder_id = ft.id
			JOIN permission_actions pa ON fp.action_id = pa.id
			LEFT JOIN user_roles ur ON ur.role_id = fp.grantee_role_id AND ur.user_id = $2
			WHERE pa.code = $3 AND (fp.grantee_user_id = $2 OR ur.id IS NOT NULL)
		)
	`
	err = r.db.QueryRowContext(ctx, queryLocal, folderID, userID, string(action)).Scan(&exists)
	return exists, err
}

func (r *postgresRepository) CheckMetadataPermission(ctx context.Context, userID string, metadataID string, action domain.Action) (bool, error) {
	// 1. Check global role-based permission for metadata_item
	queryGlobal := `
		SELECT EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN role_permissions rp ON ur.role_id = rp.role_id
			JOIN permission_actions pa ON rp.action_id = pa.id
			WHERE ur.user_id = $1 AND rp.resource_type = 'metadata_item' AND pa.code = $2
		)
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, queryGlobal, userID, string(action)).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	// 2. Check explicit metadata-level permission
	queryDirect := `
		SELECT EXISTS (
			SELECT 1 FROM metadata_permissions mp
			JOIN permission_actions pa ON mp.action_id = pa.id
			LEFT JOIN user_roles ur ON ur.role_id = mp.grantee_role_id AND ur.user_id = $1
			WHERE mp.metadata_item_id = $2 AND pa.code = $3 AND (mp.grantee_user_id = $1 OR ur.id IS NOT NULL)
		)
	`
	err = r.db.QueryRowContext(ctx, queryDirect, userID, metadataID, string(action)).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	// 3. Fallback to checking the folder's permissions
	queryFolder := `
		SELECT folder_id FROM metadata_items WHERE id = $1 AND deleted_at IS NULL
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

func (r *postgresRepository) AddPermission(ctx context.Context, perm *domain.ObjectPermission, grantedByID string) error {
	var actionID string
	err := r.db.QueryRowContext(ctx, `SELECT id FROM permission_actions WHERE code = $1`, perm.Action).Scan(&actionID)
	if err != nil {
		return err
	}

	if perm.ObjectType == "folder" {
		query := `
			INSERT INTO folder_permissions (folder_id, grantee_user_id, action_id, granted_by)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (folder_id, grantee_user_id, action_id) DO NOTHING
		`
		_, err = r.db.ExecContext(ctx, query, perm.ObjectID, perm.UserID, actionID, grantedByID)
		return err
	} else {
		query := `
			INSERT INTO metadata_permissions (metadata_item_id, grantee_user_id, action_id, granted_by)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (metadata_item_id, grantee_user_id, action_id) DO NOTHING
		`
		_, err = r.db.ExecContext(ctx, query, perm.ObjectID, perm.UserID, actionID, grantedByID)
		return err
	}
}

func (r *postgresRepository) DeletePermission(ctx context.Context, userID string, objectType string, objectID string, action string) error {
	if objectType == "folder" {
		query := `
			DELETE FROM folder_permissions
			WHERE grantee_user_id = $1 AND folder_id = $2
			  AND action_id = (SELECT id FROM permission_actions WHERE code = $3)
		`
		_, err := r.db.ExecContext(ctx, query, userID, objectID, action)
		return err
	} else {
		query := `
			DELETE FROM metadata_permissions
			WHERE grantee_user_id = $1 AND metadata_item_id = $2
			  AND action_id = (SELECT id FROM permission_actions WHERE code = $3)
		`
		_, err := r.db.ExecContext(ctx, query, userID, objectID, action)
		return err
	}
}

func (r *postgresRepository) ListPermissions(ctx context.Context) ([]*domain.ObjectPermission, error) {
	query := `
		SELECT id, grantee_user_id as user_id, 'folder' as object_type, folder_id as object_id, 
		       (SELECT code FROM permission_actions WHERE id = action_id)::text as action, granted_at as created_at
		FROM folder_permissions
		WHERE grantee_user_id IS NOT NULL
		UNION ALL
		SELECT id, grantee_user_id as user_id, 'metadata' as object_type, metadata_item_id as object_id, 
		       (SELECT code FROM permission_actions WHERE id = action_id)::text as action, granted_at as created_at
		FROM metadata_permissions
		WHERE grantee_user_id IS NOT NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := make([]*domain.ObjectPermission, 0)
	for rows.Next() {
		p := &domain.ObjectPermission{}
		err := rows.Scan(&p.ID, &p.UserID, &p.ObjectType, &p.ObjectID, &p.Action, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *postgresRepository) GetEffectivePermissions(ctx context.Context, userID string, objectType string, objectID string) ([]string, error) {
	var isAdmin bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.code = 'trainer_admin'
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
