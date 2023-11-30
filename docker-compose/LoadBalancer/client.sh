# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.167:8180/java?seed=5"
GO_API="http://128.110.96.167:8180/go?seed=5"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"
ITERATIONS=1000

# Build docker images
build_docker_images() {
# compile the docker images
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; docker build -t go-server-image ."
}

# Client code

# Send request and measure request response latencies
send_requests() {
    local api_url=$1
    local response_time_file=$2
    local execution_time_file=$3
    local size=$4

    # Update the Java code with the new array size
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Native/Java/Hello.java"

    ssh $OW_SERVER_NODE "awk '/MARKER_FOR_SIZE_UPDATE/{print;getline;print \"const ARRAY_SIZE = \" $size \";\";next}1' $OW_DIRECTORY/Native/Go/server.go > $OW_DIRECTORY/Native/Go/temp.go && mv $OW_DIRECTORY/Native/Go/temp.go $OW_DIRECTORY/Native/Go/server.go"

    # build_docker_images with new size
    build_docker_images

    for i in $(seq 1 $ITERATIONS)
    do
        # Measure the response time and capture the response
        start_time=$(date +%s.%N)
        response=$(curl -s "$api_url")
        end_time=$(date +%s.%N)
        elapsed=$(echo "scale=3; ($end_time - $start_time) * 1000" | bc)

        # Parse the response to extract executionTime in milliseconds
        execution_time=$(echo $response | jq -r '.executionTime')

        # Convert execution time to milliseconds if needed
        execution_time_ms=$(echo "scale=3; $execution_time / 1000000" | bc)

        # Record the total elapsed time
        echo "$elapsed" >> $response_time_file

        # Record the extracted executionTime in milliseconds
        echo "$execution_time_ms" >> $execution_time_file
    done

    # Move all log and image files to their respective directories
    # Check if the API is for Go and then move all log and image files to the Go directory
    if [[ "$api_url" == *"/go"* ]]; then
        mv $OW_DIRECTORY/LoadBalancer/*.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Go/$size/"
    else
        # If the API is not for Go, assume it's for Java and move files to the Java directory
        mv $OW_DIRECTORY/LoadBalancer/*.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Java/$size/"
    fi
}

# Array of sizes
sizes=(100 10000 1000000)

# Loop through each size
for size in "${sizes[@]}"; do
    # Kill the load balancer process if running
    curl http://128.110.96.167:8180/exitCall

    # Restart docker for good measure
    ssh $OW_SERVER_NODE "sudo systemctl restart docker"

    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /dev/null 2>&1 &"

    # Commands for Java API
    send_requests $JAVA_API "client_time.txt" "server_time.txt" $size
    python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Java/${size}/client_time.txt" "../Graphs/LoadBalancer/Java/${size}/server_time.txt" "../Graphs/LoadBalancer/Java/${size}/graph.pdf"

    # Commands for Go API
    send_requests $GO_API "client_time.txt" "server_time.txt" $size
    python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Go/${size}/client_time.txt" "../Graphs/LoadBalancer/Go/${size}/server_time.txt" "../Graphs/LoadBalancer/Go/${size}/graph.pdf"
done

# send_requests $JAVA_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/LoadBalancer/response_time_plotter.py ../Graphs/LoadBalancer/Java/1000000/client_time.txt ../Graphs/LoadBalancer/Java/1000000/server_time.txt ../Graphs/LoadBalancer/Java/1000000/graph.pdf

# send_requests $GO_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/LoadBalancer/response_time_plotter.py ../Graphs/LoadBalancer/Go/1000000/client_time.txt ../Graphs/LoadBalancer/Java/1000000/server_time.txt ../Graphs/LoadBalancer/Java/1000000/graph.pdf

# BACKUP #######################################################################
# Locust code
# export API_URL=$GO_API
# locust --config=./master.conf

# export API_URL=$JAVA_API
# locust --config=./master.conf
