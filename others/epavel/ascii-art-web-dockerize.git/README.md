# ASCII Art Web Dockerize

## Description

ASCII Art Web is a web-based application that generates ASCII art from user input text using different banner styles. The application provides a graphical user interface (GUI) to input text, select a banner style, and generate the corresponding ASCII art.

## Authors

- Pavlos Kerasidis
- Kimon Dimitrioy
- Edouardos Pavel

## Docker
Building your project:
 - docker build -t ascii-art-web-dockerize:latest .
Showing images:
 - docker images
Run the program:
 - docker run -p 8080:8080 ascii-art-web-dockerize:latest
Remove the project
 - docker rm ascii-art-web-dockerize

## Usage

To run the server, use the following command:

```sh
docker build -t ascii-art-web-dockerize:latest .   //Building the server
docker run -p 8080:8080 ascii-art-web-dockerize:latest    //Run the server
Then open http://localhost:8080 in your web browser to access the ASCII art generator.

Implementation Details
The application uses Go\'s net/http package to create a web server and html/template package to render HTML templates. The ASCII art generation logic is implemented in the asciiart package, which reads banner files and converts input text to ASCII art.

HTTP Endpoints
GET /: Sends the main HTML page.
POST /ascii-art: Receives text and banner data, generates ASCII art, and displays the result.
HTML Templates
templates/index.html: Main page with text input, banner selection, and submit button.
templates/result.html: Page to display the generated ASCII art.
Terminal Size
The terminal size is obtained using the syscall package from the standard library. The GetTerminalWidth function retrieves the terminal width to ensure proper formatting of the ASCII art.

Instructions
Clone the repository.
Navigate to the project directory.
Run the server using go run main.go.
Open http://localhost:8080 in your web browser.
Enter the text and select a banner style.
Click the "Generate" button to see the ASCII art.
Allowed Packages
Only the standard Go packages are allowed.

Example
Here\'s an example of how to use the application:

Start the server:

Open your web browser and go to http://localhost:8080.

Enter the text "Hello" and select the "standard" banner.

Click the "Generate" button to see the ASCII art.

HTTP Status Codes
200 OK: If everything went without errors.
404 Not Found: If nothing is found, for example, templates or banners.
400 Bad Request: For incorrect requests.
500 Internal Server Error: For unhandled errors.
```
## Implementation Details
- 1.Initialize Variables and Constants

    - Define a slice BANNERS containing banner styles.
    - Define a struct PageData to hold data for rendering the HTML      template.

- 2.Main Function

    - Initialize a boolean flag flag to track template rendering errors.
    - Set up an HTTP handler for the root path /.

- 3.HTTP Handler Function

    - Check URL Path
        - If the URL path is not /, respond with a 404 Not Found.
    - Check Error Flag
        - If flag is true, render an error page and return.
    - Initialize Context
        - Declare a variable context of type PageData.
    - Handle GET Request
        -If the request method is GET, initialize context with default values.
    - Handle POST Request
        - If the request method is POST:
            - Parse the form data.
            - If parsing fails, respond with a 400 Bad Request and log the error.
            - Retrieve selectedBanner and userInput from the form data.
            - Call utils.Ascii_art to generate ASCII art from userInput and selectedBanner.
            - If ASCII art generation fails, set asciiOutput to the error message.
            - Update context with the form data and ASCII art output.
    - Handle Unsupported Methods
        - If the request method is not GET or POST, respond with a 405 Method Not Allowed.
    - Load and Render Template
        - Parse the HTML template file home.html.
        - If parsing fails, respond with a 500 Internal Server Error and log the error.
        - Execute the template with context.
        - If rendering fails, set flag to true and log the error.
- 4.Start HTTP Server

    - Log a message indicating the server is starting.
    - Start the HTTP server on port 8080.
    - If the server fails to start, log the error and terminate the program.
