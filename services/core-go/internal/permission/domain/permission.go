package domain

import (
	"errors"
	"time"
)

type Action string

const (
	ActionRead              Action = "read"
	ActionWrite             Action = "write"
	ActionManagePermissions Action = "manage_permissions"
)

type ObjectType string

const (
	ObjectTypeFolder   ObjectType = "folder"
	ObjectTypeMetadata ObjectType = "metadata"
)

type ObjectPermission struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ObjectType string    `json:"object_type"`
	ObjectID   string    `json:"object_id"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserContext struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var (
	ErrUnauthorized      = errors.New("unauthorized action")
	ErrForbidden         = errors.New("forbidden resource access")
	ErrInvalidObjectType = errors.New("object_type must be folder, metadata, or metadata_item")
	ErrInvalidAction     = errors.New("action must be read, write, or manage_permissions")
)
