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
)

const arraySize = 0
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

    jsonResponse, executionTime, err := mainLogic(seed)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)

    log.Printf("Request processed in %v\n", executionTime)
}

func mainLogic(seed int) ([]byte, time.Duration, error) {
    start := time.Now()
    
    rand.Seed(int64(seed))

    arr := make([]int, arraySize)
    var sum int64 = 0

    for i := range arr {
        arr[i] = rand.Intn(100000)
        sum += int64(arr[i])
    }

    response := map[string]interface{}{
        "sum": sum,
    }

    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    response["heapAlloc"] = m.HeapAlloc
    response["heapSys"] = m.HeapSys
    response["heapIdle"] = m.HeapIdle
    response["heapInuse"] = m.HeapInuse

    executionTime := time.Since(start)
    response["executionTime"] = executionTime.String()

    jsonResponse, err := json.Marshal(response)
    return jsonResponse, executionTime, err
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
