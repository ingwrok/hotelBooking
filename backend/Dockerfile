# Builder Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependencies first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o hotel-api cmd/http/main.go

# Runner Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (e.g. ca-certificates for external API calls, tzdata for timezones)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/hotel-api .
COPY --from=builder /app/config.yaml . 
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8000

# Run
CMD ["./hotel-api"]
