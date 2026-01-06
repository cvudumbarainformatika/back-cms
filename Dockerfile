# Build stage
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM --platform=linux/amd64 alpine:latest

# Install ca-certificates and curl for healthchecks
RUN apk --no-cache add ca-certificates curl

# WORKDIR /root/
# 🔑 GANTI workdir ke /app (penting)
WORKDIR /app

# 🔑 TAMBAHKAN INI
RUN mkdir -p \
    /app/storage/avatars \
    /app/storage/thumbnails \
    /app/storage/berita \
    /app/storage/dokumen \
    /app/storage/galeri \
    /app/storage/attachments


# Copy the binary and required assets from builder
COPY --from=builder /app/main .
COPY --from=builder /app/database/migrations ./database/migrations

# Copy .env file if exists (optional, better to use env vars)
# COPY .env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
