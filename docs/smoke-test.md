# Phase 1 Smoke Test

This smoke test verifies the Phase 1 acceptance baseline through the GraphQL gateway and Go core service.

## Prerequisites

Start the stack first:

```bash
npm run dev:up
```

The script expects:

- GraphQL gateway: `http://localhost:4000/graphql`
- GraphQL health: `http://localhost:4000/health`
- Go service health: `http://localhost:8080/health`
- Seeded mock users: `admin_user`, `editor_user`, `viewer_user`

Override endpoints when needed:

```bash
GRAPHQL_URL=http://localhost:4000/graphql CORE_HEALTH_URL=http://localhost:8080/health NODE_HEALTH_URL=http://localhost:4000/health npm run smoke:test
```

## Run

```bash
npm run smoke:test
```

## Coverage

The script checks:

- Both service health endpoints are reachable.
- The admin user can query the folder tree.
- The admin/editor flow can create folders and metadata through GraphQL.
- The viewer user can read metadata.
- The viewer user is denied when attempting to update metadata.
- Effective permissions can be queried using `metadata_item` object type.
- Smoke-test folders and metadata are cleaned up after the run.
