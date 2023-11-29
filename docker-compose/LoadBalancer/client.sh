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
docker run -d --name my-go-server -p 9875:9875 my-go-server
curl http://128.110.96.176:9875/GoNative?seed=654
docker stop my-go-server
docker rm my-go-server


# Java
docker build -t java-server-image .
docker run -d --name my-java-server -p 9876:9876 my-java-server
curl http://128.110.96.176:9876/jsonresponse?seed=654
docker stop my-java-server
docker rm my-java-server


Ensure your Go and Java server images are built and available.
Run this load balancer Go program.
Send requests to http://localhost:8080/java?seed=654 or http://localhost:8080/go?seed=654.



# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.176:8080/java"
GO_API="http://128.110.96.176:8080/go"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"

# Build docker images
build_docker_images() {
# compile the server
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java"
# Build Java docker image
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t my-java-server ."

# Build Go docker image
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; docker build -t my-go-server ."
}

# Client code
build_docker_images
export API_URL=$GO_API
locust --config=./master.conf

export API_URL=$JAVA_API
locust --config=./master.conf
