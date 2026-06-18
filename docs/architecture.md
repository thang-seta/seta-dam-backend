# System Architecture

This document describes the high-level architecture of the SETA DAM backend.

## Architecture Diagram

```mermaid
graph TD
    Client[Client / Playground] -->|GraphQL Queries/Mutations| NodeGateway[Node.js GraphQL Gateway]
    NodeGateway -->|Internal REST API calls + User Context Headers| GoService[Go Core Service]
    GoService -->|SQL Reads/Writes| Postgres[(PostgreSQL Database)]
    
    Flyway[Flyway Migrator] -->|Runs SQL Schema & Seeds| Postgres
    
    %% Monitoring Stack
    Prometheus[Prometheus] -->|Scrapes Metrics| NodeGateway
    Prometheus -->|Scrapes Metrics| GoService
    Grafana[Grafana] -->|Visualizes Metrics| Prometheus
    Grafana -->|Visualizes Logs| Loki[Loki Log Aggregator]
```

## Request Flow
1. **Authentication (Mocked)**: The client provides user identity in the request header (e.g. `Authorization: admin_user`).
2. **GraphQL Gateway Layer**: Node.js extracts the header, builds a `Context` containing user ID, username, and role, and performs structural validation of query input.
3. **Internal API Routing**: Node.js forwards queries/mutations to the Go REST Service, propagating context details via standard headers (`X-User-ID`, `X-User-Role`, `X-User-Username`).
4. **Access Control (Go)**: The Go Service checks if the user has appropriate permissions:
   - For folders: checks the folder tree recursively (using recursive CTE) for direct/inherited permissions.
   - For metadata: checks direct metadata permissions or inherits folder-level permissions.
5. **Database Interaction**: If authorized, Go updates/fetches Postgres database values and returns the result.
