OW_SERVER_NODE="am_CU@node0"

# Kill container
# ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
# Start container
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sed -i 's/constant_pacing(1)/constant_pacing(2)/g' loadlatency.py
rm ./*.csv
locust --config=master.conf -u 1
cp ./locust_stats.csv ./Graphs/RequestArrivalRate/_2.csv
# ########################################################################

# Kill container
# ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
# Start container
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sed -i 's/constant_pacing(1)/constant_pacing(2)/g' loadlatency.py
rm ./*.csv
locust --config=master.conf -u 10 -r 10
cp ./locust_stats.csv ./Graphs/RequestArrivalRate/_20.csv
# Backup ########################################################################
# Kill container
# ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
# Start container
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sed -i 's/constant_pacing(2)/constant_pacing(10)/g' loadlatency.py
rm ./*.csv
locust --config=master.conf -u 25 -r 25
cp ./locust_stats.csv ./Graphs/RequestArrivalRate/_50.csv
# Backup ########################################################################
# Kill container
# ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
# Start container
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sed -i 's/constant_pacing(10)/constant_pacing(30)/g' loadlatency.py
rm ./*.csv
locust --config=master.conf -u 50 -r 50
cp ./locust_stats.csv ./Graphs/RequestArrivalRate/_100.csv
# Backup ########################################################################
# Kill container
# ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
# Start container
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sed -i 's/constant_pacing(30)/constant_pacing(60)/g' loadlatency.py
rm ./*.csv
locust --config=master.conf -u 100 -r 100
cp ./locust_stats.csv ./Graphs/RequestArrivalRate/_200.csv
# Backup ########################################################################

# # Constants and Variables
# 
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
#     mv $OW_DIRECTORY/RequestArrivalRate/go_response_times.txt "$OW_DIRECTORY/RequestArrivalRate/Graphs/RequestArrivalRate/Go/$size/client_time.txt"
#     mv $OW_DIRECTORY/RequestArrivalRate/go_server_times.txt "$OW_DIRECTORY/RequestArrivalRate/Graphs/RequestArrivalRate/Go/$size/server_time.txt"
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/go_heap_memory.log "./Graphs/RequestArrivalRate/Go/$size/memory.txt"
#     # SCP the server.log file along with other files
#     scp $OW_SERVER_NODE:$OW_DIRECTORY/../LoadBalancer/server.log "./Graphs/RequestArrivalRate/Go/$size/server.log"
#     # Remove file after retrieving
#     ssh $OW_SERVER_NODE "rm $OW_DIRECTORY/../LoadBalancer/*.log"

#     # Comment out Java part for now
#     # ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
#     # mv $OW_DIRECTORY/RequestArrivalRate/java_response_times.txt "$OW_DIRECTORY/Graphs/RequestArrivalRate/Java/$size/client_time.txt"
#     # mv $OW_DIRECTORY/RequestArrivalRate/java_server_times.txt "$OW_DIRECTORY/Graphs/RequestArrivalRate/Java/$size/server_time.txt"
#     # scp $OW_SERVER_NODE:$OW_DIRECTORY/LoadBalancer/java_heap_memory.log ../Graphs/RequestArrivalRate/Java/$size/memory.txt
# }


# curl $KILL_SERVER_API
# # Restart the load balancer
# ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/loadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
# sleep 10
# # Start sending requests
# locust --config=master.conf
# # Move file to recording
# mv /users/am_CU/openwhisk-devtools/docker-compose/Experiments/RequestArrivalRate/locust_stats_history.csv /users/am_CU/openwhisk-devtools/docker-compose/Experiments/RequestArrivalRate/EM_locust_stats_history.csv
# curl $KILL_SERVER_API
# # Restart the load balancer
# ssh $OW_SERVER_NODE "nohup go run /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/NOGCloadbalancer.go > /users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/server.log 2>&1 &"
# sleep 10
# # Start sending requests
# locust --config=master.conf
# # Plot time responses
# python ./loadbalancerloadlatencyplot.py "./Graphs/RequestArrivalRate/Go/${size}/client_time.txt" "./Graphs/RequestArrivalRate/Go/${size}/server_time.txt" "./Graphs/RequestArrivalRate/Go/${size}/memory.txt" "./Graphs/RequestArrivalRate/Go/${size}/distribution.pdf" "./Graphs/RequestArrivalRate/Go/${size}/latency.pdf"



# # Pure docker version code
# docker rm -vf $(docker ps -aq)
# docker rmi -f $(docker images -aq)
# cp /users/am_CU/openwhisk-devtools/docker-compose/Experiments/JavaProblem/Dockerfile /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/Dockerfile
# docker build -t java-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/
# docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image
# sleep 5
# curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
# sleep 1
# go run request_sender.go 100000
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf