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
	"bytes"
)

var iterations int = 10
var actualIterations int = 1000

// Constants for API endpoints and file names
const (
	javaAPI               = "http://node0:8180/java"
	goAPI                 = "http://node0:8080/run"
	javaResponseTimesFile = "java_response_times.txt"
	goResponseTimesFile   = "go_response_times.txt"
	javaServerTimesFile   = "java_server_times.txt"
	goServerTimesFile     = "go_server_times.txt"
	goHeapFile            = "go_heap_memory.log"
)

// Response structure for unmarshalling JSON data
type APIResponse struct {
	ExecutionTime       int64 `json:"executionTime"`
	UsedHeapSize        int64 `json:"heapUsedMemory"` // Ensure this matches the JSON key exactly
	// GC1CollectionCount  int64 `json:"gc1CollectionCount"`
	// GC1CollectionTime   int64 `json:"gc1CollectionTime"`
	// GC2CollectionCount  int64 `json:"gc2CollectionCount"`
	// GC2CollectionTime   int64 `json:"gc2CollectionTime"`
	// HeapInitMemory      int64 `json:"heapInitMemory"`      // Removed the colon and space
	// HeapCommittedMemory int64 `json:"heapCommittedMemory"` // Removed the colon and space
	// HeapMaxMemory       int64 `json:"heapMaxMemory"`       // Removed the colon and space
}

func main() {
	// Set a default value for arraysize
	defaultArraySize := 10000
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
	fmt.Printf("Arraysize: %d\n", arraysize)
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
		// Generate random seed and construct request payload
		seed := rand.Int63()
		payload := map[string]interface{}{
			"seed":       seed,
			"arraysize":  arraysize,
			"request":    i,
		}

		// Serialize payload to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error serializing payload:", err)
			continue
		}

		// Send POST request with JSON payload
		startTime := time.Now()
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}

		// Check for non-OK status codes
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Non-OK HTTP status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		// Parse the response JSON
		var apiResp APIResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&apiResp); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Record response times and other metrics
		elapsed := time.Since(startTime).Microseconds()
		responseTimes = append(responseTimes, elapsed)
		serverTimes = append(serverTimes, apiResp.ExecutionTime)
		heapSizes = append(heapSizes, apiResp.UsedHeapSize)
	}

	return responseTimes, serverTimes, heapSizes
}

func checkServerAlive(apiURL string) {
    fmt.Println("Checking server for heartbeat.")
    for i := 0; i < iterations/10; i++ {
        seed := rand.Intn(10000) // Random seed generation
        arraysize := 10          // Minimal memory usage for alive check
        payload := map[string]interface{}{
            "seed":      seed,
            "arraysize": arraysize,
        }

        // Serialize payload to JSON
        payloadBytes, err := json.Marshal(payload)
        if err != nil {
            fmt.Println("Error serializing payload:", err)
            continue
        }

        // Send POST request with JSON payload
        resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
        if err != nil {
            fmt.Println("Error sending request:", err)
            time.Sleep(time.Second)
            continue
        }

        // Log HTTP status code
        fmt.Printf("HTTP Status Code: %d\n", resp.StatusCode)

        // Read and log the response body
        responseBody, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close() // Ensure the body is closed after reading
        if err != nil {
            fmt.Println("Error reading response body:", err)
            time.Sleep(time.Second)
            continue
        }

        if len(responseBody) == 0 {
            fmt.Println("Response body is empty.")
        } else {
            fmt.Println("Response:", string(responseBody))
        }

        // Exit if the server responds with OK
        if resp.StatusCode == http.StatusOK {
            fmt.Println("Server is alive and responding.")
            break
        } else {
            fmt.Printf("Server responded with non-OK status: %d\n", resp.StatusCode)
            time.Sleep(time.Second)
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
