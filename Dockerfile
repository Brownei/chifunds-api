# Start from the official Go image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Run the Go application
CMD ["./main"]

