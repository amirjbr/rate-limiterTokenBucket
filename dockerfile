# ---------- Stage 1: Build ----------
FROM golang:1.24.4 AS builder

# Disable CGO for alpine compatibility
ENV CGO_ENABLED=0 GO111MODULE=on

WORKDIR /app

# Copy go.mod, go.sum, and vendor folder
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY .env ./

# Copy the rest of the source code
COPY . .

# Build the Go app using vendor folder
RUN go build -mod=vendor -o rateLimiter ./cmd

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install SSL certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/rateLimiter .

# Expose the port your app listens on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./rateLimiter"]
