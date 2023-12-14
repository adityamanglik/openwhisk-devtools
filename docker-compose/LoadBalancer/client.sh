# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.76:8180/java"
GO_API="http://128.110.96.76:8180/go"
KILL_SERVER_API="http://128.110.96.76:8180/exitCall"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"
ITERATIONS=5000

# Build docker images
build_docker_images() {
# compile the docker images
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; docker build -t go-server-image ."
}

# Client code

# Send request and measure request response latencies
send_requests() {
    # local base_api_url=$1
    # local response_time_file=$2
    # local execution_time_file=$3
    local size=$4

    # Update the Java code with the new array size
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Native/Java/Hello.java"

    ssh $OW_SERVER_NODE "awk '/MARKER_FOR_SIZE_UPDATE/{print;getline;print \"const ARRAY_SIZE = \" $size \";\";next}1' $OW_DIRECTORY/Native/Go/server.go > $OW_DIRECTORY/Native/Go/temp.go && mv $OW_DIRECTORY/Native/Go/temp.go $OW_DIRECTORY/Native/Go/server.go"

    # build_docker_images with new size
    build_docker_images

    go run request_sender.go
    
    # Move files for postprocessing
    mv $OW_DIRECTORY/LoadBalancer/go_response_times.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Go/$size/client_time.txt"
    mv $OW_DIRECTORY/LoadBalancer/go_server_times.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Go/$size/server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/go_heap_memory.log ../Graphs/LoadBalancer/Go/$size/memory.txt
    # Remove file after retrieving
    # ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/LoadBalancer/go_heap_memory.log"

    mv $OW_DIRECTORY/LoadBalancer/java_response_times.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Java/$size/client_time.txt"
    mv $OW_DIRECTORY/LoadBalancer/java_server_times.txt "$OW_DIRECTORY/Graphs/LoadBalancer/Java/$size/server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/java_heap_memory.log ../Graphs/LoadBalancer/Java/$size/memory.txt
}

# Array of sizes
sizes=(100 10000 1000000 3200000)
# sizes=(10000)

# for size in "${sizes[@]}"; do
#     python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Java/${size}/client_time.txt" "../Graphs/LoadBalancer/Java/${size}/server_time.txt" "../Graphs/LoadBalancer/Java/${size}/graph.pdf"

#     python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Go/${size}/client_time.txt" "../Graphs/LoadBalancer/Go/${size}/server_time.txt" "../Graphs/LoadBalancer/Go/${size}/graph.pdf"
# done

for size in "${sizes[@]}"; do
    python ../Graphs/LoadBalancer/go_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Go/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Go/${size}/memory.pdf"
    python ../Graphs/LoadBalancer/java_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Java/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Java/${size}/memory.pdf"
done

# Loop through each size
for size in "${sizes[@]}"; do
    # Kill the load balancer process if running
    curl $KILL_SERVER_API

    # Restart docker for good measure
    ssh $OW_SERVER_NODE "sudo systemctl restart docker"

    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /dev/null 2>&1 &"

    # Commands for Java API
    send_requests $JAVA_API "client_time.txt" "server_time.txt" $size
    python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Java/${size}/client_time.txt" "../Graphs/LoadBalancer/Java/${size}/server_time.txt" "../Graphs/LoadBalancer/Java/${size}/graph.pdf"
    python ../Graphs/LoadBalancer/response_time_plotter.py "../Graphs/LoadBalancer/Go/${size}/client_time.txt" "../Graphs/LoadBalancer/Go/${size}/server_time.txt" "../Graphs/LoadBalancer/Go/${size}/graph.pdf"
    
    # Plot memory patterns
    python ../Graphs/LoadBalancer/go_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Go/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Go/${size}/memory.pdf"
    python ../Graphs/LoadBalancer/java_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Java/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Java/${size}/memory.pdf"
done

# BACKUP #######################################################################

# send_requests $JAVA_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/LoadBalancer/response_time_plotter.py ../Graphs/LoadBalancer/Java/1000000/client_time.txt ../Graphs/LoadBalancer/Java/1000000/server_time.txt ../Graphs/LoadBalancer/Java/1000000/graph.pdf

# send_requests $GO_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/LoadBalancer/response_time_plotter.py ../Graphs/LoadBalancer/Go/1000000/client_time.txt ../Graphs/LoadBalancer/Java/1000000/server_time.txt ../Graphs/LoadBalancer/Java/1000000/graph.pdf


# Locust code
# export API_URL=$GO_API
# locust --config=./master.conf

# export API_URL=$JAVA_API
# locust --config=./master.conf
