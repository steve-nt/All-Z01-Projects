FROM golang:1.23-alpine

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod/go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build backend binary
WORKDIR /app/cmd/app
RUN go build -o backend-server main.go

EXPOSE 8080
CMD ["./backend-server"]

