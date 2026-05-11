# Use the official Go image based on Alpine Linux as the build stage
FROM golang:1.22.3-alpine3.18 AS builder

# Install bash in the build stage
RUN apk add --no-cache bash

# Set the working directory inside the container to /app
WORKDIR /app

# Copy all files from the current directory to the working directory in the container
COPY . .

# Build the Go application and output the binary as 'main'
RUN go build -o main ./main.go

# Use the official Alpine Linux image as the base for the run stage
FROM alpine:3.18

# Set the working directory inside the container to /app
WORKDIR /app

# Copy the built binary from the builder stage to the current stage
COPY --from=builder /app/main .

# Copy the templates directory to the working directory in the container
COPY templates templates/

# Copy the static directory to the working directory in the container
COPY static static/

# Copy the banners directory to the working directory in the container
COPY utils/banners/ utils/banners/

# Expose port 8080 to allow external access
EXPOSE 8080

# Specify the command to run the application when the container starts
CMD ["./main"]