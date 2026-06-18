package domain

import (
	"errors"
	"time"
)

type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrFolderNotFound = errors.New("folder not found")
	ErrCycleDetected  = errors.New("cannot move folder: cycle detected")
	ErrNotEmpty       = errors.New("folder is not empty (contains subfolders or metadata)")
)
