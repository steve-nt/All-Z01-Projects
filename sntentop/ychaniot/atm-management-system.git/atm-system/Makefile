CC = gcc
CFLAGS = -Wall -g
SRC_DIR = src
OBJ_DIR = obj
BIN = atm

# Source files
SRC = $(wildcard $(SRC_DIR)/*.c)
# Object files
OBJ = $(SRC:$(SRC_DIR)/%.c=$(OBJ_DIR)/%.o)

# Default target
all: directories $(BIN)

# Create output directories
directories:
	@mkdir -p $(OBJ_DIR)
	@mkdir -p data

# Link object files to create executable
$(BIN): $(OBJ)
	$(CC) $(CFLAGS) $^ -o $@

# Compile source files to object files
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c
	$(CC) $(CFLAGS) -c $< -o $@

# Clean build files
clean:
	rm -rf $(OBJ_DIR) $(BIN)

# Run the program
run: all
	./$(BIN)

.PHONY: all clean run directories