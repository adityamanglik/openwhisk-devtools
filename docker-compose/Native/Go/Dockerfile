# Use a base image with Go installed (specify the Go version you need)
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source file(s) into the container
COPY . .

# Build the Go application
RUN go build -o server

# Expose the port the server listens on
EXPOSE 9875

# Command to run the application
CMD ["./server"]