package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/user/seta-dam-backend/services/core-go/internal/metadata/domain"
)

type MetadataRepository interface {
	Create(ctx context.Context, meta *domain.Metadata) error
	GetByID(ctx context.Context, id string) (*domain.Metadata, error)
	Update(ctx context.Context, meta *domain.Metadata) error
	Delete(ctx context.Context, id string) error
	ListByFolder(ctx context.Context, folderID string) ([]*domain.Metadata, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) MetadataRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, meta *domain.Metadata) error {
	query := `
		INSERT INTO metadata (id, folder_id, title, description, labels, category, source_url, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query, meta.ID, meta.FolderID, meta.Title, meta.Description, pq.Array(meta.Labels), meta.Category, meta.SourceURL, meta.Notes).Scan(&meta.CreatedAt)
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.Metadata, error) {
	query := `
		SELECT id, folder_id, title, description, labels, category, source_url, notes, created_at
		FROM metadata
		WHERE id = $1
	`
	m := &domain.Metadata{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.FolderID, &m.Title, &m.Description, pq.Array(&m.Labels), &m.Category, &m.SourceURL, &m.Notes, &m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMetadataNotFound
		}
		return nil, err
	}
	return m, nil
}

func (r *postgresRepository) Update(ctx context.Context, meta *domain.Metadata) error {
	query := `
		UPDATE metadata
		SET folder_id = $2, title = $3, description = $4, labels = $5, category = $6, source_url = $7, notes = $8
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, meta.ID, meta.FolderID, meta.Title, meta.Description, pq.Array(meta.Labels), meta.Category, meta.SourceURL, meta.Notes)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrMetadataNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM metadata
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
		return domain.ErrMetadataNotFound
	}
	return nil
}

func (r *postgresRepository) ListByFolder(ctx context.Context, folderID string) ([]*domain.Metadata, error) {
	query := `
		SELECT id, folder_id, title, description, labels, category, source_url, notes, created_at
		FROM metadata
		WHERE folder_id = $1
		ORDER BY title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Metadata
	for rows.Next() {
		m := &domain.Metadata{}
		err := rows.Scan(
			&m.ID, &m.FolderID, &m.Title, &m.Description, pq.Array(&m.Labels), &m.Category, &m.SourceURL, &m.Notes, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}
