package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
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
)

// Start values for port numbers
var javaRoundRobinIndex int = 8400
var goRoundRobinIndex int = 9500

const numberOfJavaContainers int = 4
const numberOfGoContainers int = 4

// Track running containers
var runningContainers = make(map[string]string)

// Global http.Client with Transport settings for high-performance
var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
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

	switch r.URL.Path {
	case "/java":
		containerName = scheduleJavaContainer()
		fmt.Println("Selected container: ", containerName)
		port = containerName[len(javaServerImage)+1:] // "+1" to skip the hyphen
		targetURL = "http://128.110.96.176:" + port + "/jsonresponse"
	case "/go":
		containerName = scheduleGoContainer()
		fmt.Println("Selected container: ", containerName)
		port = containerName[len(goServerImage)+1:] // "+1" to skip the hyphen
		targetURL = "http://128.110.96.176:" + port + "/GoNative"
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Start the container and wait for it to be ready
	fmt.Println("Checking and starting container:", containerName)

	//  Check if the container is already running
	if !isContainerRunning(containerName) {
		startNewContainer(containerName)
	}

	forwardRequest(w, r, targetURL)
}

func scheduleJavaContainer() string {
	javaRoundRobinIndex = (javaRoundRobinIndex % numberOfJavaContainers) + 8400 // Shift starting port number
	javaRoundRobinIndex++
	return javaServerImage + fmt.Sprintf("-%d", javaRoundRobinIndex)
}

func scheduleGoContainer() string {
	goRoundRobinIndex = (goRoundRobinIndex % numberOfGoContainers) + 9500
	goRoundRobinIndex++
	return goServerImage + fmt.Sprintf("-%d", goRoundRobinIndex)
}

func forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Copying response: ")
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// Check if a container is already running
func isContainerRunning(containerName string) bool {
	// Check if the container is already running
	if _, exists := runningContainers[containerName]; exists {
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
		targetURL = "http://128.110.96.176:" + containerPort + "/jsonresponse"
	} else if strings.HasPrefix(containerName, "go") {
		containerPort = containerName[len(goServerImage)+1:]
		portMapping = containerPort + ":" + goServerPort
		imageName = goServerImage
		targetURL = "http://128.110.96.176:" + containerPort + "/GoNative"
	} else {
		fmt.Println("Unknown container name:", containerName)
		return
	}

	cmd := exec.Command("docker", "run", "-d", "--rm", "--name", containerName, "-p", portMapping, imageName)
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
		runningContainers[containerName] = containerPort
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
