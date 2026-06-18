#!/bin/bash
# seed-demo.sh
# Seeds demo data into the running PostgreSQL container

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Seeding demo data into Postgres..."
docker compose exec -T db psql -U postgres -d setadam < db/seeds/demo_data.sql

echo "Demo seeding completed successfully!"
