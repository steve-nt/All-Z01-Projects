# groupie-tracker

Groupie Tracker is a web application built in Go that fetches and displays information about various artists, their concert locations and dates. This project aims to provide an interactive and user-friendly experience using client-server communication and data visualizations.

## Features
- Fetch and display artist details (name, image, members, first album, etc.)
- Display concert locations and dates
- Interactive UI using HTML and CSS
- Handles client-server interactions via API requests
- Implements an event-driven system (e.g., user-triggered actions)
- Error handling to prevent crashes and ensure smooth user experience
- Unit tests for core functionalities


## Installation

**Clone the repository:**
```bash
git clone https://platform.zone01.gr/git/ttarara/groupie-tracker
```
**Navigate to the project directory:**
```bash
cd groupie-tracker
```
**Run the Server:**
```bash
go run . 
```
The server will start on http://localhost:8080/.

## Usage

1. Browse through different artists, their concert dates, and locations.

2. Click on specific elements (e.g., artist cards) to trigger actions like displaying detailed information.

3. Experience the event-driven system by interacting with UI elements.


### **Deployment**:
The project is deployed on cloud platforms:

-   Railway.

     https://groupie-tracker-production-60e6.up.railway.app/

- Render:

    https://test1-m7di.onrender.com/


**API Structure**:
The project fetches data from a given API, which consists of four endpoints:

- /artists → Fetches artist details
- /locations → Fetches artist locations
- /dates → Fetches concert dates
- /relation → Links the above data

Each artist has a unique ID, and data is displayed dynamically on the frontend.

**Event System**:
The project implements client-server interactions triggered by events, such as:

- Clicking an artist’s name fetches additional details dynamically.
- Searching for a band filters results in real-time.
- A form submission sends data to the server for processing.

Asynchronous Go routines and channels ensure smooth execution.

**Testing**:
Unit tests ensure the reliability of core functionalities.

Run tests with:
```bash
go test
```
or 
```bash
go test -v
```


## Authors
   
    Theocharoula Tarara 🎵

    🎵 Stefanos Ntentopoulos

    Sofia Busho 🎵


## 💃  Enjoy exploring the world of music with Groupie Tracker! 🕺