@echo off
echo Running Flyway migrations...
docker compose run --rm flyway migrate
echo Migrations completed successfully!
