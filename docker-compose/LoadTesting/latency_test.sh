# Commnad:
# taskset -c 1 locust -f load_gen.py --headless --users 10000 --spawn-rate 2000 -H http://128.110.96.62 --run-time 60 --csv=locust
# TODO: Extend script with peak bandwidth calculation

# Constants and Variables
OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.176:9876/jsonresponse"
GO_API="http://128.110.96.176:9875/GoNative"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
ITERATIONS=100000
MaxGCPauseMillis=50
Xmx="256m"
GC_FLAGS="-Xms$Xmx -Xmx$Xmx -XX:MaxGCPauseMillis=$MaxGCPauseMillis -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/Native/Java/gc_log"
echo $GC_FLAGS

start_java_server() {
# compile the server
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java"
# start Java server
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; taskset -c 1 java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer > /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/server_log 2>&1 &"
}

# Kill java server IF IT IS RUNNING
kill_java_server() {
    local PID=$(ssh $OW_SERVER_NODE "jps | awk '/JsonServer/ {print \$1}'")
    if [ -z "$PID" ]; then
        echo "JavaServer is not running."
    else
        ssh $OW_SERVER_NODE "kill -9 $PID"
        echo "Killed JsonServer with PID $PID."
    fi
}

# start Go server
start_go_server() {
ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Go/; taskset -c 2 go run server.go &"
}

# Function to kill Go server
kill_go_server() {
    local PID=$(ssh $OW_SERVER_NODE "pgrep -f server.go")  # Replace 'server.go' with the Go server's process name if different
    if [ -z "$PID" ]; then
        echo "Go server is not running."
    else
        ssh $OW_SERVER_NODE "kill -9 $PID"
        echo "Killed Go server with PID $PID."
    fi
}

# Function to warm up server
warm_up_server() {
    local API_URL=$1
    local RETRY_COUNT=0
    local MAX_RETRIES=5

    while [ "$RETRY_COUNT" -lt "$MAX_RETRIES" ]; do
        HTTP_STATUS=$(curl -o /dev/null -s -w '%{http_code}' "$API_URL")
        if [ "$HTTP_STATUS" -eq 200 ]; then
            echo "Received valid response (HTTP Status 200) from $API_URL"
            return 0
        else
            echo "Invalid response from $API_URL, retrying... (Attempt $((RETRY_COUNT+1)) of $MAX_RETRIES)"
            ((RETRY_COUNT++))
            sleep 1
        fi
    done

    echo "Failed to receive a valid response from $API_URL after $MAX_RETRIES attempts. Exiting script."
    exit 1
}

# Kill any previous running instances of server
kill_java_server
start_java_server
# Warm up servers
warm_up_server "$JAVA_API"
# Java Load Processing
export API_URL=$JAVA_API
taskset -c 3 locust --config=./master.conf
# Enable file flush
sleep 1
# Move file for postprocessing
mv locust_stats_history.csv ../Graphs/LoadTesting/Java/LoadLatencyCurve.csv
# TODO: From this file, determine the median, p99, p999 request latency 
# TODO: From this file, determine peak throughput supported by the server
# Kill server after execution
kill_java_server

# Go Processing
kill_go_server
start_go_server
warm_up_server "$GO_API"
# Go Load Processing
export API_URL=$GO_API
taskset -c 3 locust --config=./master.conf
# Enable file flush
sleep 1
# Move file for postprocessing
mv locust_stats_history.csv ../Graphs/LoadTesting/Go/LoadLatencyCurve.csv
kill_go_server


# Start experiments
# bash ParallelExperiment.sh $JAVA_API NativeJava $ITERATIONS $RATE
# locust --headless --users $RATE --spawn-rate $RATE -H http://128.110.96.176 --run-time 60s --csv=locust -f load_gen.py --skip-log --reset-stats
# taskset -c 4 locust --config=./worker.conf &
# taskset -c 5 locust --config=./worker.conf &
# taskset -c 6 locust --config=./worker.conf &
# taskset -c 7 locust --config=./worker.conf &
# taskset -c 8 locust --config=./worker.conf &
