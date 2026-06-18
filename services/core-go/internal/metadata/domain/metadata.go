package domain

import (
	"errors"
	"time"
)

type Metadata struct {
	ID          string    `json:"id"`
	FolderID    string    `json:"folder_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Labels      []string  `json:"labels"`
	Category    string    `json:"category"`
	SourceURL   string    `json:"source_url"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	ErrMetadataNotFound = errors.New("metadata item not found")
	ErrTitleRequired    = errors.New("title is a required field")
	ErrFolderRequired   = errors.New("folder_id is a required field")
)
