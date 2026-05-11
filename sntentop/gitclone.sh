#!/bin/bash

# --- CONFIGURATION ---
GITEA_URL="https://platform.zone01.gr/git"    # No trailing slash!
GITEA_USER="sntentop"                         # Gitea username
GITEA_PASS="1qaz!QAZ"                         # Gitea password
REQUIRED_SPACE_MB=500                         # Minimum space (MB)
DRY_RUN=false                                 # "true"=dry run, "false"=clone

prompt_for_credentials() {
    if [ -z "$GITEA_USER" ]; then
        read -p "Enter Gitea username: " GITEA_USER
    fi
    if [ -z "$GITEA_PASS" ]; then
        read -s -p "Enter Gitea password: " GITEA_PASS
        echo
    fi
}

check_api_credentials() {
    while true; do
        test_response=$(curl -su "$GITEA_USER:$GITEA_PASS" -f "$API_URL&page=1" 2>/dev/null)
        if [ -z "$test_response" ]; then
            echo "Could not fetch repositories. Please check your credentials and try again."
            GITEA_USER=; GITEA_PASS=
            prompt_for_credentials
        else
            if echo "$test_response" | grep -q '"message":.*"invalid credentials"'; then
                echo "Authentication failed (invalid credentials). Please try again."
                GITEA_USER=; GITEA_PASS=
                prompt_for_credentials
            else
                break
            fi
        fi
    done
}

# --- MAIN SCRIPT ---
API_URL="$GITEA_URL/api/v1/user/repos?limit=1000&sort=updated"
prompt_for_credentials
check_api_credentials

# --- Fetch all repository clone URLs with pagination ---
repo_urls=""
page=1
while : ; do
    response=$(curl -su "$GITEA_USER:$GITEA_PASS" -f "$API_URL&page=$page" 2>/dev/null)
    count=$(echo "$response" | jq 'length')
    if (( count == 0 )); then
        break
    fi
    urls=$(echo "$response" | jq -r '.[].clone_url')
    repo_urls="$repo_urls"$'\n'"$urls"
    ((page++))
done

repo_urls=$(echo "$repo_urls" | grep .)   # Remove empty lines

# --- Clone each repository ---
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

# Commented-out user prompt for batch/unattended cloning:    
#    while true; do
#        read -p "Proceed with cloning $repo_name? (y = yes / n = skip / q = quit): " answer
#        case $answer in
#            [Yy]* ) break ;;
#            [Nn]* ) echo "Skipping $repo_name"; continue 2 ;;
#            [Qq]* ) echo "Quitting."; exit 0 ;;
#            * ) echo "Please answer y, n, or q." ;;
#        esac
#    done

    if [ "$DRY_RUN" = true ]; then
        echo "[DRY RUN] Would clone $url into $repo_name"
    else
        # Authentication handling for git clone
        while true; do
            GIT_ASKPASS_SCRIPT=$(mktemp)
            chmod +x "$GIT_ASKPASS_SCRIPT"
            cat > "$GIT_ASKPASS_SCRIPT" <<EOF
#!/bin/sh
case "\$1" in
    Username*) echo "$GITEA_USER" ;;
    Password*) echo "$GITEA_PASS" ;;
esac
EOF
            GIT_ASKPASS="$GIT_ASKPASS_SCRIPT" GIT_TERMINAL_PROMPT=0 git clone "$url" "$repo_name" 2>git_clone_err.log
            result=$?
            rm -f "$GIT_ASKPASS_SCRIPT"

            if [ $result -eq 0 ]; then
                echo "Cloned successfully."
                break
            elif grep -qi "authentication failed" git_clone_err.log || grep -qi "failed to authenticate" git_clone_err.log; then
                echo "Authentication failed for repo: $repo_name"
                GITEA_USER=; GITEA_PASS=
                prompt_for_credentials
            else
                echo "Cloning failed for $repo_name. Check git_clone_err.log for details."
                break
            fi
        done
        rm -f git_clone_err.log
    fi
done