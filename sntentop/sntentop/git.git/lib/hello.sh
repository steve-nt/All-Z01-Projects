#!/bin/bash

# Default is World
# Author: Jim Weirich <jim.weirich@gmail.com>
<<<<<<< HEAD
source lib/greeter.sh

=======
<<<<<<< HEAD
#name=${1:-"World"}
#
#echo "Hello, $name"

echo "What's your name"
read my_name

echo "Hello, $my_name"
=======
source lib/greeter.sh

>>>>>>> main
name="$1"
if [ -z "$name" ]; then
    name="World"
fi

Greeter "$name"
<<<<<<< HEAD
=======
>>>>>>> greet
>>>>>>> main
