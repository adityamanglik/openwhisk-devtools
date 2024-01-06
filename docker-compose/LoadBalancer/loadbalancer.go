package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	javaServerPort   = "9876"
	goServerPort     = "9875"
	loadBalancerPort = ":8180"
	javaServerImage  = "java-server-image"
	goServerImage    = "go-server-image"
	waitTimeout      = 10 * time.Second
	serverIP         = "http://128.110.96.59:"

	maxGoHeapSize             = 6692864 // Max size of the heap
	GoGCTriggerThreshold      = 0.60    // GC is triggered at 55% utilization
	resumeGoRequestsThreshold = 0.90    // Resume normal operations at 90% idle
)

// Add scheduling policy selection logic
type SchedulingPolicy int

const (
	RoundRobin    SchedulingPolicy = 1
	HeapSizeBased SchedulingPolicy = 2
)

// Start values for port numbers
var javaRoundRobinIndex int = 8400
var goRoundRobinIndex int = 9500

const numberOfJavaContainers int = 1
const numberOfGoContainers int = 1

// Allocate dedicated CPU for container
var currentCPUIndex int = 10 + rand.Intn(10)

// Track running containers
var aliveContainers = make(map[string]string)
var containerHeapUsage = make(map[string]int64)

// Track current scheduling policy
var currentSchedulingPolicy SchedulingPolicy = RoundRobin

// Global http.Client with Transport settings for high-performance
var client = &http.Client{
	Timeout: 5 * time.Second, // Set the timeout to 5 seconds
	Transport: &http.Transport{
		MaxIdleConns:        99999,
		MaxIdleConnsPerHost: 99999,
		IdleConnTimeout:     90 * time.Second,
	},
}

type JavaResponse struct {
	Sum                 int64 `json:"sum"`
	ExecutionTime       int64 `json:"executionTime"`
	Gc1CollectionCount  int   `json:"gc1CollectionCount"`
	Gc1CollectionTime   int   `json:"gc1CollectionTime"`
	Gc2CollectionCount  int   `json:"gc2CollectionCount"`
	Gc2CollectionTime   int   `json:"gc2CollectionTime"`
	HeapInitMemory      int64 `json:"heapInitMemory: "`
	HeapUsedMemory      int64 `json:"heapUsedMemory: "`
	HeapCommittedMemory int64 `json:"heapCommittedMemory: "`
	HeapMaxMemory       int64 `json:"heapMaxMemory: "`
}

type GoResponse struct {
	ExecutionTime int64 `json:"executionTime"`
	HeapAlloc     int64 `json:"heapAlloc"`
	HeapIdle      int64 `json:"heapIdle"`
	HeapInuse     int64 `json:"heapInuse"`
	HeapSys       int64 `json:"heapSys"`
	Sum           int64 `json:"sum"`
}

func main() {

	// Stop all running Docker containers
	stopAllRunningContainers()

	http.HandleFunc("/", handleRequest)
	fmt.Println("Load Balancer is running on port", loadBalancerPort)

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := http.ListenAndServe(loadBalancerPort, nil); err != nil {
			fmt.Println("Error starting server:", err)
			stopChan <- os.Interrupt
		}
	}()

	// Block until a signal is received.
	<-stopChan

	// Stop all running Docker containers
	stopAllRunningContainers()

	fmt.Println("Shutting down load balancer server...")
}

// Stop all running Docker containers
func stopAllRunningContainers() {
	fmt.Println("Stopping all running Docker containers...")

	// Get the list of all container IDs
	cmd := exec.Command("docker", "ps", "-aq")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting container IDs:", err)
		return
	}
	containerIDs := strings.Fields(string(output))

	// Remove each container
	for _, id := range containerIDs {
		cmd = exec.Command("docker", "rm", "-vf", id)
		if err := cmd.Run(); err != nil {
			fmt.Println("Error removing container:", id, err)
		}
	}

	fmt.Println("All containers stopped successfully")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL, containerName, port string

	// Extract seed value from the query parameters
	seedValue := r.URL.Query().Get("seed")
	arraysizeValue := r.URL.Query().Get("arraysize")

	switch r.URL.Path {
	case "/java":
		containerName = scheduleJavaContainer()
		fmt.Println("Selected container: ", containerName)
		port = containerName[len(javaServerImage)+1:] // "+1" to skip the hyphen
		targetURL = serverIP + port + "/jsonresponse"
	case "/go":
		containerName = scheduleGoContainer()
		fmt.Println("Selected container: ", containerName)
		port = containerName[len(goServerImage)+1:] // "+1" to skip the hyphen
		targetURL = serverIP + port + "/GoNative"
	case "/exitCall":
		fmt.Println("Exit call received. Initiating shutdown...")
		go func() {
			stopAllRunningContainers()
			os.Exit(0)
		}()
		return
	default:
		http.Error(w, "Requested API Not found", http.StatusNotFound)
		return
	}

	// Append seed value to the targetURL if it's present
	if seedValue != "" {
		targetURL += "?seed=" + seedValue
	}

	if arraysizeValue != "" {
		targetURL += "&arraysize=" + arraysizeValue
	}

	// fmt.Println("TargetURL: ", targetURL)

	// Start the container and wait for it to be ready
	fmt.Println("Checking and starting container:", containerName)

	//  Check if the container is already running
	if !isContainerRunning(containerName) {
		startNewContainer(containerName)
	}

	forwardRequest(w, r, targetURL, containerName)
}

func scheduleJavaContainer() string {
	switch currentSchedulingPolicy {
	case RoundRobin:
		javaRoundRobinIndex = (javaRoundRobinIndex % numberOfJavaContainers) + 8400 // Shift starting port number
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	case HeapSizeBased:
		// TODO
		javaRoundRobinIndex = (javaRoundRobinIndex % numberOfJavaContainers) + 8400 // Shift starting port number
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	default:
		// Default to Round Robin if the policy is not implemented
		javaRoundRobinIndex = (javaRoundRobinIndex % numberOfJavaContainers) + 8400 // Shift starting port number
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	}
}

func scheduleGoContainer() string {
	switch currentSchedulingPolicy {
	case RoundRobin:
		goRoundRobinIndex = (goRoundRobinIndex % numberOfGoContainers) + 9500
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	case HeapSizeBased:
		// Choose container based on current heap utilization
		for containerName, heapIdle := range containerHeapUsage {
			if strings.HasPrefix(containerName, goServerImage) {
				heapUtilization := float64(maxGoHeapSize-heapIdle) / float64(maxGoHeapSize)
				fmt.Println("Container: %s, Heap Utilization: %f", containerName, heapUtilization)
				if heapUtilization < GoGCTriggerThreshold {
					return containerName
				} else {
					// Take container offline and send fake requests to trigger GC
					go handleGCForGoContainers(containerName)
					continue
				}
			}
		}
		// If all containers are above the threshold, use Round Robin as fallback
		goRoundRobinIndex = (goRoundRobinIndex % numberOfGoContainers) + 9500
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	default:
		goRoundRobinIndex = (goRoundRobinIndex % numberOfGoContainers) + 9500
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	}
}

func forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string, containerName string) {
	// Send request to container
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := client.Do(req) // Use the global client
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response header and status code to the client
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Read the response body into a buffer
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	// Create two readers from the buffer: one for forwarding, one for logging
	reader1 := bytes.NewReader(bodyBytes)
	reader2 := bytes.NewReader(bodyBytes)

	// Forward the response to the client
	_, err = io.Copy(w, reader1)
	if err != nil {
		fmt.Println("Error forwarding response body: ", err)
		return
	}

	// Extract and log heap info
	extractAndLogHeapInfo(reader2, containerName)
}

func extractAndLogHeapInfo(responseBody io.Reader, containerName string) {
	bodyBytes, err := ioutil.ReadAll(responseBody)

	if err != nil {
		fmt.Println("Error reading response body for metrics: ", err)
		return
	}

	var heapInfo string
	if strings.Contains(containerName, "java") {

		var javaResp JavaResponse
		if err := json.Unmarshal(bodyBytes, &javaResp); err != nil {
			fmt.Println("JSON unmarshalling error:", err)
		} else {
			heapInfo = fmt.Sprintf("HeapUsedMemory: %d, HeapCommittedMemory: %d, HeapMaxMemory: %d\n", javaResp.HeapUsedMemory, javaResp.HeapCommittedMemory, javaResp.HeapMaxMemory)
			// fmt.Println(heapInfo)
			logHeapInfo("java_heap_memory.log", heapInfo)
			containerHeapUsage[containerName] = javaResp.HeapUsedMemory
		}
	} else if strings.Contains(containerName, "go") {
		var goResp GoResponse
		if err := json.Unmarshal(bodyBytes, &goResp); err != nil {
			fmt.Println("JSON unmarshalling error:", err)
		} else {
			heapInfo = fmt.Sprintf("HeapAlloc: %d, HeapIdle: %d, HeapInuse: %d\n", goResp.HeapAlloc, goResp.HeapIdle, goResp.HeapInuse)
			// fmt.Println(heapInfo)
			containerHeapUsage[containerName] = goResp.HeapIdle
			logHeapInfo("go_heap_memory.log", heapInfo)
		}
	}
}

// New function to handle fake requests and GC triggering
func handleGCForGoContainers(containerName string) {
	requestCounter := 0
	for {
		// Fetch the current heap idle value
		heapIdle := containerHeapUsage[containerName]
		heapUtilization := float64(maxGoHeapSize-heapIdle) / float64(maxGoHeapSize)

		// Check if the heap utilization is within the target range
		if heapUtilization < resumeGoRequestsThreshold {
			break // Exit the loop if the condition is met
		}
		fmt.Println("Sending fake requests to tip over the server")
		// Send a fake request if heap utilization is above the trigger threshold
		if heapUtilization >= GoGCTriggerThreshold {
			seed := rand.Intn(10000)
			requestURL := serverIP + containerName[len(goServerImage)+1:] + "/GoNative?seed=" + strconv.Itoa(seed)

			// Process the response to get the latest heap idle value
			resp, err := http.Get(requestURL)
			if err != nil {
				fmt.Println("Error sending fake request:", err)
				continue
			}

			// Read and unmarshal the response body
			responseBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close() // Ensure response body is closed
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}

			var goResp GoResponse
			if err := json.Unmarshal(responseBody, &goResp); err != nil {
				fmt.Println("Error unmarshalling response:", err)
				continue
			}

			// Update the heap idle value
			containerHeapUsage[containerName] = goResp.HeapIdle

			requestCounter++
			if requestCounter > 10000 {
				break // prevent infinite loop
			}
		}

		// time.Sleep(1 * time.Second) // Throttle the loop
	}
	fmt.Println("Go container is clean and ready for use again")
}

func logHeapInfo(filename, info string) {
	fullPath := "/users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/" + filename
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(info); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// Check if a container is already running
func isContainerRunning(containerName string) bool {
	// Check if the container is already running
	if _, exists := aliveContainers[containerName]; exists {
		fmt.Println("Container already running: ", containerName)
		return true // Container is already running
	}
	return false

	// EDGE CASE
	// IMPLEMENT if containers start becoming unresponsive over 20-30 minutes runtime

	// // Check if the container is already running
	// cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
	// output, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println("Error checking running container:", err)
	// 	return false
	// }
	// if string(output) != "" {
	// 	fmt.Println("Container already running:", containerName)
	// 	return true // Container is already running
	// }

	// return false
}

// Start a new container
func startNewContainer(containerName string) {
	var portMapping, imageName, containerPort, targetURL string
	var w http.ResponseWriter

	// Define the port mapping and image name based on container prefix
	if strings.HasPrefix(containerName, "java") {
		containerPort = containerName[len(javaServerImage)+1:]
		portMapping = containerPort + ":" + javaServerPort
		imageName = javaServerImage
		targetURL = serverIP + containerPort + "/jsonresponse"
	} else if strings.HasPrefix(containerName, "go") {
		containerPort = containerName[len(goServerImage)+1:]
		portMapping = containerPort + ":" + goServerPort
		imageName = goServerImage
		targetURL = serverIP + containerPort + "/GoNative"
	} else {
		fmt.Println("Unknown container name:", containerName)
		return
	}

	// Assign a specific CPU to the container and increment the CPU index
	cpuSet := strconv.Itoa(currentCPUIndex)
	currentCPUIndex++

	// Ensure currentCPUIndex doesn't exceed your system's CPU count
	if currentCPUIndex >= 31 { // assuming you have 32 CPUs
		currentCPUIndex = 10 + rand.Intn(10) // reset to 11 or handle as needed
	}

	cmd := exec.Command("docker", "run", "--cpuset-cpus", cpuSet, "-d", "--rm", "--name", containerName, "-p", portMapping, imageName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting container:", containerName, err)
		// Check if a stopped container with the name exists
		cmd = exec.Command("docker", "ps", "-a", "-q", "-f", "name="+containerName)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error checking stopped container:", err)
		}
		if string(output) != "" {
			// Remove the existing container
			fmt.Println("Removing existing container:", containerName)
			cmd = exec.Command("docker", "rm", containerName)
			if err := cmd.Run(); err != nil {
				fmt.Println("Error removing container:", err)
			}
		}
		return
	} else { // container successfully started
		if !waitForServerReady(targetURL) {
			http.Error(w, "Server is not ready", http.StatusServiceUnavailable)
			return
		}
		aliveContainers[containerName] = containerPort
		containerHeapUsage[containerName] = 0
		fmt.Println("Container started:", containerName)
		return
	}
}

func waitForServerReady(url string) bool {
	deadline := time.Now().Add(waitTimeout)
	for time.Now().Before(deadline) {
		fmt.Println("Waiting for container: ", url)
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}
