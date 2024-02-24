# Perform multiple throughput testing
ab -n 10000 -c 1 "http://node0:8180/go?seed=1000&arraysize=10000&requestnumber=56" | grep "Requests per second"
ab -n 10000 -c 1 "http://node0:8180/go?seed=1001&arraysize=10000&requestnumber=57" | grep "Requests per second"
ab -n 10000 -c 1 "http://node0:8180/go?seed=1002&arraysize=10000&requestnumber=58" | grep "Requests per second"
ab -n 10000 -c 1 "http://node0:8180/go?seed=1003&arraysize=10000&requestnumber=56" | grep "Requests per second"
ab -n 10000 -c 1 "http://node0:8180/go?seed=1004&arraysize=10000&requestnumber=57" | grep "Requests per second"
ab -n 10000 -c 1 "http://node0:8180/go?seed=1005&arraysize=10000&requestnumber=58" | grep "Requests per second"
