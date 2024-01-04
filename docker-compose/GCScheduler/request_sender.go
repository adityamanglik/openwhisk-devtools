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

    // ensure server is alive
    checkServerAlive(goAPI)
    // javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
    goResponseTimes, goServerTimes := sendRequests(goAPI, arraysize) 
    writeTimesToFile(goResponseTimesFile, goResponseTimes)
    writeTimesToFile(goServerTimesFile, goServerTimes)

    // ensure server is alive
    checkServerAlive(javaAPI)
    // javaResponseTimes, javaServerTimes := sendRequests(javaAPI)
    javaResponseTimes, javaServerTimes := sendRequests(javaAPI, arraysize) 
    
    // Write time data to files
    writeTimesToFile(javaResponseTimesFile, javaResponseTimes)
    writeTimesToFile(javaServerTimesFile, javaServerTimes)
}

func sendRequests(apiURL string, arraysize int) ([]int64, []int64) {
    var responseTimes []int64
    var serverTimes []int64

    for i := 0; i < iterations; i++ {
        seed := rand.Intn(10000) // Example seed generation
        requestURL1 := fmt.Sprintf("%s?seed=%d", apiURL, seed)
        requestURL := fmt.Sprintf("%s&arraysize=%d", requestURL1, arraysize)

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
