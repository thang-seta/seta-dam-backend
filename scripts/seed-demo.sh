#!/bin/bash
# seed-demo.sh
# Seeds demo data into the running PostgreSQL container

set -e

echo "Seeding demo data into Postgres..."
docker compose exec -T db psql -U postgres -d setadam -f /flyway/seeds/demo_data.sql

echo "Demo seeding completed successfully!"
