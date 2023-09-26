function updateArraySize() {
    local size=$1
    ssh am_CU@apt069.apt.emulab.net "sed -i 's/const ARRAY_SIZE = [0-9]\+;/const ARRAY_SIZE = ${size};/' /users/am_CU/openwhisk-devtools/docker-compose/Functions/wordcount.js"
    ssh am_CU@apt069.apt.emulab.net "sed -i 's/private static final int ARRAY_SIZE = [0-9]\+;/private static final int ARRAY_SIZE = ${size};/' /users/am_CU/openwhisk-devtools/docker-compose/Functions/Hello.java"
}

function runExperiment() {
    # Update the Java code with the new array size
    updateArraySize $1

    # Make sure we get fresh data by resetting functions via update
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/Functions/; javac -cp gson-2.10.1.jar Hello.java"
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/Functions/; jar cvf hello.jar Hello.class"
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update helloJava Functions/hello.jar --main Hello"
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/; WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update hello Functions/wordcount.js"

    # Start generating load
    source Experiment.sh http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world JS

    # Retrieve warm/cold status of each activation
    scp JSactivation_ids.txt am_CU@apt069.apt.emulab.net:/users/am_CU/openwhisk-devtools/docker-compose/Scripts/
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/Scripts/; bash ./activation_status_checker.sh ./JSactivation_ids.txt"
    scp am_CU@apt069.apt.emulab.net:/users/am_CU/openwhisk-devtools/docker-compose/Scripts/JSactivation_ids.txt_startStates.txt ./ 

    # Start generating load
    source Experiment.sh http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world Java

    # Retrieve warm/cold status of each activation
    scp Javaactivation_ids.txt am_CU@apt069.apt.emulab.net:/users/am_CU/openwhisk-devtools/docker-compose/Scripts/
    ssh am_CU@apt069.apt.emulab.net "cd /users/am_CU/openwhisk-devtools/docker-compose/Scripts/; bash ./activation_status_checker.sh ./Javaactivation_ids.txt"
    scp am_CU@apt069.apt.emulab.net:/users/am_CU/openwhisk-devtools/docker-compose/Scripts/Javaactivation_ids.txt_startStates.txt ./ 

    # Plot response curves
    python response_time_plotter.py JSOutputTime.txt JSactivation_ids.txt_startStates.txt
    python response_time_plotter.py JavaOutputTime.txt Javaactivation_ids.txt_startStates.txt
    # Clean up after each iteration
    rm *.txt || true
}

# Run the experiments for the three array sizes
for size in 100 10000 1000000; do
    runExperiment $size
    # Move all PNG files to that directory
    mv /users/am_CU/openwhisk-devtools/docker-compose/Graphs/*.png "/users/am_CU/openwhisk-devtools/docker-compose/Graphs/$size/"
done