
#!/bin/bash

# Look for all files ending with .sh in the current directory and its subfolders,
# remove the .sh extension, and display the filenames in descending order



find . -type f -name "*.sh" | sed 's/\.sh$//' | sed 's|.*/||' | sort -r 


