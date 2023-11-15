package main

import (
    "context"
    "encoding/json"
    "log"
    "math/rand"
    "net"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "strconv"
    "syscall"
    "time"
)

import "sync"
import "io/ioutil"

var (
    executionTimes []time.Duration
    timesMutex     sync.Mutex
)

const arraySize = 0
const serverPort = ":9875"

func main() {
    // Check if the port is already in use
    ln, err := net.Listen("tcp", serverPort)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }

    // Create an HTTP server
    server := &http.Server{Addr: serverPort, Handler: nil}

    // Handle routes
    http.HandleFunc("/GoNative", jsonHandler)
    log.Println("Server listening on http://localhost" + serverPort + "/GoNative")

    // Start server in a goroutine
    go func() {
        if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    // Graceful shutdown
    gracefulShutdown(server)
}

func saveExecutionTimesToFile(filename string) {
    var data string
    for _, t := range executionTimes {
        data += t.String() + "\n"
    }

    if err := ioutil.WriteFile(filename, []byte(data), 0644); err != nil {
        log.Fatalf("Failed to write execution times to file: %v", err)
    }
}

func gracefulShutdown(server *http.Server) {
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    log.Println("Shutting down server...")

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    saveExecutionTimesToFile("execution_times.txt")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    seed := 42 // default seed value

    seedStr := params.Get("seed")
    if seedStr != "" {
        var err error
        seed, err = strconv.Atoi(seedStr)
        if err != nil {
            http.Error(w, "Invalid seed value", http.StatusBadRequest)
            return
        }
    }

    start := time.Now()
    jsonResponse, err := mainLogic(seed)
    elapsed := time.Since(start)

    timesMutex.Lock()
    executionTimes = append(executionTimes, elapsed)
    timesMutex.Unlock()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func mainLogic(seed int) ([]byte, error) {
    rand.Seed(int64(seed))

    arr := make([]int, arraySize)
    var sum int64 = 0

    for i := range arr {
        arr[i] = rand.Intn(100000) // random integers between 0 and 99999
        sum += int64(arr[i])
    }

    response := map[string]int64{
        "sum": sum,
    }

    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    response["heapAlloc"] = int64(m.HeapAlloc)
    response["heapSys"] = int64(m.HeapSys)
    response["heapIdle"] = int64(m.HeapIdle)
    response["heapInuse"] = int64(m.HeapInuse)

    return json.Marshal(response)
}
