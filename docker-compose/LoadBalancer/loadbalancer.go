package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// SCHEDULING POLICY DATA STRUCTURES//////////////////////////////////////////////////////////////////////

// func newGoGCStructure() *GoGCStructure {
// 	GoGC := GoGCStructure{}
// 	GoGC.currentHeapIdle = 0
// 	GoGC.currentHeapAlloc = 0
// 	GoGC.HeapAllocThreshold = 0
// 	GoGC.GCThreshold = 0.0
// 	return &GoGC
// }

// 	maxGoHeapSize             = 6692864 // Max size of the heap
// 	GoGCTriggerThreshold      = 0.60    // GC is triggered at 55% utilization
// 	resumeGoRequestsThreshold = 0.90    // Resume normal operations at 90% idle
// }

// Add scheduling policy selection logic
type SchedulingPolicy int

const (
	RoundRobin   SchedulingPolicy = 1
	GCMitigation SchedulingPolicy = 2
)

// Track current scheduling policy
var currentSchedulingPolicy SchedulingPolicy = GCMitigation
var handlingGCForGoContainers bool
var GoGCTriggerThreshold float32
var GoGCIdleHeapThreshold int64

// type GoGCStructure struct {
// 	currentHeapIdle    int64
// 	currentHeapAlloc   int64
// 	HeapAllocThreshold int64
// 	GCThreshold        float32
// }

// Track heap across active go containers
// var GoContainerHeapTracker = make(map[string]*GoGCStructure)
var mutexHandlingGCForGoContainers sync.Mutex

// Fake request array size
var fakeRequestArraySize int

// var mutexGoContainerHeapTracker sync.Mutex

// NETWORK CONNECTION DATA STRUCTURES//////////////////////////////////////////////////////////////////////

// Start values for port numbers
const javaPortStart = 8400
const goPortStart = 9500

var javaRoundRobinIndex int = javaPortStart
var goRoundRobinIndex int = goPortStart

// Global http.Client with Transport settings for high-performance
var client = &http.Client{
	Timeout: 60 * time.Second, // Set the timeout to 5 seconds
	Transport: &http.Transport{
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Data structures to parse the JSON
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
	NextGC        int64 `json:"NextGC"`
	NumGC         int64 `json:"NumGC"`
	Sum           int64 `json:"sum"`
	RequestNumber int64 `json:"requestNumber"`
}

const (
	javaServerPort   = "9876"
	goServerPort     = "9875"
	loadBalancerPort = ":8180"
	javaServerImage  = "java-server-image"
	goServerImage    = "go-server-image"
	waitTimeout      = 10 * time.Second
	serverIP         = "http://node0:"
)

// OPERATION DATA STRUCTURES//////////////////////////////////////////////////////////////////////

// We do not need 64 containers for the paper, only 2 should suffice to show the idea works
const maxNumberOfJavaContainers int = 2
const maxNumberOfGoContainers int = 2

// Global request identifier
var globalRequestCounter int64

// Since there are only two containers, we do not need to worry about assigning both to same CPU
// There is plenty of space among 22 CPUs

// Track running containers
var aliveContainers = make(map[string]string)

// Allocate dedicated CPU for container
var currentCPUIndex int = 10 + rand.Intn(10)

// Log file handler
var (
	logChannel chan string
)

// MAIN   //////////////////////////////////////////////////////////////////////
func init() {
	// Stop all running Docker containers
	stopAllRunningContainers()

	// Initialize the request counter variable
	globalRequestCounter = 0
	// Initialize GC thresholds
	GoGCTriggerThreshold = 0.935
	GoGCIdleHeapThreshold = 100000

	fakeRequestArraySize = 100000

	// If GCMitigation Policy, start and warm the containers
	if currentSchedulingPolicy == GCMitigation {
		container1 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
		container2 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
		startNewContainer(container1)
		startNewContainer(container2)
		mutexHandlingGCForGoContainers.Lock()
		handlingGCForGoContainers = false
		mutexHandlingGCForGoContainers.Unlock()

		// Send 10000 request to warm up containers
		for j := 0; j <= 2000; j++ {
			seed := rand.Intn(10000)
			arraysize := fakeRequestArraySize
			requestURL := serverIP + aliveContainers[container1] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(arraysize)
			// Send fake request
			resp, err := http.Get(requestURL)
			if err != nil {
				fmt.Println("Error sending fake request:", err)
				continue
			} else {
				resp.Body.Close() // Ensure response body is closed
			}

			requestURL = serverIP + aliveContainers[container2] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(arraysize)
			// Send fake request
			// Send fake request
			resp, err = http.Get(requestURL)
			if err != nil {
				fmt.Println("Error sending fake request:", err)
				continue
			} else {
				resp.Body.Close() // Ensure response body is closed
			}
		}
		// initialize GCTracker values
		SendFakeRequest(container1)
		SendFakeRequest(container2)
		time.Sleep(5 * time.Second)
		// fmt.Println("Sent request to initialize GC data structure")
		// fmt.Printf("HeapIdle: %d, HeapAlloc: %d GCThresh %f \n", GoContainerHeapTracker[container1].currentHeapIdle, GoContainerHeapTracker[container1].currentHeapAlloc, GoContainerHeapTracker[container1].GCThreshold)
	} else if currentSchedulingPolicy == RoundRobin {
		container1 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
		container2 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
		startNewContainer(container1)
		startNewContainer(container2)
		handlingGCForGoContainers = false

		// Send 10000 request to warm up containers
		for j := 0; j <= 2000; j++ {
			seed := rand.Intn(10000)
			arraysize := fakeRequestArraySize
			requestURL := serverIP + aliveContainers[container1] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(arraysize)
			// Send fake request
			resp, err := http.Get(requestURL)
			if err != nil {
				fmt.Println("Error sending fake request:", err)
				continue
			} else {
				resp.Body.Close() // Ensure response body is closed
			}

			requestURL = serverIP + aliveContainers[container2] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(arraysize)
			// Send fake request
			// Send fake request
			resp, err = http.Get(requestURL)
			if err != nil {
				fmt.Println("Error sending fake request:", err)
				continue
			} else {
				resp.Body.Close() // Ensure response body is closed
			}
		}
		// initialize GCTracker values
		SendFakeRequest(container1)
		SendFakeRequest(container2)
		time.Sleep(5 * time.Second)
	}
	// Initialize the log channel with a buffer size of 100
	logChannel = make(chan string, 110)

	// Start the logger goroutine
	go loggerRoutine()
}

func main() {
	// Inform go runtime that we are constrained to a single CPU
	// runtime.GOMAXPROCS(1)

	http.HandleFunc("/", handleRequest)
	fmt.Println("Load Balancer is running on port", loadBalancerPort)

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start the listener in a goroutine to provide concurrency needed to listen for the signal below
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
	close(logChannel)
	fmt.Println("Load balancer server shot down.")
}

// Stop all running Docker containers
func stopAllRunningContainers() {
	// fmt.Println("Stopping all running Docker containers...")

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

	fmt.Println("All docker containers stopped successfully.")
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
		// add container to heap tracker
		// mutexGoContainerHeapTracker.Lock()
		// GoContainerHeapTracker[containerName] = newGoGCStructure()
		// mutexGoContainerHeapTracker.Unlock()
	} else {
		fmt.Println("Unknown container name:", containerName)
		// Die fast
		panic(1)
	}

	// Assign a specific CPU to the container and increment the CPU index
	cpuSet := strconv.Itoa(currentCPUIndex)

	currentCPUIndex++
	// Ensure currentCPUIndex doesn't exceed your system's CPU count
	if currentCPUIndex >= 31 { // assuming you have 32 CPUs
		currentCPUIndex = 10 + rand.Intn(10) // reset to 11 or handle as needed
	}

	cmd := exec.Command("docker", "run", "--cpuset-cpus", cpuSet, "--memory=128m", "-d", "--rm", "--name", containerName, "-p", portMapping, imageName)
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
		// Die fast
		fmt.Println("Error running container:", containerName)
		panic(1)
	} else { // container successfully started
		if !waitForServerReady(targetURL) {
			http.Error(w, "Server is not ready", http.StatusServiceUnavailable)
			return
		}
		aliveContainers[containerName] = containerPort
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

// REQUEST HANDLER //////////////////////////////////////////////////////////////////////

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL, containerName, port string
	// Increment counter for every valid request
	globalRequestCounter += 1

	switch r.URL.Path {
	case "/java":
		containerName = scheduleJavaContainer()
		fmt.Printf("GRequest: %d, Selected container: %s\n", globalRequestCounter, containerName)
		port = containerName[len(javaServerImage)+1:] // "+1" to skip the hyphen
		targetURL = serverIP + port + "/jsonresponse"
	case "/go":
		containerName = scheduleGoContainer()
		fmt.Printf("GRequest: %d, Selected container: %s\n", globalRequestCounter, containerName)
		port = containerName[len(goServerImage)+1:] // "+1" to skip the hyphen
		targetURL = serverIP + port + "/GoNative"
	case "/exitCall":
		fmt.Println("Exit call received. Initiating shutdown...")
		stopAllRunningContainers()
		os.Exit(0)
	default:
		http.Error(w, "Requested API Not found", http.StatusNotFound)
		return
	}

	// Extract request number for passing to local functions

	// Extract seed value from the query parameters
	seedValue := r.URL.Query().Get("seed")
	arraysizeValue := r.URL.Query().Get("arraysize")
	requestNumber := r.URL.Query().Get("requestnumber")
	// Append seed value to the targetURL if it's present
	if seedValue != "" {
		targetURL += "?seed=" + seedValue
	}

	if arraysizeValue != "" {
		targetURL += "&arraysize=" + arraysizeValue
	}
	// Append request number to targetURL
	if requestNumber != "" {
		targetURL += "&requestnumber=" + requestNumber
	}

	// print(targetURL)

	// Start the container and wait for it to be ready
	// fmt.Println("Checking and starting container:", containerName)

	//  Check if the container is already running
	if !isContainerRunning(containerName) {
		// startNewContainer(containerName)
		// FOR NOW Assume containers stay alive for execution
		fmt.Print("Container " + containerName + " is not alive, KILLING LoadBalancer\n")
		panic(1)
	}

	forwardRequest(w, r, targetURL, containerName, requestNumber)
}

func forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string, containerName string, requestNumber string) {
	// Send request to container
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy header from incoming request to new request
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

	// Extract and log heap info for each request
	// Off the critical path
	extractAndLogHeapInfo(reader2, containerName, requestNumber)
}

// SCHEDULING POLICY //////////////////////////////////////////////////////////////////////

func scheduleJavaContainer() string {
	switch currentSchedulingPolicy {
	case RoundRobin:
		javaRoundRobinIndex = (javaRoundRobinIndex % maxNumberOfJavaContainers) + javaPortStart
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	case GCMitigation:
		// TODO
		javaRoundRobinIndex = (javaRoundRobinIndex % maxNumberOfJavaContainers) + javaPortStart // Shift starting port number
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	default:
		// Default to Round Robin if the policy is not implemented
		javaRoundRobinIndex = (javaRoundRobinIndex % maxNumberOfJavaContainers) + javaPortStart // Shift starting port number
		javaRoundRobinIndex++
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
	}
}

func scheduleGoContainer() string {
	switch currentSchedulingPolicy {
	case RoundRobin:
		goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)

	case GCMitigation:
		// fmt.Println("In GCMITIGATION Sched policy")
		// fmt.Printf("goRoundRobinIndex: %d\n", goRoundRobinIndex)
		targetContainer := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
		// if we are performing cleanup, send requests to other containers
		mutexHandlingGCForGoContainers.Lock()
		localReadValue := handlingGCForGoContainers
		mutexHandlingGCForGoContainers.Unlock()
		if localReadValue == true {
			fmt.Println("handlingGCForGoContainers is True")
			targetContainer = goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
		}
		return targetContainer

	default:
		// use Round Robin as default
		goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	}
}

// Send a single fake requests to targetContainer to trigger GC
func SendFakeRequest(containerName string) {
	// Send a fake request if heap utilization is above the trigger threshold
	fmt.Printf("Sending fake request to tip over the container %s\n", containerName)
	// Generate fake request
	seed := rand.Intn(10000)
	arraysize := fakeRequestArraySize
	if strings.Contains(containerName, "go") {
		requestURL := serverIP + aliveContainers[containerName] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(arraysize)

		// Send fake request
		resp, err := http.Get(requestURL)
		if err != nil {
			fmt.Println("Error sending fake request:", err)
		}
		defer resp.Body.Close() // Ensure response body is closed

		// Read and unmarshal the response body
		responseBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		reader1 := bytes.NewReader(responseBody)
		// Extract and log heap info for each FAKE request to trigger again if still no GC
		extractAndLogHeapInfo(reader1, containerName, strconv.Itoa(math.MaxInt32))
	}
}

// This function is OFF the critical path for every request, legit or fake
func extractAndLogHeapInfo(responseBody io.Reader, containerName string, requestNumber string) {
	bodyBytes, err := ioutil.ReadAll(responseBody)
	if err != nil {
		fmt.Println("Error reading response body for metrics: ", err)
		return
	}
	var heapInfo string
	var builder strings.Builder
	if strings.Contains(containerName, "java") {
		var javaResp JavaResponse
		if err := json.Unmarshal(bodyBytes, &javaResp); err != nil {
			fmt.Println("Java JSON unmarshalling error:", err)
		} else {
			builder.WriteString("HeapUsedMemory: ")
			builder.WriteString(strconv.FormatInt(javaResp.HeapUsedMemory, 10))
			builder.WriteString(", HeapCommittedMemory: ")
			builder.WriteString(strconv.FormatInt(javaResp.HeapCommittedMemory, 10))
			builder.WriteString(", HeapMaxMemory: ")
			builder.WriteString(strconv.FormatInt(javaResp.HeapMaxMemory, 10))
			heapInfo = builder.String()
			logHeapInfo("java_heap_memory.log", heapInfo)

			// heapInfo = fmt.Sprintf("HeapUsedMemory: %d, HeapCommittedMemory: %d, HeapMaxMemory: %d\n", javaResp.HeapUsedMemory, javaResp.HeapCommittedMemory, javaResp.HeapMaxMemory)
			// logHeapInfo("java_heap_memory.log", heapInfo)
		}
	} else if strings.Contains(containerName, "go") {
		var goResp GoResponse
		if err := json.Unmarshal(bodyBytes, &goResp); err != nil {
			fmt.Println("Go JSON unmarshalling error:", err)
		} else {
			fmt.Printf("Request: %d, Container: %s, HeapAlloc: %d, HeapIdle: %d, NextGC: %d, NumGC: %d RN:%d\n", goResp.RequestNumber, containerName, goResp.HeapAlloc, goResp.HeapIdle, goResp.NextGC, goResp.NumGC, goResp.RequestNumber)
			// Fake requests have invalid request number
			if goResp.RequestNumber != math.MaxInt32 {

				builder.WriteString("Request: ")
				builder.WriteString(strconv.FormatInt(goResp.RequestNumber, 10))
				builder.WriteString(", Container: ")
				builder.WriteString(containerName)
				builder.WriteString(", HeapAlloc: ")
				builder.WriteString(strconv.FormatInt(goResp.HeapAlloc, 10))
				builder.WriteString(", HeapIdle: ")
				builder.WriteString(strconv.FormatInt(goResp.HeapIdle, 10))
				builder.WriteString(", HeapInuse: ")
				builder.WriteString(strconv.FormatInt(goResp.HeapInuse, 10))
				builder.WriteString(", NextGC: ")
				builder.WriteString(strconv.FormatInt(goResp.NextGC, 10))
				builder.WriteString(", NumGC: ")
				builder.WriteString(strconv.FormatInt(goResp.NumGC, 10))
				builder.WriteString("\n")

				heapInfo := builder.String()
				// fmt.Println(heapInfo)
				logHeapInfo("go_heap_memory.log", heapInfo)

				// heapInfo = fmt.Sprintf("Request: %d, Container: %s, HeapAlloc: %d, HeapIdle: %d, HeapInuse: %d, NextGC: %d, NumGC: %d\n", goResp.RequestNumber, containerName, goResp.HeapAlloc, goResp.HeapIdle, goResp.HeapInuse, goResp.NextGC, goResp.NumGC)
				// fmt.Println(heapInfo)
				// logHeapInfo("go_heap_memory.log", heapInfo)
			}
			// track heap stats in struct

			// GoContainerHeapTracker[containerName].currentHeapAlloc = goResp.HeapAlloc
			// GoContainerHeapTracker[containerName].currentHeapIdle = goResp.HeapIdle
			GCThresh := float32(goResp.HeapAlloc) / float32(goResp.NextGC)
			// GoContainerHeapTracker[containerName].GCThreshold = GCThresh

			// print the tracked stats
			// fmt.Printf("Updated tracker from extractAndLogHeapInfo for %s \n", containerName)
			// fmt.Printf("HeapIdle: %d, HeapAlloc: %d GCThresh %f \n", goResp.HeapIdle, goResp.HeapAlloc, float32(goResp.HeapAlloc)/float32(goResp.NextGC))

			// if target container is likely to undergo GC, schedule to alternate and force GC on target
			if goResp.HeapIdle < int64(GoGCIdleHeapThreshold) {
				fmt.Printf("targetContainer: %s\t", containerName)
				fmt.Printf("HeapIdle < %d = %d\n", GoGCIdleHeapThreshold, goResp.HeapIdle)
				// Make sure to signal in process
				mutexHandlingGCForGoContainers.Lock()
				handlingGCForGoContainers = true
				mutexHandlingGCForGoContainers.Unlock()
				go func() {
					SendFakeRequest(containerName)
				}()
				return
			}
			if GCThresh >= GoGCTriggerThreshold {
				fmt.Printf("GCThreshold >= GoGCTriggerThreshold %f\n", GCThresh)
				// Make sure to signal in process
				mutexHandlingGCForGoContainers.Lock()
				handlingGCForGoContainers = true
				mutexHandlingGCForGoContainers.Unlock()
				go func() {
					SendFakeRequest(containerName)
				}()
				return
			}
			// If both conditions are false, set lever to 9500
			mutexHandlingGCForGoContainers.Lock()
			handlingGCForGoContainers = false
			mutexHandlingGCForGoContainers.Unlock()
		}
	}
}

// logHeapInfo sends log information to the logChannel.
func logHeapInfo(filename, info string) {
	var builder strings.Builder
	builder.WriteString("/users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/")
	builder.WriteString(filename)
	builder.WriteString("+ ")
	builder.WriteString(info)

	// Send the constructed log entry to the channel
	logChannel <- builder.String()
}

// writeLogToFile writes the log entry to the specified file.
func writeLogToFile(logEntry string) {
	// Extract filename and info from logEntry
	parts := strings.SplitN(logEntry, "+ ", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid log entry format")
		return
	}
	filename, info := parts[0], parts[1]

	// Open the file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the info to the file
	if _, err := file.WriteString(info + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// loggerRoutine handles writing log messages to files.
func loggerRoutine() {
	for logEntry := range logChannel {
		writeLogToFile(logEntry)
	}
}
