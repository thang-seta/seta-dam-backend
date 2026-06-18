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
	List(ctx context.Context) ([]*domain.Metadata, error)
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
		INSERT INTO metadata_items (
			id, folder_id, title, description, labels, category, 
			external_source, external_id, source_url, thumbnail_url, 
			license, author, metadata_json, notes, created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	metaJSON := meta.MetadataJSON
	if metaJSON == "" {
		metaJSON = "{}"
	}
	return r.db.QueryRowContext(ctx, query,
		meta.ID, meta.FolderID, meta.Title, meta.Description, pq.Array(meta.Labels), meta.Category,
		meta.ExternalSource, meta.ExternalID, meta.SourceURL, meta.ThumbnailURL,
		meta.License, meta.Author, metaJSON, meta.Notes, meta.CreatedBy,
	).Scan(&meta.CreatedAt, &meta.UpdatedAt)
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.Metadata, error) {
	query := `
		SELECT id, folder_id, title, COALESCE(description, ''), COALESCE(labels, ARRAY[]::text[]), COALESCE(category, ''), 
		       COALESCE(external_source, ''), COALESCE(external_id, ''), COALESCE(source_url, ''), COALESCE(thumbnail_url, ''), 
		       COALESCE(license, ''), COALESCE(author, ''), COALESCE(metadata_json, '{}'::jsonb), COALESCE(notes, ''), created_by, updated_by, created_at, updated_at, deleted_at
		FROM metadata_items
		WHERE id = $1 AND deleted_at IS NULL
	`
	m := &domain.Metadata{}
	var metaJSONBytes []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.FolderID, &m.Title, &m.Description, pq.Array(&m.Labels), &m.Category,
		&m.ExternalSource, &m.ExternalID, &m.SourceURL, &m.ThumbnailURL,
		&m.License, &m.Author, &metaJSONBytes, &m.Notes, &m.CreatedBy, &m.UpdatedBy, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMetadataNotFound
		}
		return nil, err
	}
	m.MetadataJSON = string(metaJSONBytes)
	return m, nil
}

func (r *postgresRepository) Update(ctx context.Context, meta *domain.Metadata) error {
	query := `
		UPDATE metadata_items
		SET folder_id = $2, title = $3, description = $4, labels = $5, category = $6, 
		    external_source = $7, external_id = $8, source_url = $9, thumbnail_url = $10, 
		    license = $11, author = $12, metadata_json = $13, notes = $14, updated_by = $15, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	metaJSON := meta.MetadataJSON
	if metaJSON == "" {
		metaJSON = "{}"
	}
	res, err := r.db.ExecContext(ctx, query,
		meta.ID, meta.FolderID, meta.Title, meta.Description, pq.Array(meta.Labels), meta.Category,
		meta.ExternalSource, meta.ExternalID, meta.SourceURL, meta.ThumbnailURL,
		meta.License, meta.Author, metaJSON, meta.Notes, meta.UpdatedBy,
	)
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
		UPDATE metadata_items
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
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

func (r *postgresRepository) List(ctx context.Context) ([]*domain.Metadata, error) {
	query := `
		SELECT id, folder_id, title, COALESCE(description, ''), COALESCE(labels, ARRAY[]::text[]), COALESCE(category, ''), 
		       COALESCE(external_source, ''), COALESCE(external_id, ''), COALESCE(source_url, ''), COALESCE(thumbnail_url, ''), 
		       COALESCE(license, ''), COALESCE(author, ''), COALESCE(metadata_json, '{}'::jsonb), COALESCE(notes, ''), created_by, updated_by, created_at, updated_at, deleted_at
		FROM metadata_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC, title ASC
	`
	return r.list(ctx, query)
}

func (r *postgresRepository) ListByFolder(ctx context.Context, folderID string) ([]*domain.Metadata, error) {
	query := `
		SELECT id, folder_id, title, COALESCE(description, ''), COALESCE(labels, ARRAY[]::text[]), COALESCE(category, ''), 
		       COALESCE(external_source, ''), COALESCE(external_id, ''), COALESCE(source_url, ''), COALESCE(thumbnail_url, ''), 
		       COALESCE(license, ''), COALESCE(author, ''), COALESCE(metadata_json, '{}'::jsonb), COALESCE(notes, ''), created_by, updated_by, created_at, updated_at, deleted_at
		FROM metadata_items
		WHERE folder_id = $1 AND deleted_at IS NULL
		ORDER BY title ASC
	`
	return r.list(ctx, query, folderID)
}

func (r *postgresRepository) list(ctx context.Context, query string, args ...interface{}) ([]*domain.Metadata, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*domain.Metadata, 0)
	for rows.Next() {
		m := &domain.Metadata{}
		var metaJSONBytes []byte
		err := rows.Scan(
			&m.ID, &m.FolderID, &m.Title, &m.Description, pq.Array(&m.Labels), &m.Category,
			&m.ExternalSource, &m.ExternalID, &m.SourceURL, &m.ThumbnailURL,
			&m.License, &m.Author, &metaJSONBytes, &m.Notes, &m.CreatedBy, &m.UpdatedBy, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		m.MetadataJSON = string(metaJSONBytes)
		list = append(list, m)
	}
	return list, nil
}
