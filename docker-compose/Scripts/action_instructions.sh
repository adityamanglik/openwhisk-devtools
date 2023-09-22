# Create function
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action create hello hello.js
# Invoke function
# -i flag solves x509 certificate problem
res=$(WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action invoke hello --blocking --result)
# Display result
echo "invocation result: $res"
# Create API for action
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update "/guest/hello" --web true
# Extract URL of action and place in variable web_action
web_action=$(WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i api create /hello /world get hello --response-type json | tail -n 1)
# Hit URL using curl
curl -sS "$web_action"
# List all prior activations and cold/warm status along with duration
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i activation list
# Fetch one instance's complete status (cold responses will have additional "key": "initTime")
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i activation get f4c37c93d16b475b837c93d16bc75bb9

# Compile Java (class name - Hello)
javac -cp gson-2.10.1.jar Hello.java
jar cvf hello.jar Hello.class
# Create Java action --> --main Hello is used to specify the main class file name
# An eligible main class is one that implements a static main method as described above. 
# If the class is not in the default package, use the Java fully-qualified class name, e.g., --main com.example.MyMain
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action create helloJava hello.jar --main Hello
# Invoke action with parameter "Workd"
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action invoke --result helloJava --param name World
# Invoke action from client with parameter "Workd"
curl -sS http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world?name=Workd