#!/usr/bin/env bash
# One-time: build the 01-edu `filler` Docker image.
# Re-run only if docker_image/Dockerfile changes.

set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT/docker_image"

missing=0
for path in linux_game_engine m1_game_engine linux_robots m1_robots maps; do
    if [ ! -e "$path" ]; then
        echo "Missing docker_image/$path from the 01-edu filler bundle." >&2
        missing=1
    fi
done
if [ "$missing" -ne 0 ]; then
    echo "Add the official bundle files to docker_image/, then rerun this script." >&2
    exit 1
fi

docker build -t filler .
