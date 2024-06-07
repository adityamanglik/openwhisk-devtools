# Use a base image with Java (specify the version you need)
FROM amazoncorretto:17

# Define build arguments
ARG GC_FLAGS

# Set the working directory inside the container
WORKDIR /app

# Copy your Java files and any other necessary files into the container
COPY . /app

# Copy the image resources into the container
COPY Resources /app/Resources

# Compile your Java code (adjust the command according to your needs)
RUN javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java ImageProcessor.java

# Print the build arguments to verify they are correctly passed
RUN echo "GC_FLAGS: $GC_FLAGS"

# Set environment variable
ENV GC_FLAGS=$GC_FLAGS

CMD java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer
