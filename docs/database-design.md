# Database Design

This document details the PostgreSQL schema designed for the SETA DAM backend.

## Entity Relationship Diagram

```mermaid
erDiagram
    users {
        uuid id PK
        varchar email UK
        varchar display_name
        boolean is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    roles {
        uuid id PK
        varchar code UK
        varchar name
        text description
        timestamp created_at
        timestamp updated_at
    }
    user_roles {
        uuid id PK
        uuid user_id FK
        uuid role_id FK
        timestamp assigned_at
    }
    permission_actions {
        uuid id PK
        permission_action_code code UK
        text description
    }
    role_permissions {
        uuid id PK
        uuid role_id FK
        uuid action_id FK
        resource_type resource_type
    }
    folders {
        uuid id PK
        uuid parent_id FK
        varchar name
        text description
        uuid created_by FK
        uuid updated_by FK
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    metadata_items {
        uuid id PK
        uuid folder_id FK
        varchar title
        text description
        text_array labels
        varchar category
        varchar external_source
        varchar external_id
        text source_url
        text thumbnail_url
        varchar license
        varchar author
        jsonb metadata_json
        text notes
        uuid created_by FK
        uuid updated_by FK
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    folder_permissions {
        uuid id PK
        uuid folder_id FK
        uuid grantee_user_id FK
        uuid grantee_role_id FK
        uuid action_id FK
        uuid granted_by FK
        timestamp granted_at
    }
    metadata_permissions {
        uuid id PK
        uuid metadata_item_id FK
        uuid grantee_user_id FK
        uuid grantee_role_id FK
        uuid action_id FK
        uuid granted_by FK
        timestamp granted_at
    }

    users ||--o{ user_roles : has
    roles ||--o{ user_roles : assigned
    roles ||--o{ role_permissions : defines
    permission_actions ||--o{ role_permissions : maps
    users ||--o{ folders : creates
    folders ||--o{ folders : contains
    folders ||--o{ metadata_items : holds
    users ||--o{ metadata_items : creates
    folders ||--o{ folder_permissions : protects
    metadata_items ||--o{ metadata_permissions : protects
    permission_actions ||--o{ folder_permissions : grants
    permission_actions ||--o{ metadata_permissions : grants
```

## Schema Details

### Tables
1. **users**: Stores identity credentials and profile info.
2. **roles**: Master list of role classifications (`trainer_admin`, `editor`, `viewer`).
3. **user_roles**: Joins users to roles.
4. **permission_actions**: Maps actions (`read`, `write`, `manage_permissions`) to specific enum values.
5. **role_permissions**: Maps roles to global actions and resource types.
6. **folders**: Represents hierarchical folder structures, supporting description and creator/updater fields, along with soft-delete metadata.
7. **metadata_items**: Stores flat, image-less metadata files mapping text descriptions, labels, categories, external source IDs, licensing details, JSON custom metadata, and soft-delete states.
8. **folder_permissions**: Direct ACL mapping of privileges (`read`, `write`, `manage_permissions`) on folders to a specific user or role.
9. **metadata_permissions**: Direct ACL mapping of privileges on metadata items to a specific user or role.

### Constraints & Indexes
- Foreign key constraints ensure strict referential integrity.
- Unique constraints (`uq_folder_permissions_user`, `uq_folder_permissions_role`, etc.) prevent duplicate permissions mappings.
- Parent ID index speeds up tree traversal.
- GIN indexes speed up `labels` array and `metadata_json` searches.
- Index filter `deleted_at IS NULL` supports efficient soft-deleted record exclusion.

