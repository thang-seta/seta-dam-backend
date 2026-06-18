@echo off
cd /d "%~dp0.."
python scripts/import-open-images-v7.py %*
