#!/bin/bash
# migrate.sh
# Triggers Flyway schema migrations

set -e

echo "Running Flyway migrations..."
docker compose run --rm flyway migrate

echo "Migrations completed successfully!"
