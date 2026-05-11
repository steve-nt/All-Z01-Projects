#!/bin/bash

# --- CONFIGURATION ---
GITEA_URL="https://platform.zone01.gr/git"
REQUIRED_SPACE_MB=500
DRY_RUN=false

# Set global Git credential helper to 'store'
git config --global credential.helper store

prompt_for_credentials() {
    if [ -z "$GITEA_USER" ]; then
        read -p "Enter Gitea username: " GITEA_USER
    fi
    if [ -z "$GITEA_PASS" ]; then
        read -s -p "Enter Gitea password: " GITEA_PASS
        echo
    fi
}

API_URL="$GITEA_URL/api/v1/user/repos?limit=1000&sort=updated"

prompt_for_credentials

# --- Fetch all repository HTTPS clone URLs with pagination ---
repo_urls=""
page=1
while : ; do
    response=$(curl -su "$GITEA_USER:$GITEA_PASS" -f "$API_URL&page=$page" 2>/dev/null)
    count=$(echo "$response" | jq 'length')
    if (( count == 0 )); then
        break
    fi
    # Clean up URLs (no trailing slash)
    urls=$(echo "$response" | jq -r '.[].clone_url' | sed 's,/*$,,')
    repo_urls="$repo_urls"$'\n'"$urls"
    ((page++))
done
repo_urls=$(echo "$repo_urls" | grep .)   # Remove empty lines

for url in $repo_urls; do
    base_repo_name=$(basename "$url" .git)
    repo_name="$base_repo_name"

    # Avoid duplicate directory names
    suffix=2
    while [ -d "$repo_name" ]; do
        repo_name="${base_repo_name}-$suffix"
        ((suffix++))
    done

    available_space=$(df -m . | awk 'NR==2 {print $4}')
    echo ""
    echo "Next repository: $repo_name"
    echo "Available space: ${available_space} MB"

    if (( available_space < REQUIRED_SPACE_MB )); then
        echo "WARNING: Less than ${REQUIRED_SPACE_MB} MB free space!"
    fi

    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] Would clone $url into $repo_name"
    else
        # Embed credentials directly into URL
        auth_url=$(echo "$url" | sed -E "s#https://#https://$GITEA_USER:$GITEA_PASS@#")

        # Try to clone and then fetch all branches/tags
        git clone "$auth_url" "$repo_name" 2>git_clone_err.log
        result=$?
        if [ $result -eq 0 ]; then
            (
                cd "$repo_name" && \
                git fetch --all && \
                git fetch --tags

                # --- Create local tracking branches for all remote branches ---
                for remote in $(git branch -r | grep 'origin/' | grep -v 'HEAD' | sed 's|origin/||'); do
                    if [ ! "$(git branch --list "$remote")" ]; then
                        git branch --track "$remote" "origin/$remote" 2>/dev/null
                    fi
                done
            )
            echo "Cloned successfully."
        elif grep -qi "authentication failed" git_clone_err.log || grep -qi "failed to authenticate" git_clone_err.log; then
            echo "Authentication failed for repo: $repo_name"
            GITEA_USER=; GITEA_PASS=
            prompt_for_credentials
            # Retry the current repo again with new credentials
            continue
        else
            echo "Cloning failed for $repo_name. Check git_clone_err.log for details."
        fi
        rm -f git_clone_err.log
    fi
done

# Remove credentials after cloning is complete for security
rm -f ~/.git-credentials