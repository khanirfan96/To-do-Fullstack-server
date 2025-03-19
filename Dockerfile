# Step 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download the Go Modules dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp .

# Step 2: Create a small image to run the app
FROM alpine:latest  

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/myapp .

# Expose port 8080 (change if your app runs on a different port)
EXPOSE 8000

# Run the binary
CMD ["./myapp"]
