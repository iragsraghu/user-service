FROM golang:1.23.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /user-service ./cmd/user-service

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /user-service /user-service
EXPOSE 8080
ENTRYPOINT ["/user-service"]
