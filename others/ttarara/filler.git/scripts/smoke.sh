#!/usr/bin/env bash
# Quick one-shot game vs wall_e on map00, printing full engine output so
# we can eyeball the protocol (coordinate order, score format, etc.).
# Used in Phase 4 to answer the `X Y` vs `Y X` question.

set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

missing=0
for path in linux_game_engine m1_game_engine linux_robots m1_robots maps; do
    if [ ! -e "$REPO_ROOT/docker_image/$path" ]; then
        echo "Missing docker_image/$path from the 01-edu filler bundle." >&2
        missing=1
    fi
done
if [ "$missing" -ne 0 ]; then
    echo "Add the official bundle files to docker_image/, run scripts/build-image.sh, then rerun this script." >&2
    exit 1
fi

if ! docker image inspect filler >/dev/null 2>&1; then
    echo "Image 'filler' not found. Run scripts/build-image.sh first." >&2
    exit 1
fi

docker run --rm \
    -v "$REPO_ROOT/docker_image/solution":/filler/solution \
    --entrypoint bash filler \
    -c '
set -euo pipefail
cd /filler/solution && rm -f Cargo.lock && cargo build --release 2>&1 | tail -3
cd /filler
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ]; then
    ENGINE=./m1_game_engine; ROBOTS=m1_robots
else
    ENGINE=./linux_game_engine; ROBOTS=linux_robots
fi
echo "===== full engine output ====="
"$ENGINE" -f maps/map00 -p1 solution/target/release/filler -p2 "$ROBOTS/wall_e"
'
