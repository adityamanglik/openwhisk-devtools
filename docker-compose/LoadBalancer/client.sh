# This script starts the load balancer server
# The load balancer server receives all requests
# As it receives a request, it checks the endpoint
# Based on the endpoint, it starts the Java or Go container
# Once the container is started, it sends the request to the container
# The container returns the request's result to the load balancer
# The container also returns the current GC stats to the LB
# The load balancer forwards the request's result to the client
# The load balancer looks at the GC data and decides where to send requests
