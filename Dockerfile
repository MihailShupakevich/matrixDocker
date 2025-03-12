
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.mod go.sum go.sum ./

RUN go mod download

COPY . .

#CMD ["go", "run", "cmd/main.go"],