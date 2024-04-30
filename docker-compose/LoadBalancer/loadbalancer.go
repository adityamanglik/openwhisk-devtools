// ASSUMPTIONS //////////////////////////////////////////////////////////////////////
// 0. Assumption: Docker container CPU allocator assumes 16 core CPU, loadbalancer is running on CPU 1
// 1. Assumption: NEVER taskset the loadbalancer
// 2. There are only two containers for ALL experiments
// 3. Fixed query size as 100K for sanity and consistency
// 4. Container numbering starts after the portStart --> Server Port = 9500, C1 = 9501, C2 = 9502

package main

import (
	"bufio"
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
	"path/filepath"
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

// type GoGCStructure struct {
// 	currentHeapIdle    int64
// 	currentHeapAlloc   int64
// 	HeapAllocThreshold int64
// 	GCThreshold        float32
// }

// // Heap to track Heap usage across multiple containers
// var GoContainerHeapTracker = make(map[string]*GoGCStructure)

// // Mutex to synchronize access to dict declared above
// var mutexGoContainerHeapTracker sync.Mutex

// Add scheduling policy selection logic
type SchedulingPolicy int

const (
	RoundRobin   SchedulingPolicy = 1
	GCMitigation SchedulingPolicy = 2
	SingleServer SchedulingPolicy = 3
)

var currentSchedulingPolicy SchedulingPolicy

// var GoGCTriggerThreshold float32
// var GoGCIdleHeapThreshold int64
var prevHeapAlloc1 int64
var prevNextGC1 int64
var prevHeapAlloc2 int64
var prevNextGC2 int64

var RequestHeapMargin int

// Fake request array size
var fakeRequestArraySize int

// Track heap across active go containers
var mutexIncrementGoPointer sync.Mutex

// var mutexHandlingGCForGoContainers2 sync.Mutex
// var handlingGCForGoContainers1 bool
// var handlingGCForGoContainers2 bool

// NETWORK CONNECTION DATA STRUCTURES//////////////////////////////////////////////////////////////////////

// Start values for port numbers
const javaPortStart = 8400
const goPortStart = 9500

// Global http.Client with Transport settings for high-performance
var client = &http.Client{
	Timeout: 60 * time.Second, // Set the timeout to 5 seconds
	Transport: &http.Transport{
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     90 * time.Second,
	},
}

const (
	javaServerPort      = "8400"
	goServerPort        = "9500"
	loadBalancerPort    = ":8180"
	javaServerImage     = "java-server-image"
	goServerImage       = "go-server-image"
	waitTimeout         = 60 * time.Second // Deadline for container to start and respond
	serverIP            = "http://node0:"
	Available_CPU_Count = 15
)

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
	Sum           int64  `json:"sum"`
	RequestNumber int64  `json:"requestNumber"`
	ExecutionTime int64  `json:"executionTime"`
	ArraySize     int    `json:"arraysize"`
	HeapAlloc     int64  `json:"heapAlloc"`
	HeapIdle      int64  `json:"heapIdle"`
	NextGC        int64  `json:"NextGC"`
	NumGC         int64  `json:"NumGC"`
	GOGC          string `json:"GOGC"`
	GOMEMLIMIT    string `json:"GOMEMLIMIT"`
}

// OPERATION DATA STRUCTURES//////////////////////////////////////////////////////////////////////

// We do not need 64 containers for the paper, only 2 should suffice to show the idea works
const maxNumberOfJavaContainers int = 2
const maxNumberOfGoContainers int = 2

// Allocate dedicated CPU for container
var currentCPUIndex int

// Indices for scheduling containers
var javaRoundRobinIndex int
var goRoundRobinIndex int

// Track running containers
var aliveContainers = make(map[string]string)

// Log file handler
var (
	logChannel chan string
)

// Global request identifier
var globalRequestCounter int64

// MAIN   //////////////////////////////////////////////////////////////////////
func init() {
	// Stop all running Docker containers
	stopAllRunningContainers()

	// Delete logs from previous execution
	deleteExistingLogFiles()

	// Since there are only two containers, we do not need to worry about assigning both to same CPU
	// There is plenty of space among 16 CPUs
	// Allocate dedicated CPU for container
	currentCPUIndex = 2 + rand.Intn(Available_CPU_Count)

	// Indices for scheduling containers
	javaRoundRobinIndex = javaPortStart + 1
	goRoundRobinIndex = goPortStart + 1

	// Initialize the request counter variable
	globalRequestCounter = 0

	// Set default scheduling policy
	currentSchedulingPolicy = RoundRobin

	// Read command line parameters to set scheduling policy
	if len(os.Args) > 1 {
		policy := os.Args[1]
		if policy == "GCMitigation" {
			currentSchedulingPolicy = GCMitigation
		} else if policy == "RoundRobin" {
			currentSchedulingPolicy = RoundRobin
		} else if policy == "SingleServer" {
			currentSchedulingPolicy = SingleServer
		}
	} else {
		fmt.Printf("Usage: loadbalancer <Scheduling Policy = RoundRobin GCMitigation SingleServer>\n")
		fmt.Printf("Invalid or NO scheduling policy provided, using default value: RoundRobin.\n")
	}
	fmt.Printf("Scheduling policy selected: %d\n", currentSchedulingPolicy)

	// Initialize the log channel with a buffer size of 100
	logChannel = make(chan string, 110)

	// Start the logger goroutine
	go loggerRoutine()
}

func main() {
	// Start request handler
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
	} else { // start the containers
		if currentSchedulingPolicy == RoundRobin {
			// Start twp containers
			container1 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
			container2 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
			startNewContainer(container1)
			startNewContainer(container2)
		} else if currentSchedulingPolicy == SingleServer {
			// Start one container
			container1 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
			startNewContainer(container1)
		} else if currentSchedulingPolicy == GCMitigation {
			// Start two containers
			container1 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
			container2 := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
			startNewContainer(container1)
			startNewContainer(container2)
	
			// Initialize GC thresholds based on GOGC value from Dockerfile
			SetGoGCThresholds()
	
			// Warm up containers
			// numWarmUpRequests := 100
	
			// // TODO: Move warm up request to client instead of server
			// // Warm up containers
			// for j := 0; j <= numWarmUpRequests; j++ {
			// 	// Send same request to both containers
			// 	seed := rand.Intn(10000)
			// 	// Container 1
			// 	requestURL := serverIP + aliveContainers[container1] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(fakeRequestArraySize1)
			// 	// Send fake request
			// 	resp, err := http.Get(requestURL)
			// 	if err != nil {
			// 		fmt.Println("Error sending fake request:", err)
			// 		continue
			// 	} else {
			// 		resp.Body.Close() // Ensure response body is closed
			// 	}
	
			// 	// Container 2
			// 	requestURL = serverIP + aliveContainers[container2] + "/GoNative?seed=" + strconv.Itoa(seed) + "&arraysize=" + strconv.Itoa(fakeRequestArraySize2)
			// 	// Send fake request
			// 	resp, err = http.Get(requestURL)
			// 	if err != nil {
			// 		fmt.Println("Error sending fake request:", err)
			// 		continue
			// 	} else {
			// 		resp.Body.Close() // Ensure response body is closed
			// 	}
			// }
	
			// initialize GC Structure values
			// SendFakeRequest(container1)
			// SendFakeRequest(container2)
			fmt.Println("Sent requests to initialize GC data structure")
		}
		return true
	}
	return false

	// EDGE CASE
	// IMPLEMENT if containers start becoming unresponsive over 20-30 minutes runtime
	// Solution: create a coroutine that routinely polls containers under idle time (no request received at loadbalancer) to check if they are dead or alive. For any successful request returned to client, the trigger is reset to avoid polluting experiments

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
	// Build container image before starting
	if strings.HasPrefix(containerName, "go") {
		cmd := exec.Command("docker", "build", "-t", "go-server-image", "/users/am_CU/openwhisk-devtools/docker-compose/Native/Go/")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error building container image:", containerName, err)
			panic(1)
		}
	} else if strings.HasPrefix(containerName, "java") {
		cmd := exec.Command("docker", "build", "-t", "java-server-image", "/users/am_CU/openwhisk-devtools/docker-compose/Native/Java/")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error building container image:", containerName, err)
			panic(1)
		}
	} // Fresh contianer images are now built and ready to launch

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
		// Die fast
		panic(1)
	}

	// Assign a specific CPU to the container and increment the CPU index
	cpuSet := strconv.Itoa(currentCPUIndex)
	// Increment CPU allocation pointer
	currentCPUIndex++
	// Ensure currentCPUIndex doesn't exceed your system's CPU count
	if currentCPUIndex > Available_CPU_Count {
		// Since there are only two containers, we do not need to worry about assigning both to same CPU
		// There is plenty of space among 16 CPUs
		currentCPUIndex = 2 + rand.Intn(Available_CPU_Count)
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

	// Extract query parameters
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

	// fmt.Printf("targetURL: %s, Selected container: %s\n", targetURL, containerName)

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

	extractAndLogHeapInfo(reader2, containerName, requestNumber)
}

// SCHEDULING POLICY //////////////////////////////////////////////////////////////////////

func scheduleJavaContainer() string {
	switch currentSchedulingPolicy {
	case SingleServer:
		// server ports are always one ahead of port start
		return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex+1)
	case RoundRobin:
		javaRoundRobinIndex = (javaRoundRobinIndex % maxNumberOfJavaContainers) + javaPortStart
		// server ports are always one ahead of port start
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
	if currentSchedulingPolicy == SingleServer {
		// server ports are always one ahead of port start
		fmt.Println("In SingleServer Sched policy")
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	} else if currentSchedulingPolicy == RoundRobin {
		fmt.Println("In RoundRobin Sched policy")
		goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
		// server ports are always one ahead of port start
		goRoundRobinIndex++
		return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
	} else if currentSchedulingPolicy == GCMitigation {
		fmt.Println("In GCMITIGATION Sched policy")
		fmt.Printf("goRoundRobinIndex: %d\n", goRoundRobinIndex)
		// Default container for requests
		targetContainer := goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
		// Try and acquire lock, if can acquire, send request?
		// mutexHandlingGCForGoContainers1.Lock()
		// localReadValue := handlingGCForGoContainers1
		// mutexHandlingGCForGoContainers1.Unlock()
		// // if we are performing cleanup, send requests to other containers
		// if localReadValue == true {
		// 	fmt.Println("handlingGCForGoContainers is True")
		// 	targetContainer = goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex+1)
		// }
		return targetContainer
	} else {
		// use Round Robin as default
		goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
		// server ports are always one ahead of port start
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
	// Log the memory statistics
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
		// Java logging ends here ////////////////////////////////////////////////////////////////////////
	} else if strings.Contains(containerName, "go") {
		var goResp GoResponse
		if err := json.Unmarshal(bodyBytes, &goResp); err != nil {
			fmt.Println("Go JSON unmarshalling error:", err)
		} else {
			fmt.Printf("Request: %d, Container: %s, HeapAlloc: %d, NextGC: %d, NumGC: %d, HeapIdle: %d\n", goResp.RequestNumber, containerName, goResp.HeapAlloc, goResp.NextGC, goResp.NumGC, goResp.HeapIdle)
			// Fake requests have invalid request number
			if goResp.RequestNumber != math.MaxInt32 {
				builder.WriteString("Request: ")
				builder.WriteString(strconv.FormatInt(goResp.RequestNumber, 10))
				builder.WriteString(", Container: ")
				builder.WriteString(containerName)
				builder.WriteString(", HeapAlloc: ")
				builder.WriteString(strconv.FormatInt(goResp.HeapAlloc, 10))
				builder.WriteString(", NextGC: ")
				builder.WriteString(strconv.FormatInt(goResp.NextGC, 10))
				builder.WriteString(", NumGC: ")
				builder.WriteString(strconv.FormatInt(goResp.NumGC, 10))
				builder.WriteString(", HeapIdle: ")
				builder.WriteString(strconv.FormatInt(goResp.HeapIdle, 10))
				builder.WriteString("\n")

				heapInfo := builder.String()
				// fmt.Println(heapInfo)
				logHeapInfo("go_heap_memory.log", heapInfo)

			}
			// GCMitigation Policy ////////////////////////////////////////////////////////////////////////
			if currentSchedulingPolicy == GCMitigation {
				if strings.Contains(containerName, "1") {
					// GC is triggered for HeapAlloc breaching NextGC
					// marginAvailable := int64(RequestHeapMargin) * (goResp.HeapAlloc - prevHeapAlloc1)
					marginAvailable := int64(RequestHeapMargin) * prevHeapAlloc1
					fmt.Printf("BEFORE prevHeapAlloc1: %d, Margin: %d, nextGC: %d\n", prevHeapAlloc1, marginAvailable, prevNextGC1)
					// Current container likely under heap memory pressure
					if (goResp.HeapAlloc + marginAvailable) > prevNextGC1 {
						fmt.Printf("targetContainer: %s\t", containerName)
						fmt.Printf("(currHeapAlloc + margin) > prevNextGC) --> %d + %d > %d\n", goResp.HeapAlloc, marginAvailable, prevNextGC1)
						if goResp.RequestNumber != math.MaxInt32 {
							// Increment pointer to start sending requests to second container
							mutexIncrementGoPointer.Lock()
							goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
							// server ports are always one ahead of port start
							goRoundRobinIndex++
							mutexIncrementGoPointer.Unlock()
						}
						// Send fake request in parallel while updating pointer
						go func() {
							SendFakeRequest(containerName)
						}()
						// Update heap tracker data structure for second container to prevent deadlock
						// Send fake request to second container before returning
						// SendFakeRequest(goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex))
						return // because we need heap statistics from new container, NOT dead one
					}
					// Update latest HeapAlloc and NextGC otherwise
					// prevHeapAlloc1 = goResp.HeapAlloc
					prevNextGC1 = goResp.NextGC
					fakeRequestArraySize = goResp.ArraySize

					fmt.Printf("AFTER prevHeapAlloc1: %d, nextGC1: %d, arraysize: %d\n", prevHeapAlloc1, prevNextGC1, fakeRequestArraySize)
				} else { // second container
					// GC is triggered for HeapAlloc breaching NextGC
					// marginAvailable := int64(RequestHeapMargin) * (goResp.HeapAlloc - prevHeapAlloc2)
					marginAvailable := int64(RequestHeapMargin) * prevHeapAlloc2
					fmt.Printf("BEFORE prevHeapAlloc2: %d, Margin: %d, nextGC: %d\n", prevHeapAlloc2, marginAvailable, prevNextGC2)
					// Current container likely under heap memory pressure
					if (goResp.HeapAlloc + marginAvailable) > prevNextGC2 {
						fmt.Printf("targetContainer: %s\t", containerName)
						fmt.Printf("(currHeapAlloc + margin) > prevNextGC) --> %d + %d > %d\n", goResp.HeapAlloc, marginAvailable, prevNextGC2)
						// Increment pointer if NOT Fake request
						if goResp.RequestNumber != math.MaxInt32 {
							mutexIncrementGoPointer.Lock()
							goRoundRobinIndex = (goRoundRobinIndex % maxNumberOfGoContainers) + goPortStart
							// server ports are always one ahead of port start
							goRoundRobinIndex++
							mutexIncrementGoPointer.Unlock()
						}
						// Send fake request in parallel while updating pointer
						go func() {
							SendFakeRequest(containerName)
						}()
						// Update heap tracker data structure for second container to prevent deadlock
						// Send fake request to second container before returning
						// SendFakeRequest(goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex))
						return // because we need heap statistics from new container, NOT dead one
					}
					// Update latest HeapAlloc and NextGC otherwise
					// prevHeapAlloc2 = goResp.HeapAlloc
					prevNextGC2 = goResp.NextGC
					fakeRequestArraySize = goResp.ArraySize

					fmt.Printf("AFTER prevHeapAlloc2: %d, nextGC2: %d, arraysize: %d\n", prevHeapAlloc2, prevNextGC2, fakeRequestArraySize)

				}
			}
		}

	} // Go logging ends here ////////////////////////////////////////////////////////////////////////
}

func SetGoGCThresholds() {
	// Read dockerfile for GOGC value
	dockerfilePath := "/users/am_CU/openwhisk-devtools/docker-compose/Native/Go/Dockerfile"
	file, err := os.Open(dockerfilePath)
	if err != nil {
		fmt.Printf("Error opening Dockerfile: %s\n", err)
		return
	}
	defer file.Close()
	// Parse Dockerfile to determine GOGC value
	detectedGOGC := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ENV") {
			envVars := strings.Fields(line)
			for i := 1; i < len(envVars); i++ { // Start from 1 to skip "ENV"
				parts := strings.Split(envVars[i], "=")
				if len(parts) == 2 && parts[0] == "GOGC" {
					detectedGOGC1, convErr := strconv.Atoi(parts[1])
					if convErr != nil {
						fmt.Printf("Error converting GOGC value to int: %s\n", convErr)
						fmt.Printf("Assuming default GOGC value: %d\n", detectedGOGC)
						break
					}
					detectedGOGC = detectedGOGC1
					break
				}
			}
		}
	}
	fmt.Printf("GOGC value found: %d\n", detectedGOGC)
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading Dockerfile: %s\n", err)
	}

	// Set initial values assuming GOGC
	if detectedGOGC == 1000 {
		prevHeapAlloc1 = 100000
		prevNextGC1 = 41943040
		prevHeapAlloc2 = 100000
		prevNextGC2 = 41943040
		RequestHeapMargin = 3
		fakeRequestArraySize = 10000
	} else if detectedGOGC == 100 {
		prevHeapAlloc1 = 100000
		prevNextGC1 = 4194304
		prevHeapAlloc2 = 100000
		prevNextGC2 = 4194304
		RequestHeapMargin = 5
		fakeRequestArraySize = 100
	} else if detectedGOGC == 1 {
		prevHeapAlloc1 = 100000
		prevNextGC1 = 1160000
		prevHeapAlloc2 = 100000
		prevNextGC2 = 1160000
		RequestHeapMargin = 5
		fakeRequestArraySize = 100
	} else { // if GOGC not found, change scheduling policy
		fmt.Println("Error detecting GOGC. Switching policy to RoundRobin.")
		currentSchedulingPolicy = RoundRobin
	}
	// Set different margins for different containers
	// RequestHeapMargin = 3

	// Initialize fake request size
	// fakeRequestArraySize = 10000
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

func deleteExistingLogFiles() {
	logFiles, err := filepath.Glob("/users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/*.log")
	if err != nil {
		fmt.Println("Error finding log files:", err)
		return
	}
	// fmt.Println("Found files:", logFiles)
	for _, logFile := range logFiles {
		// path := filepath.Join("/users/am_CU/openwhisk-devtools/docker-compose/LoadBalancer/", logFile)
		if err := os.Remove(logFile); err != nil {
			if !os.IsNotExist(err) { // Ignore error if file doesn't exist
				fmt.Printf("Failed to delete log file %s: %v\n", logFile, err)
			}
		} else {
			fmt.Printf("Deleted existing log file: %s\n", logFile)
		}
	}
}
