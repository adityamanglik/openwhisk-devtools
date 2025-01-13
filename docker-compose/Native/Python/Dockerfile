# Use a base image with Python installed. Adjust the version as needed.
FROM python:3.10

# Set the working directory inside the container
WORKDIR /app

# Copy the source file(s) into the container
COPY . .

# Copy the resources directory into the container
COPY Resources /app/Resources

# (Optional) If your server requires any external Python packages, 
# you can list them in a requirements.txt file and uncomment the following line:
RUN pip install --no-cache-dir -r requirements.txt

# Accept a build argument for the memory limit and set it as an environment variable
ARG MEMORY_LIMIT

# Command to run the application. Replace 'server.py' with the name of your Python server script.
CMD ["python", "pyaes.py"]
