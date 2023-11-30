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
    "strconv"
    "syscall"
    "time"
    "runtime"
)

// MARKER_FOR_SIZE_UPDATE
const ARRAY_SIZE = 1000000

const serverPort = ":9875"

func main() {
    ln, err := net.Listen("tcp", serverPort)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }

    server := &http.Server{Addr: serverPort, Handler: nil}
    http.HandleFunc("/GoNative", jsonHandler)
    log.Println("Server listening on http://localhost" + serverPort)

    go func() {
        if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    gracefulShutdown(server)
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

    jsonResponse, err := mainLogic(seed)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)

    // log.Printf("Request processed in %v\n", executionTime)
}

func mainLogic(seed int) ([]byte, error) {
    start := time.Now().UnixNano()
    
    rand.Seed(int64(seed))

    arr := make([]int, ARRAY_SIZE)
    var sum int64 = 0

    for i := range arr {
        arr[i] = rand.Intn(100000)
    }

    for i := range arr {
        sum += int64(arr[i])
    }

    executionTime := time.Now().UnixNano() - start

    response := map[string]interface{}{
        "sum": sum,
        "executionTime": executionTime, // Include raw execution time in nanoseconds
    }
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    response["heapAlloc"] = m.HeapAlloc
    response["heapSys"] = m.HeapSys
    response["heapIdle"] = m.HeapIdle
    response["heapInuse"] = m.HeapInuse
    jsonResponse, err := json.Marshal(response)
    return jsonResponse, err
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
}
