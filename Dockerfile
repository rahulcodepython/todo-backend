# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# -ldflags="-w -s" strips debugging information and symbols, reducing binary size.
# CGO_ENABLED=0 creates a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o todo-backend ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/todo-backend .

# Copy the production environment file
# Note: In a real production scenario, you would manage secrets using
# environment variables injected by your orchestration tool (e.g., Kubernetes, Docker Swarm), not by copying a .env file.
COPY .env .env

# Expose the port the app runs on
EXPOSE 8000

# Set the command to run the application
# APP_ENV is set here to ensure the production environment is loaded
CMD ["/app/todo-backend"]