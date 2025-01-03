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

    # compile the docker images
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/../Native/Go/; docker build -t go-server-image ."

    # # Restart docker for good measure
    # ssh $OW_SERVER_NODE "sudo systemctl restart docker"

    # Change fakerequestarraysize
    ssh am_CU@node0 "sed -i 's/fakeRequestArraySize = [^ ]*/fakeRequestArraySize = $size/' /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go"


    # Restart the load balancer
    ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
    sleep 5
    # Start sending requests
    taskset -c 2 go run request_sender.go $size
    
    # Move files for postprocessing
    mv $OW_DIRECTORY/GCScheduler/go_response_times.txt "$OW_DIRECTORY/GCScheduler/Graphs/GCScheduler/Go/$size/client_time.txt"
    mv $OW_DIRECTORY/GCScheduler/go_server_times.txt "$OW_DIRECTORY/GCScheduler/Graphs/GCScheduler/Go/$size/server_time.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/GCScheduler/Go/$size/memory.txt"
    # SCP the server.log file along with other files
    scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/GCScheduler/Go/$size/server.log"
    # Remove file after retrieving
    ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/../LoadBalancer/*.log"

    # Comment out Java part for now
    # ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
    # mv $OW_DIRECTORY/GCScheduler/java_response_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Java/$size/client_time.txt"
    # mv $OW_DIRECTORY/GCScheduler/java_server_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Java/$size/server_time.txt"
    # scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/java_heap_memory.log ../Graphs/GCScheduler/Java/$size/memory.txt
}

# Array of sizes
# sizes=(100 100 100 1000 1000 1000 10000 10000 10000 50000 50000 50000 100000 100000 100000 1000000 1000000 1000000)
# sizes=(10000 10000 10000 10000 10000)
sizes=(100000)

# for size in "${sizes[@]}"; do
#     # python ../Graphs/GCScheduler/response_time_plotter.py "../Graphs/GCScheduler/Java/${size}/client_time.txt" "../Graphs/GCScheduler/Java/${size}/server_time.txt" "../Graphs/GCScheduler/Java/${size}/graph.pdf"

#     python ./Graphs/GCScheduler/response_time_plotter.py "./Graphs/GCScheduler/Go/${size}/client_time.txt" "./Graphs/GCScheduler/Go/${size}/server_time.txt" "./Graphs/GCScheduler/Go/${size}/memory.txt" "./Graphs/GCScheduler/Go/${size}/NOGCdistribution.pdf" "./Graphs/GCScheduler/Go/${size}/NOGClatency.pdf"
# done
# echo "Graphs done"

# for size in "${sizes[@]}"; do
#     python ../Graphs/GCScheduler/go_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Go/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Go/${size}/memory.pdf"
#     python ../Graphs/GCScheduler/java_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Java/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Java/${size}/memory.pdf"
# done

# Loop through each size
for size in "${sizes[@]}"; do
    # Kill the load balancer process if running
    curl $KILL_SERVER_API
    send_requests $size
    # Kill the load balancer process if running
    curl $KILL_SERVER_API

    # Plot time responses
    # python ./Graphs/GCScheduler/response_time_plotter.py "./Graphs/GCScheduler/Java/${size}/client_time.txt" "./Graphs/GCScheduler/Java/${size}/server_time.txt" "./Graphs/GCScheduler/Java/${size}/graph.pdf"
    python ./Graphs/GCScheduler/response_time_plotter.py "./Graphs/GCScheduler/Go/${size}/client_time.txt" "./Graphs/GCScheduler/Go/${size}/server_time.txt" "./Graphs/GCScheduler/Go/${size}/memory.txt" "./Graphs/GCScheduler/Go/${size}/NOGCdistribution.pdf" "./Graphs/GCScheduler/Go/${size}/NOGClatency.pdf"
    
    # Plot memory patterns
    python ./Graphs/GCScheduler/go_mem_plotter.py "./Graphs/GCScheduler/Go/${size}/memory.txt" "./Graphs/GCScheduler/Go/${size}/NOGCmemory.pdf"
    # python ./Graphs/GCScheduler/java_mem_plotter.py "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Java/${size}/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/GCScheduler/Java/${size}/memory.pdf"

    # Calculate impact of GC
    # python ./Graphs/GCScheduler/analyzer.py "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/server_time.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/client_time.txt" >> analyzer3.log
    python ./Graphs/GCScheduler/analyzer.py "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/server_time.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/client_time.txt"
    # tail -n 1 "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/$size/latencies.csv" >> "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/latencies.csv"
    # echo "$size" >> analyzer3.log
done

# CSV post processing

# Initialize an empty array to hold file paths
# filepaths=()

# # Loop through each size and add its file path to the array
# for size in "${sizes[@]}"; do
#     filepaths+=("./Graphs/GCScheduler/Go/$size/latencies.csv")
# done

# # Combine the CSV files
# (head -n 1 "${filepaths[0]}" && tail -n +2 -q "${filepaths[@]}") > "./Graphs/GCScheduler/Go/latencies.csv"

# BACKUP #######################################################################

# send_requests $JAVA_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/GCScheduler/response_time_plotter.py ../Graphs/GCScheduler/Java/1000000/client_time.txt ../Graphs/GCScheduler/Java/1000000/server_time.txt ../Graphs/GCScheduler/Java/1000000/graph.pdf

# send_requests $GO_API "client_time.txt" "server_time.txt" 1000000
# python ../Graphs/GCScheduler/response_time_plotter.py ../Graphs/GCScheduler/Go/1000000/client_time.txt ../Graphs/GCScheduler/Java/1000000/server_time.txt ../Graphs/GCScheduler/Java/1000000/graph.pdf


# Locust code
# export API_URL=$GO_API
# locust --config=./master.conf

# export API_URL=$JAVA_API
# locust --config=./master.conf