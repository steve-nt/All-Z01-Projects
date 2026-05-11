# Build stage: compile the Go backend (CGO required for go-sqlite3)
FROM golang:1.24.1-bookworm AS builder

ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o server .

# Runtime stage: minimal image to run the server
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary, migrations, and assets from builder
COPY --from=builder /build/server .
COPY --from=builder /build/pkg/db/migrations/sqlite /app/pkg/db/migrations/sqlite
COPY --from=builder /build/frontend ./frontend

# Create directories for database and uploads (writable at runtime)
RUN mkdir -p data frontend/uploads/images frontend/uploads/thumbnails

EXPOSE 8080

# Run as non-root when possible (optional: add user and chown)
CMD ["./server"]
