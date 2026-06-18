package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
)

type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) error
	GetByID(ctx context.Context, id string) (*domain.Folder, error)
	Update(ctx context.Context, folder *domain.Folder) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Folder, error)
	IsDescendantOf(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error)
	IsEmpty(ctx context.Context, id string) (bool, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) FolderRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, folder *domain.Folder) error {
	query := `
		INSERT INTO folders (id, name, parent_id)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query, folder.ID, folder.Name, folder.ParentID).Scan(&folder.CreatedAt)
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.Folder, error) {
	query := `
		SELECT id, name, parent_id, created_at
		FROM folders
		WHERE id = $1
	`
	f := &domain.Folder{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&f.ID, &f.Name, &f.ParentID, &f.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrFolderNotFound
		}
		return nil, err
	}
	return f, nil
}

func (r *postgresRepository) Update(ctx context.Context, folder *domain.Folder) error {
	query := `
		UPDATE folders
		SET name = $2, parent_id = $3
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, folder.ID, folder.Name, folder.ParentID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrFolderNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM folders
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrFolderNotFound
	}
	return nil
}

func (r *postgresRepository) List(ctx context.Context) ([]*domain.Folder, error) {
	query := `
		SELECT id, name, parent_id, created_at
		FROM folders
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*domain.Folder
	for rows.Next() {
		f := &domain.Folder{}
		err := rows.Scan(&f.ID, &f.Name, &f.ParentID, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *postgresRepository) IsDescendantOf(ctx context.Context, potentialDescendantID string, ancestorID string) (bool, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT id, parent_id
			FROM folders
			WHERE id = $1
			UNION ALL
			SELECT f.id, f.parent_id
			FROM folders f
			JOIN folder_tree ft ON f.id = ft.parent_id
		)
		SELECT EXISTS (
			SELECT 1 FROM folder_tree WHERE id = $2
		)
	`
	var isDescendant bool
	err := r.db.QueryRowContext(ctx, query, potentialDescendantID, ancestorID).Scan(&isDescendant)
	return isDescendant, err
}

func (r *postgresRepository) IsEmpty(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM folders WHERE parent_id = $1
			UNION ALL
			SELECT 1 FROM metadata WHERE folder_id = $1
		)
	`
	var containsElements bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&containsElements)
	if err != nil {
		return false, err
	}
	return !containsElements, nil
}
