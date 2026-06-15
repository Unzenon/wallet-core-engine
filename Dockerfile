# Stage 1: Build binary file Go
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copy file dependency dulu (untuk optimasi cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# UBAH BARIS INI: Arahkan ke folder tempat main.go berada
RUN go build -o main-app ./cmd/api

# Stage 2: Jalankan aplikasi dengan OS minimal (Alpine)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main-app .

EXPOSE 8080
CMD ["./main-app"]

