# Use a minimal base image with Go
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the entire application code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal image for runtime
FROM alpine:latest

# Metadata for the Docker image
LABEL maintainer="Your Name <theocharoula.tarara@gmail.com>"
LABEL description="ASCII Art Web Application for converting text to ASCII art."
LABEL version="1.0"

# Set the working directory for runtime
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy required assets (templates, static files)
COPY templates ./templates
COPY static ./static
COPY handlers ./handlers 

# Expose the port the application runs on
EXPOSE 8080

# Set the entry point for the application
CMD ["./main"]
