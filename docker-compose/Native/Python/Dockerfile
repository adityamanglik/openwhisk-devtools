# Use a base image with Python installed. Adjust the version as needed.
FROM python:3.10

# Set the working directory inside the container
WORKDIR /app

# Copy the source file(s) into the container
COPY . .

# (Optional) If your server requires any external Python packages, 
# you can list them in a requirements.txt file and uncomment the following line:
RUN pip install --no-cache-dir -r psutil

# There are no direct Python equivalents for the Node.js-specific environment variables 
# and command-line flags shown in your JavaScript Dockerfile example.
# Python runtime behavior can be influenced by other means if necessary.

# Command to run the application. Replace 'server.py' with the name of your Python server script.
CMD ["python", "server.py"]
