# Use a base image with NodeJS
FROM node:14

# Set the working directory inside the container
WORKDIR /app

# Copy the source file(s) into the container
COPY . .

# Set environment variable
# ENV GOMAXPROCS=1
# ENV UV_THREADPOOL_SIZE=1

# Command to run the application
CMD ["node", "--max-old-space-size=96", "--max-semi-space-size=16", "--trace_gc", "--nouse-idle-notification", "server.js"]