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
    javaAPI               = "http://128.110.96.76:8180/java"
    goAPI                 = "http://128.110.96.76:8180/go"
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
    // Example usage
    sendRequests(javaAPI, javaResponseTimesFile, javaServerTimesFile) // Replace 100 with the actual size
    sendRequests(goAPI, goResponseTimesFile, goServerTimesFile)    // Replace 100 with the actual size
}

func sendRequests(apiURL, responseTimeFile string, serverTimeFile string, ) {
    // Logic to update server code and rebuild Docker images

    for i := 0; i < iterations; i++ {
        seed := rand.Intn(10000) // Example seed generation
        requestURL := fmt.Sprintf("%s?seed=%d", apiURL, seed)

        startTime := time.Now()
        resp, err := http.Get(requestURL)
        if err != nil {
            fmt.Println("Error sending request:", err)
            continue
        }
        defer resp.Body.Close()

        // Read and unmarshal the response body
        responseBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response body:", err)
            continue
        }

        var apiResp APIResponse
        if err := json.Unmarshal(responseBody, &apiResp); err != nil {
            fmt.Println("Error unmarshalling response:", err)
            continue
        }

        endTime := time.Now()
        elapsed := endTime.Sub(startTime)

        // Convert elapsed time to milliseconds and log to file
        logTime(responseTimeFile, elapsed.Nanoseconds())

        // Log the extracted executionTime (in milliseconds)
        logTime(serverTimeFile, apiResp.ExecutionTime) // Assuming executionTime is in nanoseconds
    }

    // Logic to move log files
}

// Function to log time values to a file
func logTime(filename string, time int64) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(strconv.FormatInt(time, 10) + "\n")
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }
}
