# Forum — Go Web Forum with SQLite & Docker

This project is a **web forum** built in **Go**, using **SQLite** for data storage and **Docker** for containerization.

It allows users to register, log in, create posts, comment, react (like/dislike), and browse posts by category, while enforcing authentication, session handling, and proper HTTP error management.

---

## ✨ Features

### 👤 Authentication
- Register with **email**, **username**, and **password**
- Passwords hashed using **bcrypt**
- Login / logout
- One active session per user (cookie-based)
- Duplicate email and username detection

### 📝 Forum Functionality
- Create posts (logged-in users only)
- Comment on posts (logged-in users only)
- View posts and comments (everyone)
- Predefined categories:
  - General
  - Technology
  - Gaming
  - Movies & TV
  - Music
  - Sports
  - Help
- Filter posts by category
- Like / dislike posts and comments

### 🗄 Database
- SQLite database file: `forum.db`
- Tables:
  - `users`
  - `sessions`
  - `posts`
  - `comments`
  - `categories`
  - `post_categories`
  - `reactions`
- Uses `CREATE`, `INSERT`, and `SELECT` SQL queries

### ⚠️ Error Handling
- Custom error pages for:
  - 400 Bad Request
  - 404 Not Found
  - 500 Internal Server Error
- Proper HTTP status codes returned

### 🎨 UI
- Dark theme (eye-friendly)
- Responsive layout
- Hover animations
- Emoji icons (no frontend frameworks)

---

## 🚀 Run Locally (Without Docker)

### Requirements
- Go installed
- SQLite installed (optional, for inspecting the DB)

### Steps 
```bash
git clone https://platform.zone01.gr/git/sgougoul/forum
cd forum
go mod tidy
go run .
```
### 🚀 Run With Docker

### Requirements
- Docker installed

### Steps

### Build the Docker image
```bash   docker image build -t forum .```

### Run the container (with database mounted)
```bash  
docker rm -f forum 2>/dev/null || true

docker container run --name forum \
  -p 8080:8080 \
  -v "$(pwd)/forum.db:/app/forum.db" \
  forum
```
### The forum will be available at : 
     http://localhost:8080
     
### Stop the container /remove the container 
```bash  
  docker stop forum
  docker rm forum 
```

### Check running containers 
```bash  
  docker ps
```
### Show disk usage
```bash
  docker system df
```
### Clean unused Docker resources
```bash
  docker system prune -f  
```
### IF YOU DON'T WANT TO RUN THE COMMANDS BY YOURSELF, run the following command to enable a script that runs the docker automatically  :
```bash
chmod +x run_docker.sh
./run_docker.sh
```