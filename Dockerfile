# Stage 1: Build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO_ENABLED=0 penting untuk aplikasi statis yang ringan
RUN CGO_ENABLED=0 GOOS=linux go build -o main-app ./cmd/api

# Stage 2: Runtime
FROM alpine:latest
# Keamanan: Tambahkan sertifikat & user non-root
RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/main-app .

# Jalankan sebagai non-root
USER appuser

EXPOSE 8080
CMD ["./main-app"]