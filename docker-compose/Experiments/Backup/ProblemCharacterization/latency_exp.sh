# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://node0:8180/java"
GO_API="http://node0:8180/go"
KILL_SERVER_API="http://node0:8180/exitCall"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose/Experiments"
JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# Send request and measure request response latencies
send_requests() {
    local size=$1

    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go GCMitigation > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    sleep 10
    
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

# Send request and measure request response latencies
send_requests_NOGC() {
    local size=$1

    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    sleep 10
    
    # Start sending requests
    taskset -c 2 go run request_sender.go $size
    
    # Move files for postprocessing
    mv $OW_DIRECTORY/ProblemCharacterization/go_response_times.txt "$OW_DIRECTORY/ProblemCharacterization/Graphs/Go/$size/NOGC_client_time.txt"
    mv $OW_DIRECTORY/ProblemCharacterization/go_server_times.txt "$OW_DIRECTORY/ProblemCharacterization/Graphs/Go/$size/NOGC_server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/Go/$size/NOGC_memory.txt"
    # SCP the server.log file along with other files
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/Go/$size/NOGC_server.log"
    
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

    send_requests_NOGC $size

    # Plot time responses
    python ./Graphs/EM_NOGC_plotter.py "./Graphs/Go/${size}/client_time.txt" "./Graphs/Go/${size}/server_time.txt" "./Graphs/Go/${size}/NOGC_client_time.txt" "./Graphs/Go/${size}/NOGC_server_time.txt" "./Graphs/Go/${size}/memory.txt" "./Graphs/Go/${size}/NOGC_memory.txt" "./Graphs/Go/${size}/EM_NOGC_latency.pdf"
done