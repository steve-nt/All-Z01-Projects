#!/usr/bin/env bash
set -euo pipefail

# 1) Read the key from secret.key
if [[ ! -r secret.key ]]; then
  echo "Error: secret.key not found or unreadable" >&2
  exit 1
fi
KEY=$(<secret.key)

# 2) Build, injecting it via -ldflags
go build -ldflags "-X 'main.MapQuestKey=${KEY}'" -o myapp main.go
echo "✅ Built myapp with MAPQUEST_KEY injected"
