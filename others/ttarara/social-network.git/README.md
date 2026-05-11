# Social Network

## Project Description

This project consists of building a **Facebook-like social network** that allows users to interact, share content, communicate in real time, and organize themselves in groups. The application implements core social networking features such as **profiles, followers, posts, groups, chats, and notifications**.

The goal of the project is to practice full‑stack development by combining **frontend technologies**, a **Go backend**, **SQLite database management**, **WebSocket communication**, and **Docker containerization**.

The system follows a **client–server architecture** where the frontend communicates with the backend through HTTP requests and WebSocket connections for real-time features.

---

# Social Network Application

## Overview

This project is a full‑stack **Social Network** application that enables users to register, authenticate, share content, interact with other users, and communicate in real time. The system follows a **client–server architecture** with a Go backend and a modern web frontend.

The application provides a complete social platform including user profiles, posts, comments, groups, messaging, and notifications.

The backend is implemented in **Go**, using **SQLite** as the database and **WebSockets** for real‑time communication. The frontend communicates with the backend through HTTP APIs and WebSocket connections. The entire system can be deployed easily using **Docker and Docker Compose**.

---

# Key Features

### User Management

* User registration and authentication
* Session‑based authentication
* Secure password hashing using **bcrypt**
* User profiles with optional avatar, nickname, and biography

### Social Interaction

* Create and comment on posts
* Follow and interact with other users
* Upload images with posts

### Groups

* Create and manage groups
* Join groups and interact with members
* Group posts and comments
* Group events and event responses

### Messaging

* Private messaging between users
* Real‑time message updates using **WebSockets**

### Notifications

* Follow notifications
* Group invitations and join requests
* Interaction notifications

---

# System Architecture

The application follows a **modular layered architecture**:

* **Frontend** – User interface and interaction
* **Backend API** – Business logic and request handling
* **Database Layer** – SQLite with migration management
* **Real‑time Communication** – WebSockets for live updates

```
Frontend  →  Backend API  →  Database
      ↘           ↓
      WebSocket Communication
```

---

# Project Structure

```
social-network/
├── backend/
│   ├── main.go
│   ├── server.go
│   ├── middleware/
│   ├── pkg/
│   │   ├── authentication/
│   │   ├── db/
│   │   │   ├── migrations/
│   │   │   └── sqlite/
│   │   ├── groups/
│   │   ├── messages/
│   │   ├── notifications/
│   │   ├── posts/
│   │   ├── profile/
│   │   └── websocket/
│   └── utils/
│
├── frontend/
│   ├── src/
│   ├── public/
│   └── package.json
│
├── docker-compose.yml
├── README.md
└── docs/
```

### Backend

Handles:

* API endpoints
* authentication and session management
* database access
* migrations
* WebSocket communication

### Frontend

Responsible for:

* user interface
* API communication
* real‑time updates

---

# Technologies Used

## Backend

* **Go (Golang)**
* **SQLite** database
* **Gorilla WebSocket** for real‑time communication
* **golang‑migrate** for database migrations
* **bcrypt** for password hashing
* **UUID** generation for unique identifiers

## Frontend

* JavaScript
* Modern frontend tooling
* REST API communication

## DevOps

* Docker
* Docker Compose

---

# Installation and Running the Project

## Prerequisites

Make sure the following tools are installed:

* Go
* Node.js and npm
* Docker
* Docker Compose
* SQLite (optional for manual inspection)

---

# Running with Docker (Recommended)

Navigate to the project root where `docker-compose.yml` is located.

Stop any running containers:

```
docker-compose down
```

(Optional) Remove previous containers:

```
docker rm -f social-network-backend social-network-frontend 2>/dev/null || true
```

Build and start the containers:

```
docker-compose up --build
```

After the containers start successfully, open the application in your browser:

```
http://localhost/
```

### Viewing Logs

```
docker-compose logs -f
```

### Stopping Containers

```
docker-compose down
```

---

# Running Without Docker

## Backend

```
cd backend
go run .
```

The backend server will start on:

```
http://localhost:8080
```

## Frontend

```
cd frontend
npm install
npm run dev
```

The frontend will run on:

```
http://localhost:5173
```

---

# Database

The application uses **SQLite** as its primary database.

Database schema changes are managed through **golang‑migrate**, ensuring that database structure is version‑controlled and automatically applied when the backend starts.

### Inspect the Database

```
sqlite3 <database_name>.db
```

### List Tables

```
.tables
```

### Check Migration Status

```
SELECT * FROM schema_migrations;
```

Example result:

```
10|0
```

This confirms that all migrations have been applied successfully.

---

# Authentication System

The application implements **session‑based authentication**.

### Authentication Features

* User registration
* Secure login/logout
* Password hashing with bcrypt
* Session cookies stored in the database

### Registration Fields

Required fields:

* Email
* Password
* First Name
* Last Name
* Date of Birth

Optional fields:

* Avatar / Profile Image
* Nickname
* About Me

---

# Docker Containers

The project runs using **two containers**:

| Container               | Purpose                     |
| ----------------------- | --------------------------- |
| social-network-backend  | Runs Go backend server      |
| social-network-frontend | Serves frontend application |

Check containers:

```
docker ps -a
```

Expected containers:

```
social-network-backend
social-network-frontend
```

---

# Audit Compliance

This project satisfies the main requirements of the social network specification:

* Allowed Go packages respected
* SQLite used as the primary database
* Database migrations implemented and applied
* Session‑based authentication system
* Organized backend and frontend architecture
* Docker containerization for deployment

---

# Contributors

Georgia
Theochara
Sofia
Andriana
Iana

