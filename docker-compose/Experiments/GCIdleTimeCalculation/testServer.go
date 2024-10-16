package main

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

const serverPort = ":9500"

func init() {
	// debug.SetGCPercent(-1) // Disable the garbage collector
	// os.Setenv("GOGC", "500")
	// runtime.GOMAXPROCS(2)
}

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

	var m runtime.MemStats
	var prevMemStats runtime.MemStats
	gogcValue := os.Getenv("GOGC")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	iterationNumber := 0 // Initialize iteration counter
	for range ticker.C {
		iterationNumber++ // Increment iteration number on each tick
		runtime.ReadMemStats(&m)
		if m.HeapAlloc != prevMemStats.HeapAlloc || m.HeapIdle != prevMemStats.HeapIdle ||
			m.NextGC != prevMemStats.NextGC || m.NumGC != prevMemStats.NumGC {
			fmt.Printf("Iteration: %d------------------------------\n", iterationNumber)
			fmt.Printf("heapAlloc: %d\n", m.HeapAlloc)
			fmt.Printf("heapIdle: %d\n", m.HeapIdle)
			fmt.Printf("NextGC: %d\n", m.NextGC)
			fmt.Printf("NumGC: %d\n", m.NumGC)
			fmt.Printf("GOGC: %s\n", gogcValue)

			// Update previous values for next comparison
			prevMemStats = m
		}
	}
	gracefulShutdown(server)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	seed := 42               // default seed value
	ARRAY_SIZE := 10000      // default array size value
	REQ_NUM := math.MaxInt32 // default request number

	seedStr := params.Get("seed")
	if seedStr != "" {
		var err error
		seed, err = strconv.Atoi(seedStr)
		if err != nil {
			http.Error(w, "Invalid seed value", http.StatusBadRequest)
			return
		}
	}

	arrayStr := params.Get("arraysize")
	if arrayStr != "" {
		var err error
		ARRAY_SIZE, err = strconv.Atoi(arrayStr)
		if err != nil {
			http.Error(w, "Invalid array size value", http.StatusBadRequest)
			return
		}
	}

	reqNumStr := params.Get("requestnumber")
	if reqNumStr != "" {
		var err error
		REQ_NUM, err = strconv.Atoi(reqNumStr)
		if err != nil {
			http.Error(w, "Invalid request number value", http.StatusBadRequest)
			return
		}
	}

	jsonResponse, err := mainLogic(seed, ARRAY_SIZE, REQ_NUM)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	// log.Printf("Request processed in %v\n", executionTime)
}

func mainLogic(seed int, ARRAY_SIZE int, REQ_NUM int) ([]byte, error) {
	start := time.Now().UnixMicro()

	rand.Seed(int64(seed))

	lst := list.New()

	for i := 0; i < ARRAY_SIZE; i++ {
		// Inserting integers directly, assuming payload simulation isn't the focus
		lst.PushFront(rand.Intn(seed)) // Use integers for direct summation
		// Stress GC with nested list
		if i%5 == 0 {
			nestedList := list.New()
			for j := 0; j < rand.Intn(5); j++ {
				nestedList.PushBack(rand.Intn(seed))
			}
			lst.PushBack(nestedList)
		}
		// Immediate removal after insertion to stress GC
		if i%5 == 0 {
			e := lst.PushFront(rand.Intn(seed))
			lst.Remove(e)
		}

	}

	// Sum values and return result
	var sum int64 = 0
	for e := lst.Front(); e != nil; e = e.Next() {
		if val, ok := e.Value.(int); ok {
			sum += int64(val)
		}
	}

	executionTime := time.Now().UnixMicro() - start

	response := map[string]interface{}{
		"sum":           sum,
		"executionTime": executionTime, // Include raw execution time in microseconds
		"requestNumber": REQ_NUM,
	}

	gogcValue := os.Getenv("GOGC")
	// gomemlimitValue := os.Getenv("GOMEMLIMIT")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response["heapAlloc"] = m.HeapAlloc
	response["NextGC"] = m.NextGC
	// response["heapSys"] = m.HeapSys
	response["heapIdle"] = m.HeapIdle
	// response["heapInuse"] = m.HeapInuse

	response["NumGC"] = m.NumGC
	response["GOGC"] = gogcValue
	// response["GOMEMLIMIT"] = gomemlimitValue
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
