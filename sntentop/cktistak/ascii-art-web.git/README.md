# ASCII Art Web

A web application that converts text into ASCII art using different banner styles. Users can input text and select from various fonts to generate stylized ASCII art output.

## How to Run

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Project Structure

```
ascii-art-web/
├── main.go              # Entry point
├── main_test.go         # Integration and main package tests
├── go.mod               # Go module definition
├── LICENSE              # MIT License
├── handlers/
│   ├── handlers.go      # HTTP request handlers
│   └── handlers_test.go # Handler unit tests
├── services/
│   ├── ascii-art.go     # ASCII art generation logic
│   └── ascii-art_test.go # Service unit tests
├── templates/
│   ├── index.html       # Main page template
│   └── error.html       # Error page template
└── banners/
    ├── standard.txt     # ASCII art fonts
    ├── shadow.txt
    ├── thinkertoy.txt
    └── zigzag.txt
```

## Error Status Testing

### 400 Bad Request
- Submit empty text
- Submit text with non-ASCII characters

### 404 Not Found
- Visit any URL other than `/` or `/ascii-art`
- Example: `http://localhost:8080/nonexistent`

### 500 Internal Server Error
- Delete the `/templates/` directory

## Contributors

- ### Constantine Ktistakis
- ### Giorgos Koutzos
- ### Ioannis Kountouris