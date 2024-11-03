# The build stage
FROM golang:1.23
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

# Set environment variables for your app
ENV PORT=8000

# Expose the port your app will run on
EXPOSE 8000

# Start the application
CMD ["/app/app"]