# Use bash for better scripting
set shell := ["bash", "-c"]

# Default recipe
default: dev

# Run the app in dev mode with auto-reload
dev:
    cargo watch -x run

# Run the app normally (no watch)
start:
    cargo run

# Run tests
test:
    cargo test --all

# Format and lint code
fmt:
    cargo fmt --all
    cargo clippy --all-targets --all-features -- -D warnings

# Run SQLx migrations
migrate:
    sqlx migrate run

# Create a new migration: just new-migration name="add_users"
new-migration name:
    sqlx migrate add {{name}}

# Build for release
build:
    cargo build --release

# Build Docker image
docker-build:
    docker build -t ai_document_backend .

# Run docker-compose in dev mode
docker-up:
    docker compose -f compose.dev.yml --env-file .env.development up

# Stop docker containers
docker-down:
    docker compose -f compose.dev.yml --env-file .env.development down
