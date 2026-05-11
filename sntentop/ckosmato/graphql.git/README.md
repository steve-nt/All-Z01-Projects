# GraphQL Profile Viewer

A full-stack web application that displays user statistics and data visualizations by querying a GraphQL API. Built with Go backend and vanilla JavaScript frontend.

## Features

- **JWT Authentication**: Secure login system with token-based authentication
- **GraphQL Proxy**: Custom Go proxy for querying the GraphQL API
- **Data Visualization**: Interactive graphs and charts using D3.js
- **User Profile**: Display comprehensive user statistics and progress
- **Responsive Design**: Modern CSS styling for all screen sizes

## Tech Stack

### Backend
- **Go** - Server and API proxy
- **JWT** - Authentication tokens
- **GraphQL** - Data querying

### Frontend
- **HTML/CSS** - UI structure and styling
- **Vanilla JavaScript** - Client-side logic
- **D3.js** - Data visualization library

## Project Structure

```
graphql/
├── server/
│   ├── main.go           # Server entry point and routing
│   ├── go.mod            # Go dependencies
│   ├── handlers/
│   │   └── login.go      # JWT authentication handler
│   └── proxy/
│       └── proxy.go      # GraphQL proxy implementation
├── templates/
│   ├── index.html        # Login page
│   └── profile.html      # Profile dashboard
├── static/
│   └── style.css         # Application styles
├── scripts/
│   ├── main.js           # Main application logic
│   ├── profile.js        # Profile data fetching
│   ├── DrawGraphs.js     # D3.js visualization
│   └── utils.js          # Utility functions
└── README.md
```

## Installation

### Prerequisites
- Go 1.16 or higher
- A modern web browser

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd graphql
```

2. Install Go dependencies:
```bash
cd server
go mod download
```

3. Run the server:
```bash
go run main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

1. **Login**: Enter your credentials on the login page
2. **View Profile**: After authentication, you'll see your user dashboard
3. **Explore Data**: Interactive graphs display your progress and statistics
4. **Logout**: Click logout to end your session

## GraphQL API

The application proxies requests to a GraphQL endpoint to fetch:
- User information
- Project statistics
- Progress data
- XP and audit ratios

## Development

### Running locally
```bash
cd server
go run main.go
```

### Building for production
```bash
cd server
go build -o graphql-server
./graphql-server
```

## License

This project is part of the Athens Zone 01 curriculum.
