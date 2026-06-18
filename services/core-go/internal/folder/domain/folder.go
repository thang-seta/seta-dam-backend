package domain

import (
	"errors"
	"time"
)

type Folder struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parent_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   *string    `json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

var (
	ErrFolderNotFound = errors.New("folder not found")
	ErrNameRequired   = errors.New("name is a required field")
	ErrCycleDetected  = errors.New("cannot move folder: cycle detected")
	ErrNotEmpty       = errors.New("folder is not empty (contains subfolders or metadata)")
)
