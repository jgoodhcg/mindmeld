# Build stage (Debian-based for glibc compatibility with Tailwind binary)
FROM --platform=linux/amd64 golang:1.25 AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates && rm -rf /var/lib/apt/lists/*

# Install Tailwind CSS standalone CLI
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate code BEFORE go mod tidy (so internal/db/ exists for import resolution)
RUN go run github.com/a-h/templ/cmd/templ generate
RUN go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Ensure go.mod is tidy (must run after code generation)
RUN go mod tidy

# Build Tailwind CSS
RUN tailwindcss -i styles/input.css -o static/css/output.css --minify

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM --platform=linux/amd64 alpine:3.20

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezones
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files
COPY --from=builder /app/static ./static

# Copy migrations for goose
COPY --from=builder /app/migrations ./migrations

# Expose port (configurable via PORT env var)
EXPOSE 8080

CMD ["./server"]
