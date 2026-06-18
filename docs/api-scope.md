# API Scope

This document specifies the endpoints and operations exposed by both the GraphQL Gateway and the internal Go Service.

## GraphQL Gateway Schema (Port 4000)

### Queries
- `folderTree: [Folder!]!`: Get the full hierarchical tree structure of folders allowed for the user. Exposes `id`, `name`, `parent_id`, `description`, `created_by`, `updated_by`, `created_at`, `updated_at`.
- `folder(id: ID!): Folder`: Get detailed metadata for a specific folder.
- `metadataList(folderId: ID!): [Metadata!]!`: List all text metadata objects within a folder.
- `metadataDetail(id: ID!): Metadata`: Get the full details of a specific metadata item. Exposes `id`, `folder_id`, `title`, `description`, `labels`, `category`, `external_source`, `external_id`, `source_url`, `thumbnail_url`, `license`, `author`, `metadata_json`, `notes`, `created_by`, `updated_by`, `created_at`, `updated_at`.
- `effectivePermissions(userId: ID!, objectType: String!, objectId: ID!): EffectivePermissions!`: Get list of actions (`read`, `write`, `manage_permissions`) a user can execute on the target folder/metadata.

### Mutations
- `createFolder(name: String!, description: String, parentId: ID): Folder!`: Create a new folder under a parent.
- `updateFolder(id: ID!, name: String!, description: String): Folder!`: Update a folder's name and description.
- `moveFolder(id: ID!, parentId: ID): Folder!`: Move a folder to another parent (checking for cycles).
- `deleteFolder(id: ID!): Boolean!`: Delete a folder (only if empty).
- `createMetadata(folderId: ID!, title: String!, description: String, labels: [String!], category: String, externalSource: String, externalId: String, sourceUrl: String, thumbnailUrl: String, license: String, author: String, metadataJson: String, notes: String): Metadata!`: Add a metadata entry to a folder.
- `updateMetadata(id: ID!, folderId: ID!, title: String!, description: String, labels: [String!], category: String, externalSource: String, externalId: String, sourceUrl: String, thumbnailUrl: String, license: String, author: String, metadataJson: String, notes: String): Metadata!`: Update metadata details.
- `deleteMetadata(id: ID!): Boolean!`: Delete a metadata item.
- `grantPermission(userId: ID!, objectType: String!, objectId: ID!, action: String!): Boolean!`: Grant access. `objectType` accepts `folder`, `metadata`, or `metadata_item`; `action` accepts `read`, `write`, or `manage_permissions`.
- `revokePermission(userId: ID!, objectType: String!, objectId: ID!, action: String!): Boolean!`: Revoke access. `objectType` accepts `folder`, `metadata`, or `metadata_item`; `action` accepts `read`, `write`, or `manage_permissions`.

---

## Go REST API Endpoints (Port 8080)
These endpoints are internal and expect headers: `X-User-ID`, `X-User-Role`, `X-User-Username`.

- `POST /api/folders`: Create folder
- `GET /api/folders/tree`: Get visible folders
- `GET /api/folders/{id}`: Get folder by ID
- `PUT /api/folders/{id}`: Update folder name and description, returning the updated folder.
- `PUT /api/folders/{id}/move`: Move folder, returning the updated folder.
- `DELETE /api/folders/{id}`: Delete folder
- `POST /api/metadata`: Create metadata
- `GET /api/metadata/{id}`: Get metadata by ID
- `PUT /api/metadata/{id}`: Update metadata, returning the updated metadata item.
- `DELETE /api/metadata/{id}`: Delete metadata
- `GET /api/folders/{folderId}/metadata`: List metadata items in folder
- `POST /api/permissions`: Add permission
- `DELETE /api/permissions`: Delete permission
- `GET /api/permissions/effective`: Get effective permissions

## Error Format

The Go service returns controlled JSON errors:

```json
{ "error": "message" }
```

Common status mappings:

- `400 Bad Request`: missing required fields, invalid folder move, non-empty folder delete, invalid `metadata_json`, invalid permission `object_type`, or invalid permission `action`.
- `401 Unauthorized`: missing user context for protected operations.
- `403 Forbidden`: authenticated user does not have the required role or object-level permission.
- `404 Not Found`: folder or metadata item does not exist.

The GraphQL gateway validates non-empty folder/metadata required fields, validates `metadataJson`, and forwards user context to the Go service through `X-User-ID`, `X-User-Role`, and `X-User-Username`.
