# Create and run actions
wsk action create helloGo hello.go 
wsk action update helloGo --web true
wsk api create /helloGo /world get helloGo --response-type json
export APIHOST=http://node0:3233
export AUTH=23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP
wsk namespace list -v
curl -u $AUTH https://$APIHOST/api/v1/namespaces/_/limits
curl -u 23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP http://node0:3233/api/v1/namespaces/_/limits

# ab -n 10000 http://node0:3234/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloGo/world 