#!/usr/bin/env bash
# Build the bot and run the audit matchups inside the 01-edu filler
# container, alternating p1/p2 seats. Prints a W/L tally.
#
# Prereqs:
#   - docker_image/ contains the 01-edu bundle (Dockerfile + maps + both
#     engine variants + both robot variants)
#   - `scripts/build-image.sh` has been run (image tag: `filler`)
#   - Docker daemon is running

set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RUNS=${RUNS:-5}

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

run_in_container() {
    docker run --rm \
        -v "$REPO_ROOT/docker_image/solution":/filler/solution \
        --entrypoint bash filler \
        -c "$1"
}

echo "==> Building release binary inside container..."
run_in_container "cd /filler/solution && rm -f Cargo.lock && cargo build --release 2>&1 | tail -3"

echo
echo "==> Running bench..."
run_in_container "$(cat <<SCRIPT
set -euo pipefail
cd /filler

ARCH=\$(uname -m)
if [ "\$ARCH" = "aarch64" ]; then
    ENGINE=./m1_game_engine
    ROBOTS=m1_robots
else
    ENGINE=./linux_game_engine
    ROBOTS=linux_robots
fi
MINE=solution/target/release/filler
RUNS=$RUNS

parse_scores() {
    local all p1 p2
    all=\$(cat)
    p1=\$(echo "\$all" | sed -n 's/.*Player1 ([^)]*): \([0-9]\+\).*/\1/p' | tail -1)
    p2=\$(echo "\$all" | sed -n 's/.*Player2 ([^)]*): \([0-9]\+\).*/\1/p' | tail -1)
    echo "\${p1:-0} \${p2:-0}"
}

run_match() {
    local opponent="\$1" map="\$2"
    local wins=0
    for i in \$(seq 1 "\$RUNS"); do
        local seat_p1 seat_p2 my_seat out s1 s2
        if [ \$((i % 2)) -eq 1 ]; then
            seat_p1=\$MINE; seat_p2="\$ROBOTS/\$opponent"; my_seat=p1
        else
            seat_p1="\$ROBOTS/\$opponent"; seat_p2=\$MINE; my_seat=p2
        fi
        out=\$("\$ENGINE" -q -f "maps/\$map" -p1 "\$seat_p1" -p2 "\$seat_p2" 2>&1)
        read -r s1 s2 <<<"\$(echo "\$out" | parse_scores)"
        s1=\${s1:-0}; s2=\${s2:-0}
        if [ "\$my_seat" = p1 ] && [ "\$s1" -gt "\$s2" ]; then wins=\$((wins+1)); fi
        if [ "\$my_seat" = p2 ] && [ "\$s2" -gt "\$s1" ]; then wins=\$((wins+1)); fi
    done
    printf "%-11s on %-6s : %d/%d\n" "\$opponent" "\$map" "\$wins" "\$RUNS"
}

echo "arch=\$ARCH, runs per matchup: \$RUNS"
echo "-------------------------------------------"
run_match wall_e map00
run_match h2_d2  map01
run_match bender map02
echo
echo "Bonus:"
run_match terminator map00 || true
SCRIPT
)"
