@echo off
echo Seeding demo data into Postgres...
docker compose exec -T db psql -U postgres -d setadam < "%~dp0..\db\seeds\demo_data.sql"
echo Demo seeding completed successfully!
