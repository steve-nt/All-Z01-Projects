# Project Architecture: Backend & Frontend

## Overview

This project has **two separate servers** that work together:

1. **Backend Server** (Go) - Handles API requests, database operations, authentication
2. **Frontend Server** (HTML/CSS/JS) - Serves the user interface, makes requests to backend

## Project Structure

```
social-network/
├── backend/
│   ├── main.go              # Backend server entry point (Go HTTP server)
│   ├── pkg/
│   │   ├── authentication/  # Auth handlers (register, login, logout)
│   │   ├── db/              # Database connection & migrations
│   │   └── ...
│   └── data/                # Database files (auto-created, gitignored)
│
└── frontend/                # Frontend files (HTML, CSS, JS)
    ├── templates/           # HTML templates
    ├── static/              # CSS, JS, images
    └── uploads/             # User-uploaded images
```

## How They Connect

### Option 1: Backend Serves Everything (Recommended for Development)

The Go backend server can serve **both**:
- **API endpoints** (e.g., `/api/register`, `/api/login`)
- **Static files** (HTML, CSS, JS from `frontend/`)

**Advantages:**
- Simple setup - one server
- No CORS issues
- Easy development

**How it works:**
```go
// In backend/main.go
// API routes
http.HandleFunc("/api/register", handlers.RegisterHandler)
http.HandleFunc("/api/login", handlers.LoginHandler)

// Serve static files
http.Handle("/", http.FileServer(http.Dir("../frontend")))
```

### Option 2: Separate Servers (For Production/Docker)

- **Backend**: Port 8080 (API only)
- **Frontend**: Port 3000 (static files, e.g., nginx, Caddy, or Node.js)

**Advantages:**
- Better separation of concerns
- Can scale independently
- Matches Docker requirements

**How they connect:**
- Frontend makes HTTP requests to `http://localhost:8080/api/*`
- Backend responds with JSON data
- Frontend renders the UI

## File Naming: `main.go` vs `server.go`

✅ **Use `main.go`** (current setup)
- Standard Go convention
- When you run `go run .` or `go build`, Go automatically finds `main.go`
- More conventional in Go projects

❌ **Don't use `server.go`**
- Less conventional
- Requires `go run server.go` explicitly

## Do You Need `main.go` in Root?

**No, you don't need it!**

The root `social-network/` folder is just the project container. The actual entry points are:
- `backend/main.go` - Backend server
- `frontend/` - Frontend files (served by backend or separate server)

## Running the Project

### Development (Backend serves everything):

```bash
# Terminal 1: Start backend
cd backend
go run main.go

# Backend runs on http://localhost:8080
# - API: http://localhost:8080/api/*
# - Frontend: http://localhost:8080/*
```

### Production (Separate servers):

```bash
# Terminal 1: Backend
cd backend
go run main.go
# Runs on :8080

# Terminal 2: Frontend (e.g., using a simple HTTP server)
cd frontend
python3 -m http.server 3000
# Or use nginx, Caddy, etc.
# Frontend makes requests to http://localhost:8080/api/*
```

## Docker Setup (Future)

As per requirements, you'll have:

1. **Backend Docker Container**
   - Runs `backend/main.go`
   - Exposes port 8080
   - Contains database and Go code

2. **Frontend Docker Container**
   - Serves static files
   - Exposes port 3000 (or 80)
   - Makes HTTP requests to backend container

## Current Setup

Right now, you have:
- ✅ `backend/main.go` - Backend server entry point
- ❌ `frontend/` - Not created yet (you'll create this)

**Next Steps:**
1. Create `frontend/` directory with HTML/CSS/JS
2. Update `backend/main.go` to serve static files OR set up separate frontend server
3. Add API routes to `backend/main.go`

