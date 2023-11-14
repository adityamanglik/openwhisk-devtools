package main

import (
    "encoding/json"
    "log"
    "math/rand"
    "net/http"
    "runtime"
    "strconv"
)

const arraySize = 1000000

func main() {
    http.HandleFunc("/GoNative", jsonHandler)
    log.Println("Server listening on http://localhost:9875/GoNative")
    log.Fatal(http.ListenAndServe(":9875", nil))
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
