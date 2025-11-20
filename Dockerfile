FROM ghcr.io/go-task/task:latest AS tasker

# ==========================================
# Stage 1: Builder
# ==========================================
FROM golang:1.25-alpine AS builder

COPY --from=tasker /usr/local/bin/task /usr/local/bin/task

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
# If these files don't change, the `go mod download` step is skipped
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
# CGO_ENABLED=0: Disables CGO for a static binary (required for Alpine/Scratch)
# -ldflags="-w -s": Strips debug information to reduce binary size
RUN task build-release

# ==========================================
# Stage 2: Runner
# ==========================================
FROM alpine:3.22

WORKDIR /app

# Install CA certificates (essential for making HTTPS requests)
RUN apk --no-cache add ca-certificates

# Security: Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY .env .
COPY migrations ./migrations

# Switch to the non-root user
USER appuser

# Expose the application port (adjust as needed)
EXPOSE 7202

# Run the binary
CMD ["./main"]
