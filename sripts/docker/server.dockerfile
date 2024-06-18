# Start from the official Go image
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/server

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Set environment variables with default values
ENV DDOS_PROTECTION=OFF
ENV PORT=8080
ENV READ_TIMEOUT=750ms
ENV WRITE_TIMEOUT=750ms

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary program using environment variables directly
CMD ["sh", "-c", "./main -DDOS_PROTECTION $DDOS_PROTECTION -PORT $PORT -READ_TIMEOUT $READ_TIMEOUT -WRITE_TIMEOUT $WRITE_TIMEOUT"]
