# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
GO_API="http://128.110.96.59:8180/go"
KILL_SERVER_API="http://128.110.96.59:8180/exitCall"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GOGCSweep"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# Array of sizes
sizes=(100000 1000000)
# Array of GOGC values
GOGC=(-1 1 10 100 500 999)

# Send request and measure request response latencies
send_requests() {
    local size=$1
    # # Restart docker for good measure
    # ssh $OW_SERVER_NODE "sudo systemctl restart docker"
    # sleep 10
    # Loop through each GOGC value
    for gc in "${GOGC[@]}"; do
        # Kill the load balancer process if running
        curl $KILL_SERVER_API

        # Change GOGC value in Dockerfile
        ssh am_CU@node0 "sed -i 's/ENV GOGC=.*/ENV GOGC=$gc/' /users/am_CU/openwhisk-devtools/docker-compose/Native/Go/Dockerfile"

        # compile the docker images
        ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/../../Native/Go/; docker build -t go-server-image ."

        # Restart the load balancer
        ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    
        # Start sending requests
        taskset -c 2 go run request_sender.go $size $gc

        # Extract required data from file

        # Remove files to prevent data mix
        rm ./*.txt
        rm ./*.log
        
        # Move files for postprocessing
        scp $OW_SERVER_NODE:$OW_DIRECTORY/../../LoadBalancer/go_heap_memory.log "$OW_DIRECTORY/Data/$size_$gc_memory.txt"
        # mv $OW_DIRECTORY/Experiments/GOGCSweep/go_response_times.txt "$OW_DIRECTORY/Experiments/GOGCSweep/Data/$size_client_time.txt"
        # mv $OW_DIRECTORY/Experiments/GOGCSweep/go_server_times.txt "$OW_DIRECTORY/Experiments/GOGCSweep/Data/$size_server_time.txt"
        # scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/go_heap_memory.log ../Graphs/GCScheduler/Go/$size/memory.txt
        # SCP the server.log file along with other files
        # scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/server.log "../Graphs/GCScheduler/Go/$size/server.log"
        # Remove file after retrieving
        ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/../../LoadBalancer/*.log"

        # Kill the load balancer process if running
        curl $KILL_SERVER_API
    done
}

# Loop through each size
for size in "${sizes[@]}"; do
    
    # Execute experiment
    send_requests $size

done

python $OW_DIRECTORY/plotter.py