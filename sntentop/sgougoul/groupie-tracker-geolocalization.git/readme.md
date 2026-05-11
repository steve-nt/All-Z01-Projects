# Groupie Tracker

## Overview
Groupie Tracker is a web application that allows users to track music artists, their concert locations, and relations between different artists. It fetches data from an API and presents it in a structured format for easy access.

## Features
- View a list of artists and their details
- Retrieve artist locations and concert dates
- View relations between artists
- User-friendly web interface
- Proper error handling and API integration
- A user friendly filters menu for a more advanced search based on matching criteria such as a band's creation date, first album year , locations of concerts etc.

## Installation
### Prerequisites
- Go (1.18 or later)
- Git

### Setup Instructions
1. Clone the repository:
   ```sh
   git clone https://platform.zone01.gr/git/sgougoul/groupie-tracker-filters
   cd groupie-tracker-filters
   ```
2. Run the command:
   ```sh
   chmod +x build.sh
   ```
3. Run the command:
   ```sh
   ./build.sh
   ```   
4. Export the api key via this command:
   ```sh
   export MAPQUEST_KEY=$(<secret.key)
   ```
5. Run the command:
   ```sh
   ./myapp
   ```        
6. The website is hosted on this adress :
   ```
   http://localhost:8080/
   ```

## API Endpoints
| Endpoint        | Method | Description |
|----------------|--------|-------------|
| `/`            | GET    | Serves the home page |
| `/artists`     | GET    | Fetches all artists |


## UI Route Endpoints (with Query Parameters)

These are the routes used in the web interface where query parameters (`id`) are passed to fetch data for a specific artist:

| Route                   | Description                                  | Example URL                              |
|-------------------------|----------------------------------------------|------------------------------------------|
| `/locations.html`       | Fetches and displays locations for a specific artist based on the `id` query parameter. | `http://localhost:8080/locations.html?id=1` |
| `/relations.html`       | Fetches and displays relations for a specific artist based on the `id` query parameter. | `http://localhost:8080/relations.html?id=1` |
| `/dates.html`           | Fetches and displays concert dates for a specific artist based on the `id` query parameter. | `http://localhost:8080/dates.html?id=1`   |

## Error Handling
The application includes a structured error handling mechanism to return user-friendly error messages for issues such as invalid requests, missing data, or internal server errors.


## ⚙️ Prerequisites

* Go 1.16+ installed.
* Bash (for the `build.sh` script).
* A `secret.key` file in the project root.

---


## Authors
- **sgougoul**
- **saggelak**
- **fffffff**