# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o user-service ./cmd/user-service

# Run stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/user-service .

# Copy .env file
COPY .env .

CMD ["./user-service"]
