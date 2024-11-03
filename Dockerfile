# The build stage
FROM golang:1.22 as builder
# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the application in release mode
RUN go build -o app ./cmd/api

# Use a minimal image for the final build (scratch or Alpine)
FROM alpine:latest

# Set environment variables for your app
ENV PORT=8000

# Install SSL certificates to support HTTPS if needed
RUN apk add --no-cache ca-certificates

# Copy the built Go binary from the builder stage
COPY --from=builder /app/app /app/app

# Expose the port your app will run on
EXPOSE 8000

# Start the application
CMD ["/app/app"]