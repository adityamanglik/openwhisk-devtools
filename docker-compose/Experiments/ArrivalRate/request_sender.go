package main

// Import the required packages
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var iterations int = 500
var actualIterations int = 1000

// Constants for API endpoint and file names
const (
	goAPI               = "http://node0:3234/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/GoLL"
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
	goResponseTimes, goServerTimes, heapSizes := sendRequests(goAPI, arraysize)
	iterations = actualIterations
	fmt.Printf("Warm up done, starting measurement run\n")

	// Actual measurements
	goResponseTimes, goServerTimes, heapSizes = sendRequests(goAPI, arraysize)

	// Write time data to files
	writeTimesToFile(goResponseTimesFile, goResponseTimes)
	writeTimesToFile(goServerTimesFile, goServerTimes)
	writeTimesToFile(goHeapFile, heapSizes)
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

// Function to log time values to a file
func writeTimesToFile(filename string, times []int64) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for index, time := range times {
		_, err := file.WriteString(strconv.Itoa(index) + ", " + strconv.FormatInt(time, 10) + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}
