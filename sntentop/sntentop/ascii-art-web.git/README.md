# 🎨 **ASCII Art Web** 🎉

This README provides a comprehensive guide to the ASCII Art Web application, a dynamic tool for transforming text into captivating ASCII art.

### Description

ASCII Art Web empowers users to convert their text into visually engaging ASCII art using three distinct font styles: **Standard**, **Shadow**, and **Thinkertoy**. Built with Go, this application offers:

-   A responsive user interface for seamless interaction.
-   Custom error handling pages for HTTP status codes (`400`, `404`, `500`) to enhance user experience.
-   Comprehensive unit tests to ensure robust functionality and reliability.

Users can effortlessly input their desired text, select their preferred font, and instantly witness the generation of captivating ASCII art. In the event of an error, informative messages accompanied by engaging visuals are displayed.

---

### Authors

-   **Theocharoula Tarara** (*ttarara*)
-   **Christoforos Kissas** (*xkissas*)
-   **Stefanos Ntentopoulos** (*sntentop*)

---

### Usage: How to Run

#### Prerequisites

1.  Ensure you have **Go 1.19+** installed on your system.
2.  Clone this repository:

    ```bash
    git clone <repository-url>
    cd ascii-art-web
    ```

3.  Initialize a `go.mod` file in the root directory if it's missing:

    ```bash
    go mod init ascii-art-web
    ```

#### Steps to Run

1.  Navigate to the project root directory:

    ```bash
    cd ascii-art-web
    ```

2.  Run the application:

    ```bash
    go run .
    ```

3.  Open your browser and visit:

    ```
    http://localhost:8080
    ```

4.  To stop the server, press CTRL+C. 🛑

### Implementation Details: Algorithm

**Overview**

The project incorporates several key features:

-   **Input Validation:** Ensures the input adheres to the following constraints:
    -   Maximum length of 128 characters.
    -   Only printable ASCII characters (32–126) are permitted.
-   **ASCII Art Generation:**
    -   Reads font files from `static/fonts/`.
    -   Maps each character of the user-provided input to its corresponding ASCII art representation in the selected font.
-   **Error Handling:**
    -   Displays custom error pages for HTTP status codes (400, 404, 500) with relevant error messages and visuals.
-   **Testing:**
    -   Includes unit tests for:
        -   Template initialization.
        -   Handlers.
        -   Input validation.

**Algorithm**

-   **Font File Parsing:**
    -   Each font file (e.g., `standard.txt`, `shadow.txt`) contains 8-line ASCII art representations for printable characters.
    -   Character offset is calculated using: `(ASCII value - 32) * 9`
-   **ASCII Art Construction:**
    -   For each input character:
        -   Validate its ASCII range (32–126).
        -   Retrieve and append the corresponding ASCII art lines.
-   **Input Validation:**
    -   Checks if the input contains unsupported characters or exceeds 128 characters.
    -   Returns appropriate error codes and messages for invalid inputs.
-   **Error Handling:**
    -   Serves custom error pages (`400.html`, `404.html`, `500.html`) with visuals and error messages.

### Directory Structure
```bash
ascii-art-web/
├── handlers/
│   ├── handlers.go          # Contains HTTP handler functions and core logic
│   ├── handlers_test.go     # Unit tests for handler functions
├── templates/
│   ├── 400.html            # 400 - Bad Request page
│   ├── 404.html            # 404 - Page Not Found
│   ├── 500.html            # 500 - Internal Server Error
│   ├── index.html          # Main HTML page
├── static/
│   ├── style.css           # CSS file for styling
│   ├── images/
│   │   ├── 400.jpg         # 400 error image
│   │   ├── 404.png         # 404 error image
│   │   ├── error500.png    # 500 error image
│   ├── fonts/
│   │   ├── standard.txt     # Standard font
│   │   ├── shadow.txt       # Shadow font
│   │   ├── thinkertoy.txt   # Thinkertoy font
├── main.go                 # Entry point of the application
```

### Testing

#### How to Run Tests

1.  Navigate to the project root directory:

    ```bash
    cd ascii-art-web
    ```

2.  Run the tests:

    ```bash
    go test ./handlers
    ```

#### Expected Output
 ```bash
 ok      handlers      0.200s
 ```

#### Test Coverage

The following tests are implemented in `handlers_test.go`:

-   **404 - Page Not Found:** Ensures invalid routes return `404.html`.
-   **400 - Bad Request:** Verifies input validation and invalid font selection.
-   **500 - Internal Server Error:** Simulates a missing font file to ensure `500.html` is displayed.
-   **200 - OK:** Validates successful ASCII art generation for valid inputs.

### Features

-   **Dynamic ASCII Art:** Converts text into artistic representations using selected fonts.
-   **Responsive Error Pages:** Displays styled error messages for common HTTP status codes.
-   **Robust Testing:** Comprehensive test coverage ensures the application runs smoothly.

### License

This project is licensed under the MIT License. Feel free to use, modify, and distribute it.