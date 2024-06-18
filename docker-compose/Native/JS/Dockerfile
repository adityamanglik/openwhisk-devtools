# Use a base image with NodeJS
# Add 14 for repeating JS experiments and plot sawtooth
FROM node:14

# Set the working directory inside the container
WORKDIR /app

# Copy the source file(s) into the container
COPY . .

# Copy the resources directory into the container
COPY Resources /app/Resources

# Set environment variable
# ENV GOMAXPROCS=1
# ENV UV_THREADPOOL_SIZE=1

# Command to run the application
# CMD ["node", "--max-old-space-size=64", "--max-semi-space-size=32", "--trace_gc", "--nouse-idle-notification", "--noincremental_marking=false", " --noincremental_marking_steps", "--nolazy_sweeping", "server.js"]
CMD ["node", "--max-old-space-size=64", "--max-semi-space-size=32", "server.js"]