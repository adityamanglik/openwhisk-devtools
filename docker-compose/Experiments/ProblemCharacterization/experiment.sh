# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.87:8180/java"
GO_API="http://128.110.96.87:8180/go"
KILL_SERVER_API="http://128.110.96.87:8180/exitCall"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose/Experiments"
JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# Send request and measure request response latencies
send_requests() {
    local size=$1

    # compile the docker images
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/../Native/Go/; docker build -t go-server-image ."

    # Change fakerequestarraysize
    ssh $OW_SERVER_NODE "sed -i 's/fakeRequestArraySize = [^ ]*/fakeRequestArraySize = $size/' /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go"

    # # Restart docker for good measure
    # ssh $OW_SERVER_NODE "sudo systemctl restart docker"

    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    sleep 5
    
    # Start sending requests
    taskset -c 2 go run request_sender.go $size
    
    # Move files for postprocessing
    mv $OW_DIRECTORY/ProblemCharacterization/go_response_times.txt "$OW_DIRECTORY/ProblemCharacterization/Graphs/Go/$size/client_time.txt"
    mv $OW_DIRECTORY/ProblemCharacterization/go_server_times.txt "$OW_DIRECTORY/ProblemCharacterization/Graphs/Go/$size/server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/Go/$size/memory.txt"
    # SCP the server.log file along with other files
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/Go/$size/server.log"
    
    # Remove file after retrieving
    ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/../LoadBalancer/*.log"

}

sizes=(100000)
# sizes=(1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000)

# Loop through each size
for size in "${sizes[@]}"; do
    # Kill the load balancer process if running
    curl $KILL_SERVER_API

    send_requests $size

    # Kill the load balancer process if running
    curl $KILL_SERVER_API

    # Plot time responses
    python ./Graphs/response_time_plotter.py "./Graphs/Go/${size}/client_time.txt" "./Graphs/Go/${size}/server_time.txt" "./Graphs/Go/${size}/memory.txt" "./Graphs/Go/${size}/distribution.pdf" "./Graphs/Go/${size}/latency.pdf"
    
    # Plot memory patterns
    python ./Graphs/go_mem_plotter.py "./Graphs/Go/${size}/memory.txt" "./Graphs/Go/${size}/memory.pdf"

    # Calculate impact of GC
    python ./Graphs/analyzer.py "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/ProblemCharacterization/Graphs/Go/$size/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/ProblemCharacterization/Graphs/Go/$size/server_time.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/ProblemCharacterization/Graphs/Go/$size/client_time.txt"
done