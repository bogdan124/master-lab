# Use an official Go runtime as a parent image
FROM golang:1.18.1

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the Go application code and go.mod/go.sum files into the container
COPY go.mod .
COPY main.go .
COPY go.sum .

# Download the dependencies
RUN go mod download

# Create a non-root user
RUN useradd -u 10001 appuser

# Create the build cache directory with the correct permissions
RUN mkdir -p /home/appuser/.cache/go-build && chown -R appuser:appuser /home/appuser/.cache

# Set ownership of the entire /go directory to the non-root user during build
RUN chown -R appuser:appuser /go

# Switch to the non-root user
USER appuser

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o libp2p-app .

# Expose the necessary ports
EXPOSE 4001 5001

# Run the Go application on container startup
CMD ["./libp2p-app"]
