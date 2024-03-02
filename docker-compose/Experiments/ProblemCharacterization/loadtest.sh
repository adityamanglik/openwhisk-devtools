# Perform multiple throughput testing
size=100
ab -n 10000 -c 1 "http://node0:8801/JS?seed=1000&arraysize=$size&requestnumber=56"
ab -n 10000 -c 1 "http://node0:8801/JS?seed=1001&arraysize=$size&requestnumber=57"
ab -n 10000 -c 1 "http://node0:8801/JS?seed=1002&arraysize=$size&requestnumber=58"