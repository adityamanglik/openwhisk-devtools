# Download wsk
wget https://github.com/apache/openwhisk-cli/releases/download/latest/OpenWhisk_CLI-latest-linux-amd64.tgz
tar -zxvf OpenWhisk_CLI-latest-linux-amd64.tgz
sudo mv wsk /usr/local/bin/

# Build openwhisk from source and enable API
git clone https://github.com/apache/openwhisk.git
cd openwhisk
./gradlew core:standalone:bootRun --args='--api-gw'

# Create and run actions
wsk action create test wordcount.js 
wsk action update test --web true
wsk api create /hello /world get test --response-type json
