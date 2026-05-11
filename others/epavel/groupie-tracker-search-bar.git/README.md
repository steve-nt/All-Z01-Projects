# Groupie Tracker

Groupie Tracker is a web application designed to track and display information about various artists. This project leverages several modern technologies to ensure a seamless and efficient user experience.

## Deployed Project

### [Here](https://groupie-tracker.gr)

## Authors

- [Giannis Georgakopoulos](https://platform.zone01.gr/git/ggeorgako)
- [Giorgos Pavrianidis](https://platform.zone01.gr/git/gpavrian)
- [Edouardos Pavel](https://platform.zone01.gr/git/epavel)

## Technologies Used

- **Go**: The core programming language used for the backend.
- **Proprietary Cache**: Custom made cache used for storing frequently accessed data to improve performance.
- **Goroutines**: Utilized for concurrent processing to enhance performance.
- **HTML Templates**: Go's `html/template` package is used for rendering dynamic content.
- **Tailwind CSS**: A utility-first CSS framework for styling the frontend.
- **Docker**: Used for containerization and deployment, including DNS indexing.
- **Middlewares**: Custom middleware functions to handle cross-cutting concerns like logging, authentication, and request validation.
- **Google Maps API**: The Google Maps API is utilized to display the locations of concerts on an interactive map.

## Code Structure
> we know some(a lot) of comments in the code are abundant, but we use the comments to provide InteliSense with the description of any function we made. 
- **main.go**: The entry point of the application.
- **bin/**: Contains Go files for routing and rendering templates.
    - `routes.go`: Defines the routes and handles HTTP requests.
    - `render.go`: Contains functions for rendering templates and handling errors.
    - `fetch.go`: Handles all API calls and their respective responses.
    - `parse.go`: Formats the decoded JSON data to be used efficiently.
    - `cache.go`: Implements caching mechanisms to store frequently accessed data.
    - `filters.go`: Contains functions to filter and process data based on various criteria.
    - `global.go`: Defines global variables and configurations used throughout the application.
    - `handlers.go`: Contains HTTP handler functions for different routes.
    - `middleware.go`: Implements middleware functions for request processing.
    - `search.go`: Contains functions to handle search queries and return relevant results.
- **core/**:
    - `init.go`: Initializes the application, setting up configurations and dependencies.
    - `server.go`: Starts the HTTP server and listens for incoming requests.
    - `.txt files`: Used for initializing filters and for minimizing startup times
- **templates/**: Contains HTML templates.
    - `home.html`: The main template for the homepage.
    - `artist.html`: Diving into the details of every given artist.
    - `error.html`: Responsible for the handling of errors.
- **src/**: Contains source files for the application.
    - `input.css`: The main CSS file that includes Tailwind CSS directives.
- **public/**: Contains static assets like CSS files.
    - `app.js`: Contains JavaScript code for client-side interactions.
    - `maps.js`: Contains JavaScript code for Google Maps API interactive visualsizations also containing the exact history of the concerts themselves
    - `some icons and misc stuff`: For the small little details. :)
    - `output.css`: The compiled CSS file from Tailwind CSS.
- **Dockerfile**: Defines the Docker image for deployment.
- **tailwind.config.js**: Configuration file for customizing Tailwind CSS settings and extending the default theme.
- **package.json/-lock.json**: These files are used to manage the project's dependencies. `package.json` contains metadata about the project and lists the packages required for the project, while `package-lock.json` records the exact versions of the dependencies installed to ensure consistent installs across different environments.
- **output.css**: The compiled CSS file generated from Tailwind CSS directives in `input.css`. This file contains all the necessary styles for the application, ensuring a consistent and responsive design across different devices.

## Features

- **Concurrent Processing**: Goroutines are used to fetch and process artist data concurrently, improving the application's responsiveness.
- **Context**: The context Go package is utilized to pass data such as filtered and paginated artists, total artist count, messages, and errors between middleware and handlers using context keys. This allows for efficient and organized data handling throughout the request lifecycle.
- **Dynamic Content Rendering**: Go templates are used to dynamically render HTML content based on the data fetched from the backend.
- **Responsive Design**: Tailwind CSS ensures that the application is responsive and visually appealing across different devices.
- **Containerized Deployment**: The Dockerfile allows for easy deployment and scaling of the application, including DNS indexing for efficient routing.
- **Filters and Searching**: All actual processes happen exclusively in the back-end, taking advantage of the caching system for faster data filtering
- **Dynamic Caching**: The application employs a dynamic caching mechanism to enhance performance and reduce the load on the backend. Various caches are used to store frequently accessed data, such as artist information, locations, and filter data. The caching system includes an LRU (Least Recently Used) eviction policy to manage cache size and prevent memory leaks. Cached data is periodically refreshed to ensure that the application serves up-to-date information while minimizing API calls.

## Getting Started
> This applies for the development stage

### Prerequisites

- Go 1.23.2 or later
- Docker

### Installation

1. Clone the repository:
     ```sh
     git clone https://platform.zone01.gr/git/ggeorgako/groupie-tracker.git
     cd groupie-tracker
     ```

2. Build the Go application:
     ```sh
     go build
     ```

3. Run the application:
     ```sh
     ./groupie-tracker
     ```

### Docker Deployment

1. Build the Docker image:
     ```sh
     docker build -t groupie-tracker .
     ```

2. Run the Docker container:
     ```sh
     docker run -p 8080:8080 groupie-tracker
     ```

## Usage

Navigate to `http://localhost:8080` in your web browser to access the application. Use the search functionality to find and view information about different artists.
> Note: this applies for local testing, please refer to [Deployed Project](#deployed-project) for the actual project

### API Endpoints

The Groupie Tracker application provides several API endpoints to fetch data about artists, locations, dates, and relations. Below are the available endpoints:

- **GET /api/artists**: Retrieves a list of all artists.
- **GET /api/artists/{id}**: Retrieves detailed information about a specific artist by their ID.
- **GET /api/locations**: Retrieves a list of all locations.
- **GET /api/locations/{id}**: Retrieves detailed information about a specific location by its ID.
- **GET /api/dates**: Retrieves a list of all dates.
- **GET /api/dates/{id}**: Retrieves detailed information about a specific date by its ID.
- **GET /api/relations**: Retrieves a list of all relations.
- **GET /api/relations/{id}**: Retrieves detailed information about a specific relation by its ID.

These endpoints allow for comprehensive access to the application's data, enabling users to integrate and utilize the information in various ways.
