package main

import (
	"container/list"
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type Response struct {
	Sum             int64  `json:"sum"`
	ExecutionTime   int64  `json:"executionTime"`
	RequestNumber   int    `json:"requestNumber"`
	ArraySize       int    `json:"arraysize"`
	HeapAllocMemory uint64 `json:"heapAllocMemory"`
	GOGC            string `json:"GOGC"`
	GOMEMLIMIT      string `json:"GOMEMLIMIT"`
	NextGC          uint64 `json:"NextGC"`
	NumGC           uint32 `json:"NumGC"`
}

// MARKER_FOR_SIZE_UPDATE
// const ARRAY_SIZE = 3200000;

func init() {
	// debug.SetGCPercent(-1) // Disable the garbage collector
}

// Main is the function implementing the action
func Main(obj map[string]interface{}) map[string]interface{} {
	seed := 42               // default seed value
	ARRAY_SIZE := 10000      // default array size value
	REQ_NUM := math.MaxInt32 // default request number

	if val, ok := obj["seed"].(float64); ok {
		seed = int(val)
	}

	if val, ok := obj["arraysize"].(float64); ok {
		ARRAY_SIZE = int(val)
	}

	if val, ok := obj["requestnumber"].(float64); ok {
		REQ_NUM = int(val)
	}

	start := time.Now()

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

	executionTime := time.Since(start).Microseconds()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := Response{
		Sum:             sum,
		ExecutionTime:   executionTime,
		RequestNumber:   REQ_NUM,
		ArraySize:       ARRAY_SIZE,
		HeapAllocMemory: m.HeapAlloc,
		GOGC:            os.Getenv("GOGC"),
		GOMEMLIMIT:      os.Getenv("GOMEMLIMIT"),
		NextGC:          m.NextGC,
		NumGC:           m.NumGC,
	}

	responseMap := make(map[string]interface{})
	responseBytes, _ := json.Marshal(response)
	json.Unmarshal(responseBytes, &responseMap)

	return responseMap
}
