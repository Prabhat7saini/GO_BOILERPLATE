# Stage 1: Build
FROM golang:1.24.4 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# 👇 statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o boilerplate ./cmd/main.go

# Stage 2: Slim runtime (no glibc needed)
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=build /app/boilerplate .
COPY .env .env

EXPOSE 4000
CMD ["./boilerplate"]
