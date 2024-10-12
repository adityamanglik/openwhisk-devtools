package main

// Import the required packages
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var iterations int = 500
var actualIterations int = 1000

// Constants for API endpoint and file names
const (
	goAPI               = "http://node0:9501/GoNative"
	goResponseTimesFile = "go_response_times.txt"
	goServerTimesFile   = "go_server_times.txt"
	goHeapFile          = "go_heap_memory.log"
)

// Response structure for unmarshalling JSON data
type APIResponse struct {
	ExecutionTime int64 `json:"executionTime"`
	HeapAlloc     int64 `json:"heapAlloc"`
}

func main() {
	// Set a default value for arraysize
	defaultArraySize := 100
	arraysize := defaultArraySize

	// Check if a command line argument is provided
	if len(os.Args) > 1 {
		arraysizeStr := os.Args[1]
		if convertedSize, err := strconv.Atoi(arraysizeStr); err == nil {
			arraysize = convertedSize // Update only if conversion is successful
		} else {
			fmt.Printf("Invalid array size provided, using default value %d\n", defaultArraySize)
		}
	}
	fmt.Printf("\nArraysize: %d\n", arraysize)

	// Ensure server is alive
	checkServerAlive(goAPI)

	// Warm up
	goResponseTimes, goServerTimes, _ := sendRequests(goAPI, arraysize)
	iterations = actualIterations
	fmt.Printf("Warm up done, starting measurement run\n")

	// Actual measurements
	goResponseTimes, goServerTimes, _ = sendRequests(goAPI, arraysize)

	// Perform latency analysis
	err := latencyAnalysis(arraysize, goResponseTimes, goServerTimes)
	if err != nil {
		fmt.Println("Error during latency analysis:", err)
	}
}

func sendRequests(apiURL string, arraysize int) ([]int64, []int64, []int64) {
	var responseTimes []int64
	var serverTimes []int64
	var heapSizes []int64

	for i := 0; i < iterations; i++ {
		seed := rand.Intn(10000) // Example seed generation
		requestURL1 := fmt.Sprintf("%s?seed=%d", apiURL, seed)
		requestURL2 := fmt.Sprintf("%s&arraysize=%d", requestURL1, arraysize)
		requestURL := fmt.Sprintf("%s&requestnumber=%d", requestURL2, i)

		startTime := time.Now()
		resp, err := http.Get(requestURL)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Non-OK HTTP status code:", resp.StatusCode)
		}

		// Read and unmarshal the response body
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		var apiResp APIResponse
		if err := json.Unmarshal(responseBody, &apiResp); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			fmt.Println("Response body:", string(responseBody))
			continue
		}

		elapsed := time.Since(startTime)

		responseTimes = append(responseTimes, elapsed.Microseconds())
		serverTimes = append(serverTimes, apiResp.ExecutionTime)
		heapSizes = append(heapSizes, apiResp.HeapAlloc)
	}

	return responseTimes, serverTimes, heapSizes
}

func checkServerAlive(apiURL string) {
	fmt.Println("Checking server for heartbeat.")
	for i := 0; i < iterations/10; i++ {
		seed := rand.Intn(10000) // Random seed generation
		arraysize := 10          // Do not pollute memory for aliveCheck
		requestURL := fmt.Sprintf("%s?seed=%d&arraysize=%d", apiURL, seed, arraysize)
		resp, err := http.Get(requestURL)
		if err != nil {
			fmt.Println("Error sending request:", err)
			time.Sleep(time.Second)
			continue
		}
		defer resp.Body.Close()

		// Check if the HTTP status code is 200 (OK)
		if resp.StatusCode == http.StatusOK {
			fmt.Println("OK Response received from server.")
			// Read and print the response body
			responseBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("Response: ", string(responseBody))
			// Break out of the loop if a correct response is received
			break
		}
	}
}

func latencyAnalysis(arraySize int, responseTimes, serverTimes []int64) error {
	// Helper function to calculate percentiles
	percentile := func(times []int64, p float64) int64 {
		if len(times) == 0 {
			return 0
		}
		sortedTimes := make([]int64, len(times))
		copy(sortedTimes, times)
		sort.Slice(sortedTimes, func(i, j int) bool { return sortedTimes[i] < sortedTimes[j] })
		index := int(float64(len(sortedTimes)-1) * p)
		return sortedTimes[index]
	}

	// Calculate response time statistics
	responseP50 := percentile(responseTimes, 0.50)
	responseP90 := percentile(responseTimes, 0.90)
	responseP95 := percentile(responseTimes, 0.95)
	responseP99 := percentile(responseTimes, 0.99)
	responseP999 := percentile(responseTimes, 0.999)
	responseP9999 := percentile(responseTimes, 0.9999)

	// Calculate server time statistics
	serverP50 := percentile(serverTimes, 0.50)
	serverP90 := percentile(serverTimes, 0.90)
	serverP95 := percentile(serverTimes, 0.95)
	serverP99 := percentile(serverTimes, 0.99)
	serverP999 := percentile(serverTimes, 0.999)
	serverP9999 := percentile(serverTimes, 0.9999)

	// Calculate total server time (in microseconds)
	var totalServerTime int64 = 0
	for _, t := range serverTimes {
		totalServerTime += t
	}

	// Calculate throughput based on total server time
	// Convert totalServerTime to seconds
	totalServerTimeSeconds := float64(totalServerTime) / 1e6
	throughput := float64(len(serverTimes)) / totalServerTimeSeconds

	// Print latency statistics
	fmt.Printf("\nLatency Statistics for Array Size %d:\n", arraySize)
	fmt.Printf("Response Times (microseconds):\n")
	fmt.Printf("P50: %d, P90: %d, P95: %d, P99: %d, P99.9: %d, P99.99: %d\n",
		responseP50, responseP90, responseP95, responseP99, responseP999, responseP9999)
	fmt.Printf("Server Execution Times (microseconds):\n")
	fmt.Printf("P50: %d, P90: %d, P95: %d, P99: %d, P99.9: %d, P99.99: %d\n",
		serverP50, serverP90, serverP95, serverP99, serverP999, serverP9999)
	fmt.Printf("Throughput based on server time: %.2f requests/second\n", throughput)

	return nil
}
