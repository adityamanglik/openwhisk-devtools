# Use a base image with Go installed (specify the Go version you need)
FROM golang
# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod ./

# Copy the resources directory into the container
COPY Resources /app/Resources

# Download all dependencies
RUN go mod download

# Copy the Go source file(s) into the container
COPY . .

# Accept build arguments for GOGC and GOMEMLIMIT
ARG GOGC=1000
ARG GOMEMLIMIT=128MiB
ARG GOMAXPROCS=1

# Set environment variables
ENV GOGC=$GOGC
ENV GOMEMLIMIT=$GOMEMLIMIT
ENV GOMAXPROCS=$GOMAXPROCS

# Build the Go application
RUN go build -o server

# Command to run the application
CMD ["./server"]