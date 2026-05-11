# Ascii-art-web-stylize

This project implements a web server that handles requests for generating ASCII art from text input. It includes custom error handling for common HTTP status codes (404, 400, 500) and serves both GET and POST requests. The server validates ASCII-only input and processes the text accordingly.

## Features

- **Text Input**: Converts a given string into a visual representation using ASCII characters.

## Installation

**Clone the repository:**
```bash
git clone https://platform.zone01.gr/git/sbusho/ascii-art-web-stylize
```
**Navigate to the project directory:**
```bash
cd ascii-art-web-stylize
```
**Run the program:**
```bash
go run . 
```

## Usage

### Change Ascii Font
```bash
Switch between banners through given options, then press Generate ASCII Art button
```
**Implementation details (Algorithm)**:
- Request Handling:

The server listens for GET and POST requests.
For GET requests to the root path (/), the server returns an HTML page.
For POST requests to the root path, it processes the text form field to generate ASCII art.
- Error Handling:

Custom 404 handler: If a non-existent path is requested, the server returns a 404 error with a user-friendly message.
Custom 400 handler: If the text parameter is empty or invalid, the server returns a 400 error.
Custom 500 handler: If an internal error occurs, the server returns a 500 error.
- Validation:

The server validates that the input text contains only ASCII characters (0-127).
If the input contains non-ASCII characters, a 400 error is returned.
- Graceful Shutdown:

The server is designed to handle graceful shutdown on receiving termination signals (e.g., Ctrl+C).
It waits for ongoing requests to finish before stopping.
- Testing:

Unit tests are included to verify that the server correctly handles POST requests, including edge cases like empty input.
The server is tested using the httptest package to simulate HTTP requests and check responses.


### Supported Fonts
#### Standard
```bash
       _                         _                      _  
      | |                       | |                    | | 
 ___  | |_    __ _   _ __     __| |   __ _   _ __    __| | 
/ __| | __|  / _` | | '_ \   / _` |  / _` | | '__|  / _` | 
\__ \ \ |_  | (_| | | | | | | (_| | | (_| | | |    | (_| | 
|___/  \__|  \__,_| |_| |_|  \__,_|  \__,_| |_|     \__,_| 
```
#### Shadow
```bash
         _|                      _|
  _|_|_| _|_|_|     _|_|_|   _|_|_|   _|_|   _|      _|      _| 
_|_|     _|    _| _|    _| _|    _| _|    _| _|      _|      _| 
    _|_| _|    _| _|    _| _|    _| _|    _|   _|  _|  _|  _|   
_|_|_|   _|    _|   _|_|_|   _|_|_|   _|_|       _|      _|     
```
#### Thinkertoy
```bash
 o  o           o             o
 |  |    o      | /           |
-o- O--o   o-o  OO   o-o o-o -o- o-o o  o 
 |  |  | | |  | | \  |-' |    |  | | |  | 
 o  o  o | o  o o  o o-o o    o  o-o o--O 
                                        | 
                                     o--o 
```
#### Other
```bash                                                    
           **     **                         
  ****   ******** ******     ****   **  **** 
**    **   **     **    ** ******** ****     
**    **   **     **    ** **       **       
  ****       **** **    **   ****** **
  ```
## Authors

Thanos Ziagakis
<!-- 🍉              -->
Sofia Busho
<!-- 🌸    -->
Maria Tzemanaki
<!-- 🍓           -->