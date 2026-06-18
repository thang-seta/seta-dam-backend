#!/bin/bash
# dev-up.sh
# Starts the local development environment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Starting Docker Compose services..."
docker compose up -d --build

echo "Services started successfully!"
echo "GraphQL Gateway: http://localhost:4000/graphql"
echo "Go REST API: http://localhost:8080/health"
echo "Prometheus: http://localhost:9090"
echo "Grafana: http://localhost:3000"
