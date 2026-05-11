FROM golang:1.23.2-bullseye
LABEL maintainer="Giannis Georgakopoulos"
LABEL version="1.0"
WORKDIR /app
COPY . .
RUN go build
EXPOSE 8080
CMD ["./groupie-tracker"]