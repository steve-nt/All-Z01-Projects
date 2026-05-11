# Compiler and flags
CC = gcc
CFLAGS = -Wall -g

# Source files
SRC = src/main.c src/system.c src/auth.c
OBJ = $(SRC:.c=.o)  # Convert .c files to .o files
EXEC = atm  # The name of the output executable

# Rule to build the executable
$(EXEC): $(OBJ)
	$(CC) $(OBJ) -o $(EXEC)

# Rule to compile .c files into .o files
%.o: %.c src/header.h  # All .o files depend on their corresponding .c files and the header
	$(CC) $(CFLAGS) -c $< -o $@

# Clean rule to remove compiled files
clean:
	rm -f $(OBJ) $(EXEC)
