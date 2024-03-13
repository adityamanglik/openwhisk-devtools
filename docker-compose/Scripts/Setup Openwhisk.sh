# Download wsk
wget https://github.com/apache/openwhisk-cli/releases/download/latest/OpenWhisk_CLI-latest-linux-amd64.tgz
tar -zxvf OpenWhisk_CLI-latest-linux-amd64.tgz
sudo mv wsk /usr/local/bin/

# Build openwhisk from source and enable API
git clone https://github.com/apache/openwhisk.git
cd openwhisk
./gradlew core:standalone:bootRun --args='--api-gw'

# Create and run actions
wsk action create helloGo hello.go 
wsk action update helloGo --web true
wsk api create /helloGo /world get helloGo --response-type json
ab -n 10000 http://node0:3234/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloGo/world 