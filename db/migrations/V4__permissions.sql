-- V4__permissions.sql
-- Create object-level permissions table

CREATE TABLE IF NOT EXISTS object_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_type VARCHAR(50) NOT NULL CHECK (object_type IN ('folder', 'metadata')),
    object_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL CHECK (action IN ('read', 'write', 'manage_permissions')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_object_action UNIQUE (user_id, object_type, object_id, action)
);
