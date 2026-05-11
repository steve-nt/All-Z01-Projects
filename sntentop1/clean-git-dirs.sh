#!/bin/bash

# Base directory to search for .git folders
base_dir="$HOME/Desktop/test1"

# Check if the directory exists
if [ ! -d "$base_dir" ]; then
    echo "Error: Directory $base_dir does not exist."
    exit 1
fi

# Dry-run: Show all .git folders that would be deleted (excluding $base_dir/.git)
echo "=== DRY RUN (No files will be deleted yet) ==="
find "$base_dir" -type d -name ".git" | while read git_dir; do
    # Skip if .git is directly in $base_dir (not a subdirectory)
    if [ "$git_dir" != "$base_dir/.git" ]; then
        echo "[Would delete] $git_dir"
    fi
done

# Ask for confirmation before actual deletion
read -p "Do you want to delete these .git folders? (Y/N) " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Aborted. No files were deleted."
    exit 0
fi

# Actual deletion (skipping $base_dir/.git)
echo "=== DELETING .git FOLDERS (except $base_dir/.git) ==="
find "$base_dir" -type d -name ".git" | while read git_dir; do
    if [ "$git_dir" != "$base_dir/.git" ]; then
        echo "Deleting: $git_dir"
        rm -rf "$git_dir"
    else
        echo "Skipping (protected): $git_dir"
    fi
done

echo "Cleanup completed. All .git folders (except $base_dir/.git) removed."
