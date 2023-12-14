# Use a base image with Java (specify the version you need)
FROM openjdk:8

# Set the working directory inside the container
WORKDIR /app

# Copy your Java files and any other necessary files into the container
COPY . /app

# Compile your Java code (adjust the command according to your needs)
RUN javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java

ENV GC_FLAGS="-Xms4g -Xmx4g"

CMD java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer
