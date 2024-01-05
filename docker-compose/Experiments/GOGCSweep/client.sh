# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
GO_API="http://128.110.96.59:8180/go"
KILL_SERVER_API="http://128.110.96.59:8180/exitCall"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# Array of sizes
sizes=(100 10000 1000000)
# Array of GOGC values
GOGC=(-1 1 10 50 100 200 400 800)

# Send request and measure request response latencies
send_requests() {
    local size=$1

    # compile the docker images
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; docker build -t go-server-image ."

    # # Restart docker for good measure
    # ssh $OW_SERVER_NODE "sudo systemctl restart docker"

    # Restart the load balancer
    ssh $OW_SERVER_NODE "taskset -c 2 nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    
    # Start sending requests
    taskset -c 2 go run request_sender.go $size
    
    # Move files for postprocessing
    mv $OW_DIRECTORY/GCScheduler/go_response_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Go/$size/client_time.txt"
    mv $OW_DIRECTORY/GCScheduler/go_server_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Go/$size/server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/go_heap_memory.log ../Graphs/GCScheduler/Go/$size/memory.txt
    # SCP the server.log file along with other files
    scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/server.log "../Graphs/GCScheduler/Go/$size/server.log"
    # Remove file after retrieving
    ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/LoadBalancer/*.log"
}

# Loop through each size
for size in "${sizes[@]}"; do
    # Kill the load balancer process if running
    curl $KILL_SERVER_API
    # Execute experiment
    send_requests $size
    # Kill the load balancer process if running
    curl $KILL_SERVER_API

    # Plot time responses
    python ../Graphs/GCScheduler/response_time_plotter.py "../Graphs/GCScheduler/Go/${size}/client_time.txt" "../Graphs/GCScheduler/Go/${size}/server_time.txt" "../Graphs/GCScheduler/Go/${size}/graph.pdf"
    # Plot memory patterns
    python ../Graphs/GCScheduler/go_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Go/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Go/${size}/memory.pdf"
done

# CSV post processing
filepaths=()
# Loop through each size and add its file path to the array
for size in "${sizes[@]}"; do
    filepaths+=("../Graphs/GCScheduler/Go/$size/latencies.csv")
done
# Combine the CSV files
(head -n 1 "${filepaths[0]}" && tail -n +2 -q "${filepaths[@]}") > "../Graphs/GCScheduler/Go/latencies.csv"
