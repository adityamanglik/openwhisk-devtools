import requests
import time
import threading
import numpy as np

# Global variables to store the number of requests and total time taken
total_requests = 0
start_time = None
end_time = None

# The URL to test
NATIVE_JAVA_API="http://128.110.96.62:9876/jsonresponse"
JAVA_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world"
JAVASCRIPT_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world"
GO_API="http://128.110.96.62:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloGo/world"

# Function to send a request to the server
def is_java_response_correct(data):
    return all(key in data for key in ["sum", "gc1CollectionCount", "gc1CollectionTime", "gc2CollectionCount", "gc2CollectionTime", "heapInitMemory: ", "heapUsedMemory: ", "heapCommittedMemory: "])

def is_go_response_correct(data):
    return all(key in data for key in ["heapAllocMemory", "heapIdleMemory", "heapInuseMemory", "heapObjects", "heapReleasedMemory", "heapSysMemory", "sum"])

def is_generic_response_correct(data):
    return "payload" in data and "The sum of the array values is" in data["payload"]

lock = threading.Lock()

def send_request(url):
    global total_requests
    try:
        # Generate a random seed value
        seed_value = np.random.randint(0, 1e6)
        if url != NATIVE_JAVA_API:
        # Append the seed parameter to the URL
            url = f"{url}?seed={seed_value}"

        # Send a GET request with the seed parameter
        response = requests.get(url)

        # Check if the response has the correct structure
        data = response.json()
        if is_java_response_correct(data) or is_go_response_correct(data) or is_generic_response_correct(data):
            with lock:
                # Increment the request count only if the response is correct
                total_requests += 1
    except:
        pass


# Function to inject traffic
def inject_traffic(num_threads, url):
    global start_time
    
    threads = []
    start_time = time.time()

    for _ in range(num_threads):
        thread = threading.Thread(target=send_request, args=(url,))
        thread.start()
        threads.append(thread)

    for thread in threads:
        thread.join()

    global end_time
    end_time = time.time()

# List to store throughput values
throughput_values = []

# Number of iterations
iterations = 100

# APIs = [NATIVE_JAVA_API, JAVA_API, JAVASCRIPT_API, GO_API]
APIs = [NATIVE_JAVA_API]
# A dictionary to store throughput results for each API
api_results = {}

for URL in APIs:
    # List to store throughput values for this API
    throughput_values = []
    
    for _ in range(iterations):
        # Randomly vary the number of threads
        num_threads = np.random.randint(100, 2000)
        
        # Reset total requests for each iteration
        total_requests = 0
        
        inject_traffic(num_threads, URL)
        
        # Calculate throughput
        duration = end_time - start_time
        throughput = total_requests / duration
        throughput_values.append(throughput)

    # Calculate the median, mean, and 99th percentile throughput for this API
    median_throughput = np.median(throughput_values)
    mean_throughput = np.mean(throughput_values)
    percentile_99_throughput = np.percentile(throughput_values, 99)
    
    print(f"Results for {URL}:")
    print(f"Median Throughput: {median_throughput} requests/second")
    print(f"Mean Throughput: {mean_throughput} requests/second")
    print(f"99th Percentile Throughput: {percentile_99_throughput} requests/second")

    with open(f"throughput_values_{URL.split('/')[-2]}.txt", "w") as file:
        for value in throughput_values:
            file.write(str(value) + "\n")