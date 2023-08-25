FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /script_sight

# Copy go mod and sum files
COPY api/go.mod api/go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code
COPY api /api/

COPY frontend /frontend/

# Build the application
WORKDIR /api/cmd/
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
