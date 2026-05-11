# Ascii-Art-Web

## 🎨 Overview

Ascii-Art-Web is a web-based version of the ASCII-ART project that converts text input into ASCII art using predefined banners. It is implemented in Go and follows the Clean Architecture principle.

## 🚀 Features

- Web interface for generating ASCII art
- Multiple banner styles
- Follows Clean Architecture for maintainability and scalability
- Supports debugging with VS Code

## 🛠 Installation and Setup

### Prerequisites

- Go (1.22.9)
- Git

### Clone the Repository

```sh
git clone https://platform.zone01.gr/git/agkiata/ascii-art-web
cd ascii-art-web
```

### Install Dependencies

```sh
go mod tidy
```

### 📌 Run the Server

```sh
go run cmd/server/main.go
```

### 📌 Run the Client

```sh
go run cmd/client/main.go
```

## ⚡ Debug

To debug the project in VS Code, create `launch.json` in the `.vscode` folder.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Ascii Art Web Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "cmd/server/main.go"
        },
        {
            "name": "Ascii Art Web Client",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "cmd/client/main.go"
        }
    ]
}
```

## 🧩 Testing
Run unit tests with:
```sh
go test ./tests
```

## 📂 Project Structure

```
ascii-art-web/
├── assets/
│   ├── banners/
│   │   ├── standard.txt
│   │   ├── shadow.txt
│   │   └── thinkertoy.txt
│── cmd/
│   ├── client/    # Client-side code
│   │   └── main.go
│   ├── server/    # Server-side code
│   │   └── main.go
│── internal/      # Internal application logic
│   ├── adapter/
│   │     ├── handler/
│   │     │      └── ascii_handler.go
│   │     └── repository/
│   │            ├── banner_repository_test.go
│   │            └── banner_repository.go
│   ├── config/
│   │     └── config.go
│   ├── domain/
│   │     ├── ascii_text.go
│   │     └── banner.go
│   ├── infrastructure/
│   │     ├── router.go
│   │     └── server.go
│   ├── usecase/
│   │     ├── ascii_usecase_test.go
│   │     └── ascii_usecase.go
│── templates/     # Front End Templates
│   ├── ascii_art.html
│   ├── homePage.html
│   └── index.html
│── .gitignore 
│── config.json    # Configuration file
├── go.mod
└── README.md      # Project documentation
```
