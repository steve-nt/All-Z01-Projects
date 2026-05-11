# Push-Swap

**Push-Swap** is a Go-based project that implements a non-comparative sorting algorithm using two stacks. The project is split into two main components:
1. **Push-Swap Program**: Generates a sequence of operations to sort the stack (located in the `push-swap` folder).
2. **Checker Program**: Verifies if a given sequence of operations correctly sorts the stack (located in the `checker` folder).

## Features

- **Two Stacks**: Uses two stacks (`stackA` and `stackB`) for sorting.
- **Operations**: Supports a set of operations (`push`, `swap`, `rotate`, `reverse rotate`) on the stacks.
- **Optimized Sorting**: Minimizes the number of operations required to sort the stack.
- **Error Handling**: Handles invalid inputs, duplicates, and incorrect sequences.
- **Verification**: Verifies if the sequence of operations sorts the stack correctly.

## Project Structure

- **`push-swap/main.go`**: Main program for the Push-Swap algorithm, which parses the input and generates the sequence of operations to sort the stack.
- **`checker/main.go`**: Program that verifies if the sequence of operations correctly sorts the stack.
- **`utils/parse.go`**: Contains functions for parsing the input and initializing the stack.
- **`utils/sort.go`**: Contains the sorting logic, including functions for sorting small stacks and handling the main sorting process.
- **`operations/operations.go`**: Contains the implementation of the operations (`pa`, `pb`, `sa`, `sb`, `ra`, `rb`, `rra`, `rrb`, `ss`, `rr`, `rrr`).
- **`operations/operations_test.go`**: Contains unit tests for verifying the functionality of the operations.
- **`utils/sort_test.go`**: Contains unit tests for verifying the sorting logic.
- **`go.mod`**: Go module definition for dependency management.

**Clone the repository**:

```bash
git clone https://platform.zone01.gr/git/sbusho/push-swap
```

**Navigate to the project directory:**
```bash
cd push-swap
```

## How to Run

### For the Push-Swap Program
1. Run the Push-Swap program with a list of integers:
   go run ./push-swap "2 1 3 6 5 8"
This should output the sequence of operations required to sort the input list. For example, it should display a valid solution and less than 9 instructions.

### For the Checker Program
1. Run the Checker program with the sequence of operations:
You can provide input directly through the command line:
echo -e "sa\npb\nrrr\n" | go run ./checker "0 9 1 8 2 7 3 6 4 5"
The Checker will output "OK" if the sequence is correct or "KO" if it's incorrect.

### Error Handling
Invalid Inputs: If the input list contains non-integer values (e.g., "0 one 2 3"), the program will display an error message.
Example:
go run ./push-swap "0 one 2 3"
Error

Duplicate Numbers: If the input list contains duplicate numbers (e.g., "1 2 2 3"), the program will display an error message.
Example:
go run ./push-swap "1 2 2 3"
Error

### Validity of Solutions
The Push-Swap program should output a valid solution with the minimum number of operations for sorting the stack. For example:
go run ./push-swap "5 4 3 2 1"
sa
rra

For a sorted input list like "0 1 2 3 4 5", the program should not output any instructions.
Example:
go run ./push-swap "0 1 2 3 4 5"
Output: nothing (no operations needed).

### Example Command Sequences
Testing with Random Numbers:
You can run the program with 5 random numbers (e.g., "4 2 5 3 1") and it should display a valid solution and fewer than 12 instructions.
Example:
go run ./push-swap "4 2 5 3 1"
Output: a valid sequence with fewer than 12 operations.

Testing with Invalid Input:
If you provide invalid input, like non-integer values or duplicates, the program will display an error message.
Example:
go run ./push-swap "1 2 3 a"
Error

### Unit Testing
The project includes unit tests for the sorting and operation functions.
```bash
cd operations
``` 
or
```bash
cd utils
```
Run tests with:
```bash
go test
```
or 
```bash
go test -v
```

### License
This project is licensed under the MIT License.


## Authors

🌺 Iana Kopylova 🌷

🌷 Sofia Busho 🌺