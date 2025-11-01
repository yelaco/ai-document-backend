# ===== 1. Builder Image =====
FROM rust:1.91-slim AS builder

# Set a working directory
WORKDIR /app

# Install build dependencies (optional but common)
RUN apt-get update && apt-get install -y --no-install-recommends pkg-config libssl-dev && rm -rf /var/lib/apt/lists/*

# Cache dependencies first
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release || true

# Copy full source
COPY . .

# Build release binary
RUN cargo build --release

# ===== 2. Runtime Image =====
FROM debian:stable-slim AS runtime

WORKDIR /app

# Install system dependencies if needed (OpenSSL is common for Rust apps)
RUN apt-get update && apt-get install -y --no-install-recommends libssl3 ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/target/release/ai_document_backend /app/ai_document_backend

# Copy ONLY the config for this environment
COPY config ./config

EXPOSE 8080

CMD ["./ai_document_backend"]
