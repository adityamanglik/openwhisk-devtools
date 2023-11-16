package main

import (
	"io"
	"net/http"
	"os/exec"
	"time"
)

const (
	javaServerPort = "9876"
	goServerPort   = "9875"
	loadBalancerPort = ":9800"
	javaContainerName = "my-java-server"
	goContainerName = "my-go-server"
)

func main() {
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(loadBalancerPort, nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL string

	// Check the endpoint and start the respective container
	switch r.URL.Path {
	case "/java":
		targetURL = "http://localhost:" + javaServerPort
		startContainer(javaContainerName)
	case "/go":
		targetURL = "http://localhost:" + goServerPort
		startContainer(goContainerName)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Forward the request to the container
	resp, err := http.Get(targetURL + r.URL.String())
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the container to the client
	io.Copy(w, resp.Body)

	// Here you can retrieve and analyze GC stats, then make decisions for future requests
}

func startContainer(containerName string) {
	// Check if the container is already running
	cmd := exec.Command("docker", "ps", "-q", "-f", "name=" + containerName)
	output, _ := cmd.Output()
	if string(output) != "" {
		return // Container is already running
	}

	// Start the container
	cmd = exec.Command("docker", "run", "-d", "--rm", "--name", containerName, containerName)
	cmd.Start()
	// Note: Add error handling and possibly a waiting mechanism for the container to be ready
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
