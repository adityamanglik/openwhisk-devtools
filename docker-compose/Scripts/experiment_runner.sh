OW_SERVER_NODE="am_CU@node0"
JAVA_API="http://128.110.96.15:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world"
JAVASCRIPT_API="http://128.110.96.15:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"

function runJavaExperiment() {
    # Update the Java code with the new array size
    local size=$1    
    ssh $OW_SERVER_NODE "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Functions/Hello.java"

    # Make sure we get fresh data by resetting functions via update
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Functions/; javac -cp gson-2.10.1.jar Hello.java"
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Functions/; jar cvf hello.jar Hello.class"
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update helloJava Functions/hello.jar --main Hello"

    # Start generating load
    source Experiment.sh $JAVA_API Java

    # Retrieve warm/cold status of each activation
    scp Javaactivation_ids.txt $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; bash ./activation_status_checker.sh ./Javaactivation_ids.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/Javaactivation_ids.txt_startStates.txt ./ 

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/Java/$size/"

    # Java plotter
    python ../Graphs/java_response_time_plotter.py $size
}

function runJSExperiment() {
    # Update the JavaScript code with the new array size
    local size=$1
    ssh $OW_SERVER_NODE "sed -i 's/const ARRAY_SIZE = [0-9]\+;/const ARRAY_SIZE = ${size};/' $OW_DIRECTORY/Functions/wordcount.js"

    # Make sure we get fresh data by resetting functions via update
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update hello Functions/wordcount.js"

    # Start generating load
    source Experiment.sh $JAVASCRIPT_API JS

    # Retrieve warm/cold status of each activation
    scp JSactivation_ids.txt $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; bash ./activation_status_checker.sh ./JSactivation_ids.txt"
    scp $OW_SERVER_NODE:$OW_DIRECTORY/Scripts/JSactivation_ids.txt_startStates.txt ./ 

    # Move all log and image files to that directory
    mv $OW_DIRECTORY/Scripts/*.txt "$OW_DIRECTORY/Graphs/JS/$size/"

    # JS Plotter
    python ../Graphs/js_response_time_plotter.py $size
}

# Run the experiments for the three array sizes
for size in 100 10000 1000000 5000000 ; do
    runJavaExperiment $size
    runJSExperiment $size
done
