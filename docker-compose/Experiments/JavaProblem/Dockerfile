# Use a base image with Java (specify the version you need)
FROM amazoncorretto:17

# Set the working directory inside the container
WORKDIR /app

# Copy your Java files and any other necessary files into the container
COPY . /app

# Compile your Java code (adjust the command according to your needs)
RUN javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java

# ENV GC_FLAGS="-XX:+UseSerialGC -Xms10g -Xmx10g"
# Print the GC flags to ensure they are included
RUN echo "GC_FLAGS: $GC_FLAGS"

CMD java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer
