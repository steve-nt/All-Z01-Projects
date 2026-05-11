# Use the latest Go image for building
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all project files into the container
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Use a minimal image for runtime
FROM alpine:latest

# Set the working directory for runtime
WORKDIR /app

# Copy all files from the builder stage to runtime
COPY --from=builder /app .

# Expose the port the application runs on
EXPOSE 8080

# Set the entry point for the application
CMD ["./main"]
