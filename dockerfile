# ---------- Stage 1: Build ----------
FROM golang:1.24.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o ./cmd/main .

# ---------- Stage 2: Run ----------
FROM alpine:latest



RUN mkdir /app

WORKDIR /app
# Copy the compiled binary from builder stage
COPY --from=builder /app/app .

# Expose the port your app listens on
EXPOSE 8080
EXPOSE 8081

# Command to run the executable
ENTRYPOINT ["./app"]
