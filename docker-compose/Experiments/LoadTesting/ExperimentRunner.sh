# JS throughput testing
OW_SERVER_NODE="am_CU@node0"

ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker rmi -f $(docker images -aq)"
ssh $OW_SERVER_NODE "docker rmi -f my-js-server"
ssh $OW_SERVER_NODE "docker build -t js-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/JS/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d --rm --name my-js-server -p 8801:8800 js-server-image"
sleep 5
curl "http://node0:8801/JS?seed=999&arraysize=99&requestnumber=567"
sleep 1
locust --config=master.conf

exit 1
# ---------------------------------------------------------------
# Go throughput testing
OW_SERVER_NODE="am_CU@node0"

ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker rmi -f $(docker images -aq)"
ssh $OW_SERVER_NODE "docker rmi -f my-go-server"
ssh $OW_SERVER_NODE "docker build -t go-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Go/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 -d  --rm --name my-go-server -p 9501:9500 go-server-image"
sleep 5
curl "http://node0:9501/GoNative?seed=999&arraysize=99&requestnumber=567"
sleep 1
locust --config=master.conf

# go run request_sender.go 10000
# python Graphs/response_time_plotter.py go_response_times.txt go_server_times.txt go_heap_memory.log distribution.pdf latency.pdf latency_1.pdf sla.pdf
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf

# # Constants and Variables

# JAVA_API="http://128.110.96.59:8180/java"
# GO_API="http://128.110.96.59:8180/go"
# KILL_SERVER_API="http://128.110.96.59:8180/exitCall"
# OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose/Experiments"
# JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
# GO_RESPONSE_TIMES_FILE="go_response_times.txt"

# # Send request and measure request response latencies
# send_requests() {
#     local size=$1

#     # compile the docker images
#     ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/../Native/Go/; docker build -t go-server-image ."

#     # Change fakerequestarraysize
#     ssh am_CU@node0 "sed -i 's/fakeRequestArraySize = [^ ]*/fakeRequestArraySize = $size/' /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go"

#     # # Restart docker for good measure
#     # ssh $OW_SERVER_NODE "sudo systemctl restart docker"

#     # Restart the load balancer
#     ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
#     sleep 5
#     # Start sending requests
#     taskset -c 2 go run request_sender.go $size
    
#     # Move files for postprocessing
#     mv $OW_DIRECTORY/GCScheduler/go_response_times.txt "$OW_DIRECTORY/GCScheduler/Graphs/GCScheduler/Go/$size/client_time.txt"
#     mv $OW_DIRECTORY/GCScheduler/go_server_times.txt "$OW_DIRECTORY/GCScheduler/Graphs/GCScheduler/Go/$size/server_time.txt"
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/GCScheduler/Go/$size/memory.txt"
#     # SCP the server.log file along with other files
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/GCScheduler/Go/$size/server.log"
#     # Remove file after retrieving
#     ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/../LoadBalancer/*.log"

#     # Comment out Java part for now
#     # ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
#     # mv $OW_DIRECTORY/GCScheduler/java_response_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Java/$size/client_time.txt"
#     # mv $OW_DIRECTORY/GCScheduler/java_server_times.txt "$OW_DIRECTORY/Graphs/GCScheduler/Java/$size/server_time.txt"
#     # scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/java_heap_memory.log ../Graphs/GCScheduler/Java/$size/memory.txt
# }


# curl $KILL_SERVER_API
# # Restart the load balancer
# ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
# sleep 10
# # Start sending requests
# locust --config=master.conf
# mv /users/am_CU/openwhisk-devtools/docker-compose/Experiments/LoadTesting/locust_stats_history.csv /users/am_CU/openwhisk-devtools/docker-compose/Experiments/LoadTesting/EM_locust_stats_history.csv
# curl $KILL_SERVER_API
# # Restart the load balancer
# ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
# sleep 10
# # Start sending requests
# locust --config=master.conf
# # Plot time responses
# python ./loadbalancerloadlatencyplot.py "./Graphs/GCScheduler/Go/${size}/client_time.txt" "./Graphs/GCScheduler/Go/${size}/server_time.txt" "./Graphs/GCScheduler/Go/${size}/memory.txt" "./Graphs/GCScheduler/Go/${size}/distribution.pdf" "./Graphs/GCScheduler/Go/${size}/latency.pdf"
