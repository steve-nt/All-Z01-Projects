FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/client
RUN go build -o frontend-server main.go

EXPOSE 3000

CMD ["./frontend-server"]

