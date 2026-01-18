# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies (if needed)
RUN apk add --no-cache git

# Set destination for COPY
WORKDIR /app

# Download Go modules (leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./
COPY resources/ ./resources/

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /lz-print-agent-local

# Final stage - use scratch for minimal image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /lz-print-agent-local /lz-print-agent-local

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/lz-print-agent-local"]
