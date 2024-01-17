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
)

// Constants for API endpoints and file names
const (
	iterations            = 100000
	javaAPI               = "http://128.110.96.59:8180/java"
	goAPI                 = "http://128.110.96.59:8180/go"
	KILL_SERVER_API       = "http://128.110.96.59:8180/exitCall"
	javaResponseTimesFile = "java_response_times.txt"
	goResponseTimesFile   = "go_response_times.txt"
	javaServerTimesFile   = "java_server_times.txt"
	goServerTimesFile     = "go_server_times.txt"
)

// Response structure for unmarshalling JSON data
type APIResponse struct {
	ExecutionTime int64 `json:"executionTime"`
}

var totalExecutionTime int64

func main() {
	// Set a default value for arraysize
	arraysize := 10000
	gogc := 1

	// Check if a command line argument is provided
	if len(os.Args) > 1 {
		arraysizeStr := os.Args[1]
		GOGCStr := os.Args[2]
		if convertedSize, err := strconv.Atoi(arraysizeStr); err == nil {
			arraysize = convertedSize // Update only if conversion is successful
		} else {
			fmt.Printf("Invalid array size provided, using default value %d\n", arraysize)
		}
		if convertedgogc, err := strconv.Atoi(GOGCStr); err == nil {
			gogc = convertedgogc // Update only if conversion is successful
		} else {
			fmt.Printf("Invalid GOGC provided, using default value %d\n", gogc)
		}
	}
	fmt.Printf("Arraysize: %d\n", arraysize)
	fmt.Printf("GOGC: %d\n", gogc)
	// ensure server is alive
	checkServerAlive(goAPI)
	// javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
	goResponseTimes, goServerTimes := sendRequests(goAPI, arraysize)
	if len(goResponseTimes) == 0 || len(goServerTimes) == 0 {
		return
	}
	writeTimesToFile(goResponseTimesFile, goResponseTimes)
	writeTimesToFile(goServerTimesFile, goServerTimes)
	// calculateAndPrintStats(goResponseTimes, "Go Response Times")
	// calculateAndPrintStats(goServerTimes, "Go Server Times")
	filePath := ""
	if gogc != -1 {
		filePath = fmt.Sprintf("./Data/%d_%d_latencies.csv", arraysize, gogc)
	} else {
		filePath = fmt.Sprintf("./Data/%d_DISABLED_latencies.csv", arraysize)
	}
	err := writeToCSV(filePath, arraysize, goResponseTimes, goServerTimes)
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
	for i := 0; i < iterations; i++ {
		seed := rand.Intn(10000)      // Random seed generation
		arraysize := rand.Intn(10000) // Random seed generation
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
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
			} else {
				fmt.Println("Response Body:", string(responseBody))
			}
			// Break out of the loop if a correct response is received
			break
		}
	}
}

// Function to log time values to a file
func writeTimesToFile(filename string, times []int64) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, time := range times {
		_, err := file.WriteString(strconv.FormatInt(time, 10) + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}

func calculateAndPrintStats(times []int64, label string) {
	if len(times) == 0 {
		fmt.Println("No data to calculate statistics for", label)
		return
	}

	// Sort the slice for percentile calculation
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	// Calculate the average
	var sum int64
	for _, t := range times {
		sum += t
	}
	avg := sum / int64(len(times))

	// Helper function to calculate percentiles
	percentile := func(p float64) int64 {
		if len(times) == 0 {
			return 0
		}
		index := int(float64(len(times)-1) * p)
		return times[index]
	}

	fmt.Printf("Statistics for %s:\n", label)
	fmt.Printf("Average: %d\n", avg)
	fmt.Printf("P50 (Median): %d\n", percentile(0.50))
	fmt.Printf("P99: %d\n", percentile(0.99))
	fmt.Printf("P99.9: %d\n", percentile(0.999))
	fmt.Printf("P99.99: %d\n", percentile(0.9999))
}

func writeToCSV(fileName string, arraySize int, responseTimes, serverTimes []int64) error {
	// Create a new file or truncate existing file
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
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

	// Function to calculate the average
	average := func(times []int64) int64 {
		if len(times) == 0 {
			return 0
		}
		var sum int64
		for _, t := range times {
			sum += t
		}
		return sum / int64(len(times))
	}

	// Calculate statistics
	responseAvg := average(responseTimes)
	responseP50 := percentile(responseTimes, 0.50)
	responseP99 := percentile(responseTimes, 0.99)
	responseP999 := percentile(responseTimes, 0.999)
	responseP9999 := percentile(responseTimes, 0.9999)

	serverAvg := average(serverTimes)
	serverP50 := percentile(serverTimes, 0.50)
	serverP99 := percentile(serverTimes, 0.99)
	serverP999 := percentile(serverTimes, 0.999)
	serverP9999 := percentile(serverTimes, 0.9999)

	// Writing headers
	headers := []string{"ArraySize", "totalExecutionTime", "ClientAvg", "ClientP50", "ClientP99", "ClientP999", "ClientP9999", "ServerAvg", "ServerP50", "ServerP99", "ServerP999", "ServerP9999"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers to csv: %v", err)
	}

	// Write the data to CSV
	record := []string{
		strconv.Itoa(arraySize),
		strconv.FormatInt(totalExecutionTime, 10),
		strconv.FormatInt(responseAvg, 10),
		strconv.FormatInt(responseP50, 10),
		strconv.FormatInt(responseP99, 10),
		strconv.FormatInt(responseP999, 10),
		strconv.FormatInt(responseP9999, 10),
		strconv.FormatInt(serverAvg, 10),
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
