# Create API for JS action
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update "/guest/hello" --web true
# Extract URL of action and place in variable web_action
web_action=$(WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i api create /hello /world get hello --response-type json | tail -n 1)
# Hit URL using curl
echo $web_action

# Create API for Java action
WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i action update "/guest/helloJava" --web true
# Extract URL of action and place in variable web_action
web_action=$(WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i api create /helloJava /world get helloJava --response-type json | tail -n 1)
# Hit URL using curl
echo $web_action