OW_SERVER_NODE="am_CU@node0"
NATIVE_JAVA_API="http://128.110.96.62:9876/jsonresponse"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
ITERATIONS=10000

# Kill java server IF IT IS RUNNING
# Find the process ID
PID=$(ssh $OW_SERVER_NODE "jps | awk '/JsonServer/ {print \$1}'")

# Check if PID was found
if [ -z "$PID" ]; then
    echo "JsonServer is not running."
else
    # Kill the process
    ssh $OW_SERVER_NODE "kill $PID"
    echo "Killed JsonServer with PID $PID."
fi

MaxGCPauseMillis_values=(50 100 150 200 250 300)
Xmx_values=("64m" "128m" "256m" "512m" "1g" "2g" "4g")

# Iterate over MaxGCPauseMillis and Xmx values
for current_MaxGCPauseMillis in "${MaxGCPauseMillis_values[@]}"; do
    for current_Xmx in "${Xmx_values[@]}"; do
        GC_FLAGS="-Xmx$current_Xmx -XX:MaxGCPauseMillis=$current_MaxGCPauseMillis -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/PureJava/gc_log_$current_Xmx_$current_MaxGCPauseMillis"

        # start server
        ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/PureJava/; taskset -c 1 java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer > /users/am_CU/openwhisk-devtools/docker-compose/PureJava/server_log 2>&1 &"

        # Warm up until server is read to serve requests
        while :; do
            # Send request and store response
            RESPONSE=$(curl -s "$NATIVE_JAVA_API")

            # Check if response is valid (i.e., starts with '{')
            if [[ "$RESPONSE" == "{"* ]]; then
                echo "Received valid response!"
                # echo "$RESPONSE"
                break
            else
                # echo "$RESPONSE"
                echo "Invalid response, retrying..."
            fi

            # Optional: Sleep for a short duration before the next request
            sleep 1
        done

        # Determine peak throughput
        result=$(python load_gen.py)
        median_throughput=$(echo "$result" | grep "Median Throughput:" | awk '{print $3}')
        echo "Calculating median throughput: $median_throughput"
        RATE=$(echo "$median_throughput * 0.75" | bc)

        # Start experiments
        bash ParallelExperiment.sh $NATIVE_JAVA_API NativeJava $ITERATIONS $RATE
        # Enable file flush
        sleep 1
        # Move file for postprocessing
        mv NativeJavaOutputTime.txt ../Graphs/LoadTesting/Time_Xmx${current_Xmx}_MaxGCPauseMillis${current_MaxGCPauseMillis}.txt

        # Kill java server for next experiment
        # Find the process ID
        PID=$(ssh $OW_SERVER_NODE "jps | awk '/JsonServer/ {print \$1}'")

        # Check if PID was found
        if [ -z "$PID" ]; then
            echo "JsonServer is not running."
        else
            # Kill the process
            ssh $OW_SERVER_NODE "kill $PID"
            echo "Killed JsonServer with PID $PID."
        fi
    done
done

python ../Graphs/load_test_plotter.py