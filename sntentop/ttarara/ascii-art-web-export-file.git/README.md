ASCII Art Web Export
Description

ASCII Art Web Export is a Go-based web application that extends the functionality of the ASCII Art project. This application enables users to generate ASCII art and export the generated output as a file, specifically in .txt format. It adheres to HTTP standards by including appropriate headers such as Content-Type, Content-Length, and Content-Disposition for file transfer. The application provides a simple and user-friendly interface, including a download button, for seamless exporting of the generated ASCII art.
Features

    ASCII Art Export:
        Export the generated ASCII art output as a .txt file.

    Correct File Permissions:
        The exported file is created with read and write permissions for the user.

    HTTP Headers:
        The web server uses the required HTTP headers (Content-Type, Content-Length, and Content-Disposition) for file transfer.

    Error Handling:
        User-friendly error messages are displayed if:
            The ASCII art output is empty.
            The output file is missing or deleted.

    Standards and Best Practices:
        Built entirely with Go standard libraries, ensuring clean and maintainable code.

Requirements

    Go (version 1.16 or later)
    Browser to access the web application
    No additional Go packages are required (only standard packages are used)

How to Run

    Clone the Repository:

git clone https://github.com/yourusername/ascii-art-web-export.git
cd ascii-art-web-export

Run the Application:

    Start the Go server:

        go run main.go

    Access the Application:
        Open your browser and go to: http://localhost:8080/home

    Generate and Export ASCII Art:
        Navigate to the converter page.
        Generate ASCII art.
        Click the "Download" button to export the output as a .txt file.

How to Export ASCII Art

    Navigate to the Converter page.
    Enter the text to generate ASCII art.
    Click "Generate" to create the ASCII art output.
    Click "Download" to export the file as a .txt.

Error Handling

    Empty ASCII Art:
        Displays an alert: "You should generate your ASCII-ART before downloading."
    Missing File:
        Displays an alert: "The output.txt file is missing. Please generate your ASCII-ART again."

Project Structure

```bash
ASCII-ART-WEB-EXPORT/
├── handlers/
│   ├── handlers_download.go   # Handles file download functionality
│   ├── handlers_export.go     # Handles file export logic
│   ├── handlers_logger.go     # Handles logging
│   ├── handlers_static.go     # Serves static files
│   ├── handlers_utils.go      # Utility functions
│   └── handlers.go            # General handlers and routes
├── static/
│   ├── banners/                       # ASCII art banner files
│   │   ├── shadow.txt
│   │   ├── standard.txt
│   │   └── thinkertoy.txt
│   ├── css/                           # Styles and fonts
│   │   ├── fonts/
│   │   │   ├── Cybrpnuk2.ttf
│   │   │   ├── The_Globe.ttf
│   │   │   └── Vintage_Brother.ttf
│   │   └── style.css
│   └── img/                           # Images for the web UI
│       ├── 400.png
│       ├── 404.png
│       ├── 500.png
│       ├── snetop.jpg
│       ├── ttarara.jpg
│       └── xkissas.jpg
├── temp/                              # Temporary files directory
│   └── output.txt                     # ASCII art output file
├── templates/                         # HTML templates
│   ├── code_400.html
│   ├── code_404.html
│   ├── code_500.html
│   ├── Index.html
│   ├── PgConverter.html
│   ├── PgHome.html
│   ├── PgProject.html
│   └── PgTeam.html
├── go.mod                             # Go module file
├── main.go                            # Application entry point
└── README.md                          # Project documentation
```

Best Practices Followed

    Code Modularity:
        Handlers are organized into a dedicated handlers package for clarity and maintainability.

    Standard Go Packages:
        The project uses only Go standard packages to meet the project requirements.

    HTTP Standards:
        The application adheres to HTTP standards for headers and file transfers.

Lessons Learned

Examples of API Interaction
Example 1: Downloading a .txt File

    Endpoint: /download
    Method: GET
    Client Interaction:
        The browser or JavaScript triggers the download.
    Server Response:
        The file is served as an attachment with appropriate headers:

Content-Type: text/plain
Content-Disposition: attachment; filename="output.txt"

This project helped deepen understanding of:

    File Export Mechanisms:
        Generating and exporting files in .txt format with appropriate permissions.

    HTTP Headers:
        Implementing Content-Type, Content-Length, and Content-Disposition.

    Web Development in Go:
        Building a web application using only standard Go packages.

    Error Handling:
        Handling user interactions and providing clear feedback for errors.

License

This project is licensed under the MIT License.

Feel free to fork, modify, and contribute to the project!