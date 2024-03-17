docker build -t java-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/
docker run --cpuset-cpus 4 --memory=128m -d  --rm --name my-java-server -p 8601:8600 java-server-image
sleep 5
curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"
sleep 1
go run request_sender.go 100000