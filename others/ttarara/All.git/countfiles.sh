#!/bin/bash

# Count the number of regular files and folders in the current directory and its subfolders
count=$(find . -type f -o -type d | wc -l )

# Print the count
echo "$count"

