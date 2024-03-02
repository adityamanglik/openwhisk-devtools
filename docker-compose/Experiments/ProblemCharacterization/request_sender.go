package main

// Import the required packages
import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat"
)

var iterations int = 1000
var actualIterations int = 5000

// Constants for API endpoints and file names
const (
	javaAPI               = "http://node0:8180/java"
	goAPI                 = "http://node0:8801/JS"
	javaResponseTimesFile = "java_response_times.txt"
	goResponseTimesFile   = "go_response_times.txt"
	javaServerTimesFile   = "java_server_times.txt"
	goServerTimesFile     = "go_server_times.txt"
)

// Response structure for unmarshalling JSON data
type APIResponse struct {
	ExecutionTime int64 `json:"executionTime"`
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
	fmt.Printf("Arraysize: %d\n", arraysize)
	// ensure server is alive
	checkServerAlive(goAPI)
	// javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
	// Warm up
	goResponseTimes, goServerTimes := sendRequests(goAPI, arraysize)
	iterations = actualIterations
	// Actual measurements
	goResponseTimes, goServerTimes = sendRequests(goAPI, arraysize)

	writeTimesToFile(goResponseTimesFile, goResponseTimes)
	writeTimesToFile(goServerTimesFile, goServerTimes)
	// calculateAndPrintStats(goResponseTimes, "Go Response Times")
	// calculateAndPrintStats(goServerTimes, "Go Server Times")
	filePath := fmt.Sprintf("./Graphs/Go/%d/latencies.csv", arraysize)
	err := latencyAnalysis2(filePath, arraysize, goResponseTimes, goServerTimes)
	if err != nil {
		fmt.Println("Error writing to CSV:", err)
	}

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

func sendRequests(apiURL string, arraysize int) ([]int64, []int64) {
	var responseTimes []int64
	var serverTimes []int64

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
	}

	return responseTimes, serverTimes
}

func checkServerAlive(apiURL string) {
	fmt.Println("Checking server for heartbeat.")
	for i := 0; i < iterations/10; i++ {
		seed := rand.Intn(10000)      // Random seed generation
		arraysize := 10 // Do not pollute memory for aliveCheck
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
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func latencyAnalysis(fileName string, arraySize int, responseTimes, serverTimes []int64) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Helper function to calculate percentiles
	percentile := func(times []int64, p float64) int64 {
		if len(times) == 0 {
			return 0
		}
		sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
		index := int(float64(len(times)-1) * p)
		return times[index]
	}

	// Calculate statistics
	responseP50 := percentile(responseTimes, 0.50)
	responseP99 := percentile(responseTimes, 0.99)
	responseP999 := percentile(responseTimes, 0.999)
	responseP9999 := percentile(responseTimes, 0.9999)

	serverP50 := percentile(serverTimes, 0.50)
	serverP99 := percentile(serverTimes, 0.99)
	serverP999 := percentile(serverTimes, 0.999)
	serverP9999 := percentile(serverTimes, 0.9999)

	// Writing headers
	headers := []string{"ArraySize", "ResponseP50", "ResponseP99", "ResponseP999", "ResponseP9999", "ServerP50", "ServerP99", "ServerP999", "ServerP9999"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers to csv: %v", err)
	}

	// Write the data to CSV
	record := []string{
		strconv.Itoa(arraySize),
		strconv.FormatInt(responseP50, 10),
		strconv.FormatInt(responseP99, 10),
		strconv.FormatInt(responseP999, 10),
		strconv.FormatInt(responseP9999, 10),
		strconv.FormatInt(serverP50, 10),
		strconv.FormatInt(serverP99, 10),
		strconv.FormatInt(serverP999, 10),
		strconv.FormatInt(serverP9999, 10),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("error writing record to csv: %v", err)
	}

	return nil
}

func latencyAnalysis2(fileName string, arraySize int, responseTimes, serverTimes []int64) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Function to calculate percentiles
	percentile := func(times []int64, p float64) float64 {
		sortedTimes := make([]float64, len(times))
		for i, v := range times {
			sortedTimes[i] = float64(v)
		}
		sort.Float64s(sortedTimes)
		return stat.Quantile(p, stat.Empirical, sortedTimes, nil)
	}

	// Calculate statistics
	responseP50 := percentile(responseTimes, 0.50)
	responseP99 := percentile(responseTimes, 0.99)
	responseP999 := percentile(responseTimes, 0.999)
	responseP9999 := percentile(responseTimes, 0.9999)

	serverP50 := percentile(serverTimes, 0.50)
	serverP99 := percentile(serverTimes, 0.99)
	serverP999 := percentile(serverTimes, 0.999)
	serverP9999 := percentile(serverTimes, 0.9999)

	// Writing headers
	headers := []string{"ArraySize", "ResponseP50", "ResponseP99", "ResponseP999", "ResponseP9999", "ServerP50", "ServerP99", "ServerP999", "ServerP9999"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers to csv: %v", err)
	}

	// Write the data to CSV
	record := []string{
		fmt.Sprintf("%d", arraySize),
		fmt.Sprintf("%f", responseP50),
		fmt.Sprintf("%f", responseP99),
		fmt.Sprintf("%f", responseP999),
		fmt.Sprintf("%f", responseP9999),
		fmt.Sprintf("%f", serverP50),
		fmt.Sprintf("%f", serverP99),
		fmt.Sprintf("%f", serverP999),
		fmt.Sprintf("%f", serverP9999),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("error writing record to csv: %v", err)
	}

	return nil
}
