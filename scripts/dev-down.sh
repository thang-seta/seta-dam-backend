#!/bin/bash
# dev-down.sh
# Shuts down the development environment and cleans up volumes

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Stopping Docker Compose services..."
docker compose down -v

echo "Environment torn down successfully."
