README.md

# ASCII Art Justify!

This program generates ASCII art from a given string and banner type.

## Usage

To run the program, compile and execute the main.go file. The program accepts the following arguments:

* `--align`: Specifies the text alignment. Valid values are `left`, `center`, `right`, and `justify`. Defaults to `left` if not specified.
* `string`: The string to be converted to ASCII art.
* `banner`: The type of banner to use. Valid values are `standard`, `shadow`, and `thinkertoy`. Defaults to `standard` if not specified.

**Example Usage:**

```bash
go run . --align=right "Hello, world!" shadow


This command will generate ASCII art for the string "Hello, world!" using the shadow banner and right alignment.

Testing
The program includes a test file (01_test.go) with comprehensive tests for various functions. To run the tests, use the following command:

Bash
go test
go test -v (For more details on each test)

This command will execute all the tests and report the results.

Additional Notes
The program reads banner styles from text files located in the "banner" directory.
The program calculates the width of the terminal window to ensure proper formatting of the output.
The program includes error handling for invalid arguments, alignment values, and non-ASCII characters in the input string.
The program supports different alignment options, including left, center, right, and justify.

Contributing
Feel free to contribute to the project by submitting bug reports, feature requests, or pull requests.


