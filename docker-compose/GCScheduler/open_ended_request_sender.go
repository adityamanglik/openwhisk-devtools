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
    "sort"
)

// Constants for API endpoints and file names
const (
    iterations            = 100
    javaAPI               = "http://128.110.96.59:8180/java"
    goAPI                 = "http://128.110.96.59:8180/go"
    javaResponseTimesFile = "java_response_times.txt"
    goResponseTimesFile   = "go_response_times.txt"
	javaServerTimesFile = "java_server_times.txt"
    goServerTimesFile   = "go_server_times.txt"
)

// Response structure for unmarshalling JSON data
type APIResponse struct {
    ExecutionTime int64 `json:"executionTime"`
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
    goResponseTimes, goServerTimes := sendRequests(goAPI, arraysize) 
    writeTimesToFile(goResponseTimesFile, goResponseTimes)
    writeTimesToFile(goServerTimesFile, goServerTimes)
    calculateAndPrintStats(goResponseTimes, "Go Response Times")
    calculateAndPrintStats(goServerTimes, "Go Server Times")

    // ensure server is alive
    checkServerAlive(javaAPI)
    // javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
    javaResponseTimes, javaServerTimes := sendRequests(javaAPI, arraysize) 
    
    // Write time data to files
    writeTimesToFile(javaResponseTimesFile, javaResponseTimes)
    writeTimesToFile(javaServerTimesFile, javaServerTimes)
    calculateAndPrintStats(javaResponseTimes, "Java Response Times")
    calculateAndPrintStats(javaServerTimes, "Java Server Times")
}

func sendRequests(apiURL string, arraysize int) ([]int64, []int64) {
    responseTimesChan := make(chan int64, iterations)
    serverTimesChan := make(chan int64, iterations)

    for i := 0; i < iterations; i++ {
        go func() {
            seed := rand.Intn(10000)
            requestURL := fmt.Sprintf("%s?seed=%d&arraysize=%d", apiURL, seed, arraysize)

            startTime := time.Now()
            resp, err := http.Get(requestURL)
            if err != nil {
                fmt.Println("Error sending request:", err)
                responseTimesChan <- 0
                serverTimesChan <- 0
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                fmt.Println("Non-OK HTTP status code:", resp.StatusCode)
                responseTimesChan <- 0
                serverTimesChan <- 0
                return
            }

            responseBody, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                fmt.Println("Error reading response body:", err)
                responseTimesChan <- 0
                serverTimesChan <- 0
                return
            }

            var apiResp APIResponse
            if err := json.Unmarshal(responseBody, &apiResp); err != nil {
                fmt.Println("Error unmarshalling response:", err)
                responseTimesChan <- 0
                serverTimesChan <- 0
                return
            }

            endTime := time.Now()
            elapsed := endTime.Sub(startTime).Microseconds()

            responseTimesChan <- elapsed
            serverTimesChan <- apiResp.ExecutionTime
        }()
    }

    // Collect the times from the channels
    var responseTimes []int64
    var serverTimes []int64
    for i := 0; i < iterations; i++ {
        responseTimes = append(responseTimes, <-responseTimesChan)
        serverTimes = append(serverTimes, <-serverTimesChan)
    }

    close(responseTimesChan)
    close(serverTimesChan)

    return responseTimes, serverTimes
}


func checkServerAlive(apiURL string) {
    fmt.Println("Checking server for heartbeat.")
    for i := 0; i < iterations; i++ {
        seed := rand.Intn(10000) // Random seed generation
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
            // Break out of the loop if a correct response is received
            break
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