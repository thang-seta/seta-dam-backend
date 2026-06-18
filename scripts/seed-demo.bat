@echo off
cd /d "%~dp0.."
echo Seeding demo data into Postgres...
docker compose exec -T db psql -U postgres -d setadam < db/seeds/demo_data.sql
echo Demo seeding completed successfully!
