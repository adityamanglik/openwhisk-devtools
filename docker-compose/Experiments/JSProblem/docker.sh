OW_SERVER_NODE="am_CU@node0"
ssh $OW_SERVER_NODE "docker rm -vf $(docker ps -aq)"
ssh $OW_SERVER_NODE "docker rmi -f $(docker images -aq)"
ssh $OW_SERVER_NODE "docker stop my-js-server"
ssh $OW_SERVER_NODE "docker build -t js-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/JS/"
ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-js-server -p 8801:8800 js-server-image"
sleep 5
curl "http://node0:8801/JS?seed=1000&arraysize=10000&requestnumber=56" >> Results/"res_128.txt"
sleep 1
locust --config=master.conf
echo "" >> Results/"res_128.txt"
python analysis.py >> Results/"res_128.txt"
# go run request_sender.go 10000
# python Graphs/response_time_plotter.py go_response_times.txt go_server_times.txt go_heap_memory.log distribution.pdf latency.pdf latency_1.pdf sla.pdf
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf