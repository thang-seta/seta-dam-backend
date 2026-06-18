# SETA DAM Backend MVP

This project implements the backend MVP for the SETA Digital Asset Manager (DAM). It is a backend-only system for managing image text metadata organized in a hierarchical folder tree.

The system features:
1. **Node.js (TypeScript) GraphQL Gateway**: Handles schema parsing, structural validation, and authentication parsing.
2. **Go Core Domain Service**: Processes folder CRUD, tree operations, cycle detection, metadata entries, and evaluates access permissions.
3. **PostgreSQL**: Relational database storage.
4. **Flyway**: Auto-manages database migrations.
5. **Observability Stack**: Prometheus metrics scraping, Loki log aggregator, and Grafana visualization dashboard.

---

## Technical Architecture

```mermaid
graph TD
    Client[GraphQL Playground] -->|GraphQL query/mutation| NodeGateway[GraphQL Gateway - Port 4000]
    NodeGateway -->|HTTP REST + User Context Headers| GoService[Go Core Service - Port 8080]
    GoService -->|SQL Query| Postgres[(PostgreSQL DB - Port 5432)]
    
    Flyway[Flyway Migrator] -->|Migrations & Seeds| Postgres
    
    %% Monitoring
    Prometheus[Prometheus] -->|Scrapes Metrics| NodeGateway
    Prometheus -->|Scrapes Metrics| GoService
    Grafana[Grafana] -->|Dashboards| Prometheus
    Grafana -->|Dashboards| Loki[Loki Logs]
```

For more detailed information, please read the documents inside the `/docs` directory:
- [docs/architecture.md](docs/architecture.md)
- [docs/tech-stack.md](docs/tech-stack.md)
- [docs/api-scope.md](docs/api-scope.md)
- [docs/database-design.md](docs/database-design.md)
- [docs/permission-design.md](docs/permission-design.md)

---

## Folder Structure

```text
seta-dam-backend/
├── apps/
│   └── graphql-node/          # Node.js GraphQL API gateway
├── services/
│   └── core-go/               # Go REST core domain service
├── db/
│   ├── migrations/            # Flyway migrations
│   └── seeds/                 # Demo SQL seed files
├── infra/
│   ├── prometheus/            # Prometheus configuration
│   ├── loki/                  # Loki configuration
│   └── grafana/               # Grafana dashboards
├── docs/                      # Technical specification documents
└── scripts/                   # Local development control scripts
```

---

## Local Development Setup

### Prerequisites
- Docker & Docker Compose installed.
- Node.js (for NPM commands) and Python 3.

### Recommended: Using NPM Commands (Cross-Platform)

For convenience, you can manage the lifecycle of the application using NPM scripts from the root directory:

* **Start all services in Docker**:
  ```bash
  npm run dev:up
  ```
* **Stop all services and clean DB volume**:
  ```bash
  npm run dev:down
  ```
* **Run Database Migrations manually**:
  ```bash
  npm run db:migrate
  ```
* **Load Demo Seed Data**:
  ```bash
  npm run db:seed
  ```
* **Import Open Images Metadata**:
  ```bash
  npm run db:import -- --limit 100
  ```
* **Local Hybrid Development Mode** (Starts DB in Docker, runs Go & Node.js Gateway locally with hot-reload):
  ```bash
  npm run dev:all
  ```
* **Stop all containers (without deleting volumes)**:
  ```bash
  npm run shut:all
  ```
* **View logs**:
  ```bash
  npm run dev:logs
  ```

---

### Alternative: Using Shell Scripts

If you prefer to run the scripts directly, we provide scripts for both Linux/macOS (`.sh`) and Windows (`.bat`):

| Action | Linux / macOS | Windows (CMD/PowerShell) |
|---|---|---|
| **Start Services** | `./scripts/dev-up.sh` | `.\scripts\dev-up.bat` |
| **Stop & Clean** | `./scripts/dev-down.sh` | `.\scripts\dev-down.bat` |
| **Run Migrations** | `./scripts/migrate.sh` | `.\scripts\migrate.bat` |
| **Load Seed Data** | `./scripts/seed-demo.sh` | `.\scripts\seed-demo.bat` |
| **Import Metadata** | `./scripts/import-open-images-v7.py --limit 100` | `.\scripts\import-open-images-v7.bat --limit 100` |
| **Import Dry-run** | `./scripts/import-open-images-v7.py --limit 10 --dry-run` | `.\scripts\import-open-images-v7.bat --limit 10 --dry-run` |

The importer creates an `Open Images V7 Import` root folder, creates label-based child folders, and upserts rows into `metadata_items` using `external_source = 'open_images_v7_validation'`.

---

## Codebase Walkthrough

### Mock Authentication & Context Passing
To perform query actions as different roles:
Send an `Authorization` header with one of the mock user usernames:
- `admin_user`: Bypasses all permissions (`trainer_admin` role).
- `editor_user`: Can write inside permitted folders (`editor` role).
- `viewer_user`: Read-only user (`viewer` role).

Example GraphQL Request Header:
```json
{
  "Authorization": "admin_user"
}
```

### Access Control Rules
The Go service enforces object-level access rules:
- Permissions set on parent folders automatically **inherit downwards** to subfolders and metadata items.
- Check folder operations (CRUD) against the folder tree ancestry using a PostgreSQL recursive CTE.
- Validate folder deletions to verify they do not contain child folders or metadata (preventing unsafe deletions).
- Block folder moves that would create cycles (e.g. moving a folder into its own subfolder).
