# Use the official Golang image for building the app
FROM golang:1.22.2 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o main cmd/api/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Install required dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 3000

# Run the binary
CMD ["./main"]
