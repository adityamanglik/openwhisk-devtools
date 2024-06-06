# Use a base image with Java (specify the version you need)
FROM amazoncorretto:17

# Define build arguments
ARG GC_FLAGS
ARG GC
ARG MEM_SIZE

# Set the working directory inside the container
WORKDIR /app

# Copy your Java files and any other necessary files into the container
COPY . /app

# Compile your Java code (adjust the command according to your needs)
RUN javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java

# Print the build arguments to verify they are correctly passed
RUN echo "GC_FLAGS: $GC_FLAGS"
RUN echo "GC: $GC"
RUN echo "MEM_SIZE: $MEM_SIZE"

# Set environment variable
ENV GC_FLAGS=$GC_FLAGS

CMD java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer
