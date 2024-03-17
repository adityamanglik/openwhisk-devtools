docker build -t java-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/
docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image
sleep 5
curl "http://node0:8601/jsonresponse?seed=1000&arraysize=10000&requestnumber=56"
sleep 1
go run request_sender.go 50000