# ---------- Stage 1: Build ----------
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o main .

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install SSL certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from builder stage
COPY --from=builder /app/main .

# Expose the port your app listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
