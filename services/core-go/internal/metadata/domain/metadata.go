package domain

import (
	"errors"
	"time"
)

type Metadata struct {
	ID             string     `json:"id"`
	FolderID       string     `json:"folder_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Labels         []string   `json:"labels"`
	Category       string     `json:"category"`
	ExternalSource string     `json:"external_source"`
	ExternalID     string     `json:"external_id"`
	SourceURL      string     `json:"source_url"`
	ThumbnailURL   string     `json:"thumbnail_url"`
	License        string     `json:"license"`
	Author         string     `json:"author"`
	MetadataJSON   string     `json:"metadata_json"`
	Notes          string     `json:"notes"`
	CreatedBy      string     `json:"created_by"`
	UpdatedBy      *string    `json:"updated_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

var (
	ErrMetadataNotFound    = errors.New("metadata item not found")
	ErrTitleRequired       = errors.New("title is a required field")
	ErrFolderRequired      = errors.New("folder_id is a required field")
	ErrInvalidMetadataJSON = errors.New("metadata_json must be valid JSON")
)
