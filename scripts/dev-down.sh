#!/bin/bash
# dev-down.sh
# Shuts down the development environment and cleans up volumes

echo "Stopping Docker Compose services..."
docker compose down -v

echo "Environment torn down successfully."
