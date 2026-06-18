@echo off
cd /d "%~dp0.."
echo Stopping Docker Compose services...
docker compose down -v
echo Environment torn down successfully.
