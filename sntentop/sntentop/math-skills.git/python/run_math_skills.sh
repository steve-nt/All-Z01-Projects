#!/bin/bash

# Download the zip file
wget https://assets.01-edu.org/stats-projects/stat-bin-dockerized.zip
echo
sleep 5

# List files in the current directory
ls -la
echo
sleep 5

# Unzip the downloaded file
unzip stat-bin-dockerized.zip
echo
sleep 5

# List files again to see the contents of the unzipped folder
ls -la
echo
sleep 5

# Copy mathskills.py to the stat-bin directory
cp mathskills.py ./stat-bin/

# Copy commands-to-run.sh to the stat-bin directory
cp commands-to-run.sh ./stat-bin/

# Change to stat-bin directory
cd stat-bin/
echo
sleep 5

# List files in the stat-bin directory after copying
ls -la
echo
sleep 5

# Run the commands-to-run.sh script
bash commands-to-run.sh

# Move back to the parent directory
cd ..

# Ask user if they want to delete all files and folders starting with stat-bin
read -p "Do you want to remove all files and folders starting with 'stat-bin'? (yes/no): " choice
case "$choice" in
    yes|Yes|Y|y)
        rm -rf stat-bin*
        echo "Cleanup complete: Removed all files and folders starting with 'stat-bin'."
        ;;
    no|No|N|n)
        echo "Cleanup skipped."
        ;;
    *)
        echo "Invalid input. Cleanup skipped."
        ;;
esac