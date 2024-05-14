# GO_API="http://node0:8180/go"

# Start container on node0 with specified memory allocation
ssh $OW_SERVER_NODE "docker build -t go-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Go/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-go-server -p 9501:9500 go-server-image"
# Send traffic and record timings
sleep 5
curl "http://node0:9101/GoNative?seed=1000&arraysize=10000&requestnumber=56"
go run request_sender.go 10000

# Plot timings in SLA plot


# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf

# # Constants and Variables
# OW_SERVER_NODE="am_CU@node0"

# KILL_SERVER_API="http://node0:8180/exitCall"
# OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose/Experiments"
# JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
# GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# # Send request and measure request response latencies
# send_requests() {
#     local size=$1

#     # Start the load balancer
#     ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go SingleServer > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
#     sleep 5
    
#     # Start sending requests
#     taskset -c 2 go run request_sender.go $size
    
#     # Move files for postprocessing
#     mv $OW_DIRECTORY/IntroProblemGraph/go_response_times.txt "$OW_DIRECTORY/IntroProblemGraph/Graphs/Go/$size/client_time.txt"
#     mv $OW_DIRECTORY/IntroProblemGraph/go_server_times.txt "$OW_DIRECTORY/IntroProblemGraph/Graphs/Go/$size/server_time.txt"
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/Go/$size/memory.txt"
#     # SCP the server.log file along with other files
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/Go/$size/server.log"

# }

# sizes=(10000)
# # sizes=(1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000 1000000)

# # Loop through each size
# for size in "${sizes[@]}"; do
#     # Kill the load balancer process if running
#     curl $KILL_SERVER_API

#     send_requests $size

#     # Kill the load balancer process if running
#     curl $KILL_SERVER_API

#     # Plot time responses
#     python ./Graphs/response_time_plotter.py "./Graphs/Go/${size}/client_time.txt" "./Graphs/Go/${size}/server_time.txt" "./Graphs/Go/${size}/memory.txt" "./Graphs/Go/${size}/distribution.pdf" "./Graphs/Go/${size}/latency.pdf" "./Graphs/Go/${size}/latency_1.pdf" "./Graphs/Go/${size}/sla_plot.pdf"
    
#     # Plot memory patterns
#     python ./Graphs/go_mem_plotter.py "./Graphs/Go/${size}/memory.txt" "./Graphs/Go/${size}/memory.pdf"

#     # Calculate impact of GC
#     # python ./Graphs/analyzer.py "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/IntroProblemGraph/Graphs/Go/$size/memory.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/IntroProblemGraph/Graphs/Go/$size/server_time.txt" "/users/am_CU/openwhisk-devtools/docker-compose/Experiments/IntroProblemGraph/Graphs/Go/$size/client_time.txt"
# done