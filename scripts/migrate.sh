#!/bin/bash
# migrate.sh
# Triggers Flyway schema migrations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Running Flyway migrations..."
docker compose run --rm flyway migrate

echo "Migrations completed successfully!"
