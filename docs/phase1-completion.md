# Phase 1 Completion Notes

This project implements the Phase 1 backend-only MVP for image text metadata management.

## Completed Functional Areas

- Backend-only scope with a Node.js GraphQL gateway, Go core service, PostgreSQL, and Flyway.
- Demo users and roles for `trainer_admin`, `editor`, and `viewer`.
- Folder create, tree/list, detail, rename, move with cycle prevention, and safe soft-delete.
- Metadata create, list by folder, detail, update, and soft-delete with text-only fields.
- RBAC plus object-level permission evaluation for folders and metadata items.
- Permission grant, revoke, and effective-permissions debug flow.
- Health endpoints for both services.
- Database migrations and demo seeds for users, roles, permissions, folders, and metadata.
- Unit coverage for folder move behavior, metadata CRUD validation/authorization, and permission enforcement.
- Smoke-test script covering startup health, GraphQL operations, migration-backed seeded users, and permission denial.

## Verification

Local checks:

```bash
cd services/core-go && go test ./...
npm --prefix apps/graphql-node run build
python3 -m py_compile scripts/smoke-test.py
```

End-to-end smoke flow after Docker services are running:

```bash
npm run dev:up
npm run smoke:test
```

See [smoke-test.md](smoke-test.md) for details.
