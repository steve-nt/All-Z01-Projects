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
mkdir -p certs

if [ ! -f certs/cert.pem ] || [ ! -f certs/key.pem ]; then
  echo "Error: TLS certificate files not found in ./certs"
  echo "Generate them first, for example:"
  echo 'openssl req -x509 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"'
  exit 1
fi

docker image build -t forum .
docker rm -f forum 2>/dev/null || true

docker container run --name forum \
  -p 8080:8080 \
  -p 8443:8443 \
  --env-file .env \
  -e DB_PATH=/app/data/forum.db \
  -e TLS_CERT_FILE=/app/certs/cert.pem \
  -e TLS_KEY_FILE=/app/certs/key.pem \
  -v "$(pwd)/data:/app/data" \
  forum