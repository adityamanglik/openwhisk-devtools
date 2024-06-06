#!/bin/bash

OW_SERVER_NODE="am_CU@node0"
GC_COLLECTORS=("-XX:+UseSerialGC" "-XX:+UseParallelGC" "-XX:+UseZGC")
MEMORY_SIZES=("128m" "512m" "10240m")

for GC in "${GC_COLLECTORS[@]}"; do
  for MEM_SIZE in "${MEMORY_SIZES[@]}"; do
    echo "Running experiment with GC: $GC and Memory Size: $MEM_SIZE"

    # SSH into the server node and clean up any existing Docker containers/images
    ssh $OW_SERVER_NODE "docker rm -vf \$(docker ps -aq)"
    ssh $OW_SERVER_NODE "docker rmi -f \$(docker images -aq)"
    ssh $OW_SERVER_NODE "docker stop my-java-server"

    # Copy the Dockerfile to the target directory
    ssh $OW_SERVER_NODE "cp /users/am_CU/openwhisk-devtools/docker-compose/Experiments/JavaProblem/Dockerfile /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/Dockerfile"

    # Build the Docker image with the current GC flags and memory sizes
    ssh $OW_SERVER_NODE "docker build --build-arg GC_FLAGS='$GC -Xms$MEM_SIZE -Xmx$MEM_SIZE' -t java-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/"

    # Run the Docker container
    ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory='$MEM_SIZE' -d --rm --name my-java-server -p 8601:8600 java-server-image"

    # Allow some time for the server to start
    sleep 5

    # Make a sample request to the server
    curl "http://node0:8601/jsonresponse?seed=999&arraysize=99&requestnumber=567"

    # Allow some time for the server to process the request
    sleep 1

    # Run the locust test
    locust --config=master.conf

    # Run the analysis script
    python analysis.py >> Results/"res_$GC_$MEM_SIZE.txt"

    # Optional: collect and store results for each GC and memory size run
    # mv locust_stats_history.csv locust_stats_history_${GC}_${MEM_SIZE}.csv

  done
done
