OW_SERVER_NODE="am_CU@node0"
ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker rmi -f $(docker images -aq)"
ssh $OW_SERVER_NODE "docker stop my-python-server"
ssh $OW_SERVER_NODE "docker build --build-arg MEMORY_LIMIT=134217728 -t python-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Python/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-python-server -p 9901:9900 python-server-image"
sleep 5
curl "http://node0:9901/Python?seed=999&arraysize=5&requestnumber=567"
sleep 1
# locust --config=master.conf
go run request_sender.go
python Graphs/response_time_plotter.py go_response_times.txt go_server_times.txt go_heap_memory.log distribution.pdf latency.pdf latency_1.pdf sla.pdf
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf