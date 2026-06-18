CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE resource_type AS ENUM (
  'folder',
  'metadata_item'
);

CREATE TYPE permission_action_code AS ENUM (
  'read',
  'write',
  'manage_permissions'
);

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email varchar NOT NULL UNIQUE,
  display_name varchar NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE roles (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code varchar NOT NULL UNIQUE,
  name varchar NOT NULL,
  description text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_roles (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  role_id uuid NOT NULL REFERENCES roles(id),
  assigned_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (user_id, role_id)
);

CREATE TABLE permission_actions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code permission_action_code NOT NULL UNIQUE,
  description text
);

CREATE TABLE role_permissions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  role_id uuid NOT NULL REFERENCES roles(id),
  action_id uuid NOT NULL REFERENCES permission_actions(id),
  resource_type resource_type NOT NULL,
  UNIQUE (role_id, action_id, resource_type)
);

CREATE TABLE folders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  parent_id uuid REFERENCES folders(id),
  name varchar NOT NULL,
  description text,
  created_by uuid NOT NULL REFERENCES users(id),
  updated_by uuid REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE metadata_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  folder_id uuid NOT NULL REFERENCES folders(id),
  title varchar NOT NULL,
  description text,
  labels text[],
  category varchar,

  external_source varchar,
  external_id varchar,
  source_url text,
  thumbnail_url text,
  license varchar,
  author varchar,
  metadata_json jsonb,

  notes text,
  created_by uuid NOT NULL REFERENCES users(id),
  updated_by uuid REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE folder_permissions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  folder_id uuid NOT NULL REFERENCES folders(id),
  grantee_user_id uuid REFERENCES users(id),
  grantee_role_id uuid REFERENCES roles(id),
  action_id uuid NOT NULL REFERENCES permission_actions(id),
  granted_by uuid NOT NULL REFERENCES users(id),
  granted_at timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT chk_folder_permission_single_grantee CHECK (
    (grantee_user_id IS NOT NULL AND grantee_role_id IS NULL)
    OR
    (grantee_user_id IS NULL AND grantee_role_id IS NOT NULL)
  )
);

CREATE TABLE metadata_permissions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  metadata_item_id uuid NOT NULL REFERENCES metadata_items(id),
  grantee_user_id uuid REFERENCES users(id),
  grantee_role_id uuid REFERENCES roles(id),
  action_id uuid NOT NULL REFERENCES permission_actions(id),
  granted_by uuid NOT NULL REFERENCES users(id),
  granted_at timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT chk_metadata_permission_single_grantee CHECK (
    (grantee_user_id IS NOT NULL AND grantee_role_id IS NULL)
    OR
    (grantee_user_id IS NULL AND grantee_role_id IS NOT NULL)
  )
);

CREATE UNIQUE INDEX uq_metadata_items_external_source_id
ON metadata_items(external_source, external_id)
WHERE external_source IS NOT NULL
  AND external_id IS NOT NULL;

CREATE UNIQUE INDEX uq_folder_permissions_user
ON folder_permissions(folder_id, grantee_user_id, action_id)
WHERE grantee_user_id IS NOT NULL;

CREATE UNIQUE INDEX uq_folder_permissions_role
ON folder_permissions(folder_id, grantee_role_id, action_id)
WHERE grantee_role_id IS NOT NULL;

CREATE UNIQUE INDEX uq_metadata_permissions_user
ON metadata_permissions(metadata_item_id, grantee_user_id, action_id)
WHERE grantee_user_id IS NOT NULL;

CREATE UNIQUE INDEX uq_metadata_permissions_role
ON metadata_permissions(metadata_item_id, grantee_role_id, action_id)
WHERE grantee_role_id IS NOT NULL;

CREATE INDEX idx_folders_parent_id ON folders(parent_id);
CREATE INDEX idx_metadata_items_folder_id ON metadata_items(folder_id);
CREATE INDEX idx_metadata_items_labels_gin ON metadata_items USING gin(labels);
CREATE INDEX idx_metadata_items_metadata_json_gin ON metadata_items USING gin(metadata_json);