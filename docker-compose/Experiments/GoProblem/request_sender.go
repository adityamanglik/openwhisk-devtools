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

// Constants for API endpoints and file names
const (
	javaAPI               = "http://node0:8180/java"
	goAPI                 = "http://node0:9501/GoNative"
	javaResponseTimesFile = "java_response_times.txt"
	goResponseTimesFile   = "go_response_times.txt"
	javaServerTimesFile   = "java_server_times.txt"
	goServerTimesFile     = "go_server_times.txt"
	goHeapFile            = "go_heap_memory.log"
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
	// ensure server is alive
	checkServerAlive(goAPI)
	// javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
	// Warm up
	goResponseTimes, goServerTimes, heapSizes := sendRequests(goAPI, arraysize)
	iterations = actualIterations
	fmt.Printf("Warm up done, starting plotting run\n")
	// Actual measurements
	goResponseTimes, goServerTimes, heapSizes = sendRequests(goAPI, arraysize)
	// _ = plotTimes(goResponseTimes, heapSizes, fmt.Sprintf("Server Times for Arraysize %d", arraysize))
	// _ = plotTimes(goResponseTimes, fmt.Sprintf("Server Times for Arraysize %d", arraysize))
	// fmt.Printf("Problem plots done, starting SLA run\n")
	// iterations = 100000
	// arraysize = 10000
	// SLA measurements
	// goResponseTimes, goServerTimes, heapSizes = sendRequests(goAPI, arraysize)
	// _ = plotSLA(goResponseTimes)
	writeTimesToFile(goResponseTimesFile, goResponseTimes)
	writeTimesToFile(goServerTimesFile, goServerTimes)
	writeTimesToFile(goHeapFile, heapSizes)

	// calculateAndPrintStats(goResponseTimes, "Go Response Times")
	// calculateAndPrintStats(goServerTimes, "Go Server Times")
	// filePath := fmt.Sprintf("./Graphs/Go/%d/latencies.csv", arraysize)
	// err = latencyAnalysis2(filePath, arraysize, goResponseTimes, goServerTimes)
	// if err != nil {
	// fmt.Println("Error writing to CSV:", err)
	// }

	// // ensure server is alive
	// checkServerAlive(javaAPI)
	// // javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
	// javaResponseTimes, javaServerTimes := sendRequests(javaAPI, arraysize)

	// // Write time data to files
	// writeTimesToFile(javaResponseTimesFile, javaResponseTimes)
	// writeTimesToFile(javaServerTimesFile, javaServerTimes)
	// // calculateAndPrintStats(javaResponseTimes, "Java Response Times")
	// // calculateAndPrintStats(javaServerTimes, "Java Server Times")
	// filePath = fmt.Sprintf("../Graphs/GCScheduler/Java/%d/latencies.csv", arraysize)
	// err = writeToCSV(filePath, arraysize, javaResponseTimes, javaServerTimes)
	// if err != nil {
	//     fmt.Println("Error writing to CSV:", err)
	// }
}

func sendRequests(apiURL string, arraysize int) ([]int64, []int64, []int64) {
	var responseTimes []int64
	var serverTimes []int64
	var heapSizes []int64

	for i := 0; i < iterations; i++ {
		// fmt.Printf("Sent request: %d\n", i)
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

		endTime := time.Now()
		elapsed := endTime.Sub(startTime)

		responseTimes = append(responseTimes, elapsed.Microseconds())
		serverTimes = append(serverTimes, apiResp.ExecutionTime)
		// fmt.Println("Time:", apiResp.ExecutionTime)
		// Collect usedHeapSize along with other metrics
		// fmt.Println("UsedHeapSize:", apiResp.UsedHeapSize)
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
		// Check if the HTTP status code is 200 (OK)
		if resp.StatusCode == http.StatusOK {
			fmt.Println("OK Response received from server.")
			// Read and unmarshal the response body
			responseBody, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			fmt.Println("Response: ", string(responseBody))
			// Break out of the loop if a correct response is received
			break
		} else {
			resp.Body.Close()
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
