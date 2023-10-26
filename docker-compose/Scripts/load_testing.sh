OW_SERVER_NODE="am_CU@node0"
NATIVE_JAVA_API="http://128.110.96.62:9876/jsonresponse"
JAVA_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world"
JAVASCRIPT_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world"
GO_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloGo/world"
OW_DIRECTORY="/users/am_CU/openwhisk-devtools/docker-compose"
GC_FLAGS="-Xmx64m -XX:MaxGCPauseMillis=50 -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/PureJava/gc_log"
NO_GC_FLAGS="-Xmx4g -XX:MaxGCPauseMillis=500 -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/PureJava/no_gc_log"
ITERATIONS=1000
RATE=350

# Start experiments
bash ParallelExperiment.sh $NATIVE_JAVA_API NativeJava $ITERATIONS $RATE
