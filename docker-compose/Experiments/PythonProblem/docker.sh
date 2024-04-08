docker build -t python-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Python/
docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-python-server -p 9101:9100 python-server-image

sleep 5
curl "http://node0:9101/Python?seed=1000&arraysize=10000&requestnumber=56"
sleep 1
go run request_sender.go 1000
python SLAplotter.py ./go_response_times.txt SLAPlot.pdf