#!/bin/bash

DRY_RUN=false  # Set to false to actually move

for directory in */; do
    neodir="${directory%/}"

    # Check it's a Git repo
    if [ ! -d "$neodir/.git" ]; then
        continue
    fi

    # Get the origin URL
    remote_url=$(git -C "$neodir" remote get-url origin 2>/dev/null)
    if [ -z "$remote_url" ]; then
        echo "No origin remote found in $neodir, skipping."
        continue
    fi

    # Only process Gitea-style URLs
    echo "$remote_url" | grep -q 'platform.zone01.gr/git' || continue

    # Extract owner and repo
    owner=$(echo "$remote_url" | sed -E 's#.*/git/([^/]+)/[^/]+(\.git)?$#\1#')
    repo=$(echo "$remote_url" | sed -E 's#.*/git/[^/]+/([^/]+)(\.git)?$#\1#')

    if [[ -z "$owner" || -z "$repo" ]]; then
        echo "Could not extract owner/repo from $remote_url"
        continue
    fi

    # Determine unique destination path with suffix if needed
    dest_dir="${owner}/${repo}"
    suffix=2
    while [ -e "$dest_dir" ]; do
        dest_dir="${owner}/${repo}-${suffix}"
        ((suffix++))
    done

    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] Would move '$neodir' → '$dest_dir'"
    else
        mkdir -p "$owner"
        mv "$neodir" "$dest_dir"
        echo "Moved '$neodir' → '$dest_dir'"
    fi
done