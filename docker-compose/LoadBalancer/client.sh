# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
LOADBALANCER_API="http://128.110.96.167:8180/"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
RESPONSE_TIMES_FILE="response_times.txt"
ITERATIONS=10000

# Build docker images
build_docker_images() {
# compile the docker images
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; docker build -t java-server-image ."
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; docker build -t go-server-image ."
}

# Client code
# build_docker_images

# Send request and measure request response latencies
send_requests() {
    for i in $(seq 1 $ITERATIONS)
    do
        # Measure the response time
        start_time=$(date +%s.%N)
        curl -s "$LOADBALANCER_API" > /dev/null
        end_time=$(date +%s.%N)
        elapsed=$(echo "scale=3; $end_time - $start_time" | bc)

        # Record the response time
        echo "$elapsed" >> $RESPONSE_TIMES_FILE
    done
}

send_requests

# BACKUP #######################################################################
# Locust code
# export API_URL=$GO_API
# locust --config=./master.conf

# export API_URL=$JAVA_API
# locust --config=./master.conf
