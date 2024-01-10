# Use a base image with Go installed (specify the Go version you need)
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the Go source file(s) into the container
COPY . .

ENV GOMAXPROCS=1
# Set GOGC environment variable
ENV GOGC=1000
ENV GOMEMLIMIT 128MiB

# Build the Go application
RUN go build -o server

# Command to run the application
CMD ["./server"]