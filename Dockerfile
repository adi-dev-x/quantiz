# Use the official Golang image as a base image
FROM golang:1.22.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]