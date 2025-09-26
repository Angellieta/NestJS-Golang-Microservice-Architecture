# order-service/Dockerfile

# --- Tahap 1: Build ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 1. Copy file modul terlebih dahulu
COPY go.mod go.sum ./

# 2. Download dependensi
#    Docker akan menyimpan layer ini di cache. Jika go.mod/go.sum tidak berubah
RUN go mod download

# 3. Copy sisa source code aplikasi
COPY . .

# 4. Build aplikasi Go menjadi satu file binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server

# --- Tahap 2: Production ---
FROM alpine:latest

WORKDIR /

# Copy file binary yang sudah di-build dari tahap "builder"
COPY --from=builder /main /main

EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["/main"]