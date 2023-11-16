package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
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
	if err := http.ListenAndServe(loadBalancerPort, nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var targetURL, containerName string

	switch r.URL.Path {
	case "/java":
		targetURL = "http://localhost:" + javaServerPort + "/jsonresponse"
		containerName = javaContainerName
	case "/go":
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
	cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
	output, _ := cmd.Output()
	if string(output) != "" {
		return // Container is already running
	}

	cmd = exec.Command("docker", "run", "-d", "--rm", "--name", containerName, containerName)
	cmd.Start()
}

func waitForServerReady(url string) bool {
	deadline := time.Now().Add(waitTimeout)
	for time.Now().Before(deadline) {
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
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
