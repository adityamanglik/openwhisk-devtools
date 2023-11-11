# set -x
OW_SERVER_NODE="am_CU@node0"
NATIVE_JAVA_API="http://128.110.96.62:9876/jsonresponse"
JAVA_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world"
JAVASCRIPT_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world"
GO_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloGo/world"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
GC_FLAGS="-Xmx64m -XX:MaxGCPauseMillis=50 -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/Native/Java/gc_log"
NO_GC_FLAGS="-Xmx4g -XX:MaxGCPauseMillis=500 -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/Native/Java/no_gc_log"
ITERATIONS=5000

ssh $OW_SERVER_NODE "export OW_DIRECTORY='/users/am_CU/openwhisk-devtools/docker-compose';"

function runJavaExperiment() {
    # Update the Java code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Functions/Hello.java"

    # Make sure we get fresh data by resetting functions via update
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Functions/; javac -cp gson-2.10.1.jar Hello.java"
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Functions/; jar cvf hello.jar Hello.class"
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update helloJava Functions/hello.jar --main Hello"

    # Start generating load
    source Experiment.sh $JAVA_API Java $ITERATIONS

    # Retrieve warm/cold status of each activation
    scp Javaactivation_ids.txt $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; bash ./activation_status_checker.sh ./Javaactivation_ids.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/Javaactivation_ids.txt_startStates.txt ./

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/Java/$size/"

    # Java plotter
    python ../Graphs/java_response_time_plotter.py $size
}

function runNativeJavaExperiment() {
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

    # Update the Java code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Native/Java/Hello.java"

    # compile code and start server
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java"
    ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; taskset -c 1 java -cp .:gson-2.10.1.jar $GC_FLAGS JsonServer > /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/server_log 2>&1 &"

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

    # Start generating load
    source Experiment.sh $NATIVE_JAVA_API NativeJava $ITERATIONS

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/NativeJava/$size/"

    # Java plotter
    python ../Graphs/native_java_response_time_plotter.py $size

    # Kill java server

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

}

function runNativeJavaNoGCExperiment() {
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

    # Update the Java code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Native/Java/Hello.java"

    # compile code and start server
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; javac -cp .:gson-2.10.1.jar Hello.java JsonServer.java"
    ssh -f $OW_SERVER_NODE "cd $OW_DIRECTORY/Native/Java/; taskset -c 1 java -cp .:gson-2.10.1.jar $NO_GC_FLAGS JsonServer > /users/am_CU/openwhisk-devtools/docker-compose/Native/Java/server_log 2>&1 &"

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
            echo "Invalid response, retrying..."
        fi

        # Optional: Sleep for a short duration before the next request
        sleep 1
    done

    # Start generating load
    source Experiment.sh $NATIVE_JAVA_API NativeJava $ITERATIONS

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/NativeJava/$size/"

    # Java plotter
    python ../Graphs/native_java_response_time_plotter.py $size

    # Kill java server

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

}

function runJSExperiment() {
    # Update the JavaScript code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/const ARRAY_SIZE = [0-9]\+;/const ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Functions/wordcount.js"

    # Make sure we get fresh data by resetting functions via update
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update hello Functions/wordcount.js"

    # Start generating load
    source Experiment.sh $JAVASCRIPT_API JS $ITERATIONS

    # Retrieve warm/cold status of each activation
    scp JSactivation_ids.txt $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; bash ./activation_status_checker.sh ./JSactivation_ids.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/JSactivation_ids.txt_startStates.txt ./

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/JS/$size/"

    # JS Plotter
    python ../Graphs/js_response_time_plotter.py $size
}

function runGoExperiment() {
    # Update the Go code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "awk '/MARKER_FOR_SIZE_UPDATE/{print;getline;print \"const ARRAY_SIZE = \" size \";\";next}1' size=${size} $OW_DIRECTORY/Functions/hello.go > $OW_DIRECTORY/Functions/temp.go && mv $OW_DIRECTORY/Functions/temp.go $OW_DIRECTORY/Functions/hello.go"

    # Make sure we get fresh data by resetting functions via update
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update helloGo Functions/hello.go"

    # Start generating load
    source Experiment.sh $GO_API Go $ITERATIONS

    # Retrieve warm/cold status of each activation
    scp Goactivation_ids.txt $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; bash ./activation_status_checker.sh ./Goactivation_ids.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/Goactivation_ids.txt_startStates.txt ./

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/Go/$size/"

    # Go plotter (Assuming you have a python plotter for Go as well. If not, remove the next line)
    python ../Graphs/go_response_time_plotter.py $size
}
# Backup size: 5000000
# Run the experiments for the three array sizes
for size in 1000000; do
    echo "Size: $size"
    runNativeJavaExperiment $size
    cp -r $OW_DIRECTORY/Graphs/NativeJava/* $OW_DIRECTORY/Graphs/NativeJavaWithGC/
    runNativeJavaNoGCExperiment $size
    # runJavaExperiment $size
    # runJSExperiment $size
    # runGoExperiment $size
    # python ../Graphs/js_response_time_plotter.py $size
    # python ../Graphs/java_response_time_plotter.py $size
    # python ../Graphs/go_response_time_plotter.py $size
    # python ../Graphs/native_java_response_time_plotter.py $size
done
