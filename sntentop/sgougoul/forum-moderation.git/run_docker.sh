#!/usr/bin/env bash
set -e

# Ensure we run from the project root (where .env lives)
cd "$(dirname "$0")"

if [ ! -f .env ]; then
  echo "Error: .env file not found. Create it in the project root."
  echo "Tip: copy from .env.example -> cp .env.example .env"
  exit 1
fi

# Persist SQLite DB + WAL/SHM files across container restarts
mkdir -p data

docker image build -t forum .
docker rm -f forum 2>/dev/null || true

docker container run --name forum \
  -p 8080:8080 \
  --env-file .env \
  -e DB_PATH=/app/data/forum.db \
  -v "$(pwd)/data:/app/data" \
  forum