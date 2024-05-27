docker rm -vf $(docker ps -aq)
docker rmi -f $(docker images -aq)
docker build -t go-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Go/
docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-go-server -p 9501:9500 go-server-image
sleep 5
curl "http://node0:9501/GoNative?seed=999&arraysize=99&requestnumber=567"
sleep 1
go run request_sender.go 10000
python Graphs/response_time_plotter.py go_response_times.txt go_server_times.txt go_heap_memory.log distribution.pdf latency.pdf latency_1.pdf sla.pdf
# python SLAplotter.py ./go_response_times.txt SLAPlot.pdf