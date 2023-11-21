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
	javaServerPort    = "9876"
	goServerPort      = "9875"
	loadBalancerPort  = ":8080"
	javaContainerName = "my-java-server"
	goContainerName   = "my-go-server"
	waitTimeout       = 10 * time.Second
)

var javaRoundRobinIndex int = 0
var goRoundRobinIndex int = 0

const numberOfJavaContainers int = 4
const numberOfGoContainers int = 4

func main() {
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
	cmd := exec.Command("docker", "stop", "$(docker", "ps", "-q)")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error stopping all containers:", err)
	} else {
		fmt.Println("All containers stopped successfully")
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL, containerName string

	switch r.URL.Path {
	case "/java":
		containerName = scheduleJavaContainer()
		fmt.Println("Selected container:", containerName)
		targetURL = "http://localhost:" + javaServerPort + "/jsonresponse"
	case "/go":
		containerName = scheduleGoContainer()
		fmt.Println("Selected container:", containerName)
		targetURL = "http://localhost:" + goServerPort + "/GoNative"
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Start the container and wait for it to be ready
	CheckandStartContainer(containerName)
	if !waitForServerReady(targetURL) {
		http.Error(w, "Server is not ready", http.StatusServiceUnavailable)
		return
	}

	forwardRequest(w, r, targetURL)
}

func scheduleJavaContainer() string {
	index := javaRoundRobinIndex % numberOfJavaContainers
	javaRoundRobinIndex++
	return javaContainerName + fmt.Sprintf("-%d", index)
}

func scheduleGoContainer() string {
	index := goRoundRobinIndex % numberOfGoContainers
	goRoundRobinIndex++
	return goContainerName + fmt.Sprintf("-%d", index)
}

func forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	resp, err := http.Get(targetURL + "?" + r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
}

func CheckandStartContainer(containerName string) {
	fmt.Println("Starting container:", containerName)

	// Check if the container is already running
	if isContainerRunning(containerName) {
		fmt.Println("Container already running:", containerName)
		return // Container is already running
	}

	// Start the container
	startNewContainer(containerName)
}

// Check if a container is already running
func isContainerRunning(containerName string) bool {
	// Check if the container is already running
	cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking running container:", err)
		return false
	}
	if string(output) != "" {
		fmt.Println("Container already running:", containerName)
		return true // Container is already running
	}

	// Check if a stopped container with the name exists
	cmd = exec.Command("docker", "ps", "-a", "-q", "-f", "name="+containerName)
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error checking stopped container:", err)
		return false
	}
	if string(output) != "" {
		// Remove the existing container
		fmt.Println("Removing existing container:", containerName)
		cmd = exec.Command("docker", "rm", containerName)
		if err := cmd.Run(); err != nil {
			fmt.Println("Error removing container:", err)
			return false
		}
	}
	return false
}

// Start a new container
func startNewContainer(containerName string) {
	var portMapping, imageName string

	// Define the port mapping and image name based on container prefix
	if strings.HasPrefix(containerName, javaContainerName) {
		portMapping = "9876:9876"
		imageName = "my-java-server"
	} else if strings.HasPrefix(containerName, goContainerName) {
		portMapping = "9875:9875"
		imageName = "my-go-server"
	} else {
		fmt.Println("Unknown container name:", containerName)
		return
	}

	cmd := exec.Command("docker", "run", "-d", "--name", containerName, "-p", portMapping, imageName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error starting container:", containerName, err)
		return
	} else {
		fmt.Println("Container started:", containerName)
		return
	}
}

// func startContainer(containerName string) {
//     fmt.Println("Starting container:", containerName)

//     // Define the port mapping and image name
//     var portMapping, imageName string
//     switch containerName {
//     case javaContainerName:
//         portMapping = "9876:9876"
//         imageName = "my-java-server"
//     case goContainerName:
//         portMapping = "9875:9875"
//         imageName = "my-go-server"
//     }

//     // Start the container
//     cmd = exec.Command("docker", "run", "-d", "--name", containerName, "-p", portMapping, imageName)
//     if err := cmd.Start(); err != nil {
//         fmt.Println("Error starting container:", containerName, err)
//     } else {
//         fmt.Println("Container started:", containerName)
//     }
// }

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

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	fmt.Println("Copying response: ")
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
