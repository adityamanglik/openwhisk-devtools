# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.167:8180/java?seed=5"
GO_API="http://128.110.96.167:8180/go?seed=5"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
JAVA_RESPONSE_TIMES_FILE="java_response_times.txt"
GO_RESPONSE_TIMES_FILE="go_response_times.txt"
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
    local api_url=$1
    local response_time_file=$2
    local execution_time_file=$3

    # Update the Java code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Native/Java/Hello.java"

    local size=$1
    ssh $OW_SERVER_NODE "awk '/MARKER_FOR_SIZE_UPDATE/{print;getline;print \"const ARRAY_SIZE = \" size \";\";next}1' size=${size} $OW_DIRECTORY/Functions/hello.go > $OW_DIRECTORY/Functions/temp.go && mv $OW_DIRECTORY/Functions/temp.go $OW_DIRECTORY/Functions/hello.go"

    for i in $(seq 1 $ITERATIONS)
    do
        # Measure the response time and capture the response
        start_time=$(date +%s.%N)
        response=$(curl -s "$api_url")
        end_time=$(date +%s.%N)
        elapsed=$(echo "scale=3; ($end_time - $start_time) * 1000" | bc)

        # Parse the response to extract executionTime in milliseconds
        execution_time=$(echo $response | jq -r '.executionTime')

        # Convert execution time to milliseconds if needed
        execution_time_ms=$(echo "scale=3; $execution_time / 1000000" | bc)

        # Record the total elapsed time
        echo "$elapsed" >> $response_time_file

        # Record the extracted executionTime in milliseconds
        echo "$execution_time_ms" >> $execution_time_file
    done
}

send_requests $JAVA_API "java_response_times.txt" "java_execution_times_ms.txt"
send_requests $GO_API "go_response_times.txt" "go_execution_times_ms.txt"

# BACKUP #######################################################################
# Locust code
# export API_URL=$GO_API
# locust --config=./master.conf

# export API_URL=$JAVA_API
# locust --config=./master.conf
