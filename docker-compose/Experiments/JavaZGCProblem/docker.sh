OW_SERVER_NODE="am_CU@node0"
ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker rmi -f $(docker images -aq)"
ssh $OW_SERVER_NODE "docker stop my-java-server"
ssh $OW_SERVER_NODE "cp /users/am_CU/openwhisk-devtools/docker-compose/Experiments/JavaZGCProblem/Dockerfile /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/Dockerfile"
ssh $OW_SERVER_NODE "docker build --build-arg GC_FLAGS='-XX:+UseZGC -Xms64m -Xmx64m' -t java-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d --rm --name my-java-server -p 8601:8600 java-server-image"
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
sleep 1
# locust --config=master.conf
# python analysis.py
go run request_sender.go 10000
python Graphs/response_time_plotter.py go_response_times.txt go_server_times.txt go_heap_memory.log distribution.pdf latency.pdf latency_1.pdf sla.pdf
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf