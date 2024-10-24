# Build stage
FROM golang:1.20-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary and place it in the /bin directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/api

# Deploy stage
FROM alpine:latest

# Set up a non-root user to avoid running as root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the build stage
COPY --from=build /bin/app /bin/app

# Set permissions and make it executable
RUN chmod +x /bin/app

# Switch to non-root user
USER appuser

# Expose the application's port (default port 8080 for Railway)
EXPOSE 8080

# Run the Go binary
CMD ["/bin/app"]
