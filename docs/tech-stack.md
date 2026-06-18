# Technology Stack

This document details the technologies chosen for the SETA DAM backend Phase 1 MVP.

## GraphQL API Layer
- **Node.js (TypeScript)**: Offers a highly dynamic schema orchestration layer with Apollo Server, permitting rapid evolution of front-facing schemas.
- **Apollo Server (Express Integration)**: Industry standard GraphQL engine that integrates seamlessly with Express and supports schema modularization.

## Core Domain Service
- **Go (Golang)**: Chosen for its stellar concurrency capabilities, low resource footprint, fast compilation, and simple, robust execution model.
- **Chi Router**: A lightweight, idiomatic HTTP router for Go that compiles routes quickly and conforms to standard `net/http` signatures.

## Relational Database & Migration
- **PostgreSQL**: Robust, ACID-compliant open-source relational database that natively supports recursive CTEs (essential for folder tree hierarchies).
- **Flyway**: Standard database migration tool that applies SQL version control sequentially and makes database environments easily replicable.

## Monitoring & Observability
- **Prometheus**: Pull-based metric scrape daemon that records runtime performance.
- **Grafana**: Web interface for dashboard visualization of Prometheus metrics and Loki logs.
- **Loki**: Log aggregation system inspired by Prometheus.
