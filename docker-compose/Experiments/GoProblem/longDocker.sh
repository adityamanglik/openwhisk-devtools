# Go
OW_SERVER_NODE="am_CU@node0"
MEMORY_SIZES=("128MiB" "512MiB" "10240MiB")
GOGC_VALUES=("1" "1000")

for GOGC in "${GOGC_VALUES[@]}"; do
  for MEM_SIZE in "${MEMORY_SIZES[@]}"; do
      echo "Running experiment Memory Size: $MEM_SIZE"

      # SSH into the server node and clean up any existing Docker containers/images
      ssh $OW_SERVER_NODE "docker rm -vf \$(docker ps -aq)"
      ssh $OW_SERVER_NODE "docker rmi -f \$(docker images -aq)"
      ssh $OW_SERVER_NODE "docker stop my-go-server"
      
      # Build the Docker image with the current GC flags and memory sizes
      ssh $OW_SERVER_NODE "docker build --build-arg GOGC=$GOGC --build-arg GOMEMLIMIT=$MEM_SIZE -t go-server-image /users/am_CU/openwhisk-devtools/docker-compose/Native/Go/"

      # Run the Docker container
      ssh $OW_SERVER_NODE "docker run --cpuset-cpus 4 --memory='$MEM_SIZE' -d --rm --name my-go-server -p 9501:9500 go-server-image"

      # Allow some time for the server to start
      sleep 5

      # Make a sample request to the server
      curl "http://node0:9501/GoNative?seed=999&arraysize=99&requestnumber=567" >> Results/"res_'$GOGC'_'$MEM_SIZE'.txt"

      # Allow some time for the server to process the request
      sleep 1

      # Run the locust test
      locust --config=master.conf

      # Run the analysis script
      echo "" >> Results
      python analysis.py >> Results/"res_'$GOGC'_'$MEM_SIZE'.txt"

      # Optional: collect and store results for each GC and memory size run
      mv locust_stats_history.csv Results/stats_${GOGC}_${MEM_SIZE}.csv

    done
  done