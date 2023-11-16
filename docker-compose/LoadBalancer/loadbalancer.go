package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
    "os/exec"
    "os/signal"
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

	// Stop the Docker containers
	stopContainer(javaContainerName)
	stopContainer(goContainerName)

	fmt.Println("Shutting down load balancer server...")
}

func stopContainer(containerName string) {
	fmt.Println("Stopping container:", containerName)
	cmd := exec.Command("docker", "stop", containerName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error stopping container:", containerName, err)
	} else {
		fmt.Println("Container stopped:", containerName)
	}
}


func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL, containerName string

	switch r.URL.Path {
	case "/java":
		fmt.Println("Java request received", r.URL.Path)
		targetURL = "http://localhost:" + javaServerPort + "/jsonresponse"
		containerName = javaContainerName
	case "/go":
		fmt.Println("Go request received", r.URL.Path)
		targetURL = "http://localhost:" + goServerPort + "/GoNative"
		containerName = goContainerName
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Start the container and wait for it to be ready
	startContainer(containerName)
	if !waitForServerReady(targetURL) {
		http.Error(w, "Server is not ready", http.StatusServiceUnavailable)
		return
	}

	// Forward the request to the container
	resp, err := http.Get(targetURL + "?" + r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyResponse(w, resp)
}

func startContainer(containerName string) {
    fmt.Println("Starting container:", containerName)

    // Check if the container is already running
    cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("Error checking running container:", err)
        return
    }
    if string(output) != "" {
        fmt.Println("Container already running:", containerName)
        return // Container is already running
    }

    // Check if a stopped container with the name exists
    cmd = exec.Command("docker", "ps", "-a", "-q", "-f", "name="+containerName)
    output, err = cmd.Output()
    if err != nil {
        fmt.Println("Error checking stopped container:", err)
        return
    }
    if string(output) != "" {
        // Remove the existing container
        fmt.Println("Removing existing container:", containerName)
        cmd = exec.Command("docker", "rm", containerName)
        if err := cmd.Run(); err != nil {
            fmt.Println("Error removing container:", err)
            return
        }
    }

    // Define the port mapping and image name
    var portMapping, imageName string
    switch containerName {
    case javaContainerName:
        portMapping = "9876:9876"
        imageName = "my-java-server"
    case goContainerName:
        portMapping = "9875:9875"
        imageName = "my-go-server"
    }

    // Start the container
    cmd = exec.Command("docker", "run", "-d", "--name", containerName, "-p", portMapping, imageName)
    if err := cmd.Start(); err != nil {
        fmt.Println("Error starting container:", containerName, err)
    } else {
        fmt.Println("Container started:", containerName)
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
