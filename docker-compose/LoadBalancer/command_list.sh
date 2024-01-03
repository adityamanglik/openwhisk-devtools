# This script starts the load balancer server
# The load balancer server receives all requests
# As it receives a request, it checks the endpoint
# Based on the endpoint, it starts the Java or Go container
# Once the container is started, it sends the request to the container
# The container returns the request's result to the load balancer
# The container also returns the current GC stats to the LB
# The load balancer forwards the request's result to the client
# The load balancer looks at the GC data and decides where to send requests

# Go
docker build -t go-server-image .
docker run -d  --rm --name my-go-server -p 9875:9875 go-server-image
curl http://128.110.96.176:9875/GoNative?seed=654
docker stop my-go-server
docker rm my-go-server


# Java
docker build -t java-server-image .
docker run -d --rm --name my-java-server -p 9876:9876 java-server-image
curl http://128.110.96.176:9876/jsonresponse?seed=654
docker stop my-java-server
docker rm my-java-server


# Ensure your Go and Java server images are built and available.
# Run this load balancer Go program.
# Send requests to http://localhost:8080/java?seed=654 or http://localhost:8080/go?seed=654.
