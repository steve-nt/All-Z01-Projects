#!/bin/bash

# Prompt the user to enter the number of times to run the commands
read -p "How many times would you like to run the commands? " count

# Loop to run the commands 'count' number of times
for ((i=1; i<=count; i++)) # Loop from 1 to the specified count
do 
    echo "Running iteration $i"
    echo "------------------------------------------"
    echo
    echo "/bin/math-skills numeric results"
  
    ./bin/math-skills
    echo "------------------------------------------"
    echo "main.go numeric results"
  
    go run main.go data.txt
    echo "=========================================="
    echo
    
    # 5-second break before the next iteration
    echo "Pausing for 5 seconds before the next iteration..."
    sleep 5
    echo
done