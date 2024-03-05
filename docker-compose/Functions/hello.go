package main

import (
	"container/list"
	"encoding/json"
	"math"
	"math/rand"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
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
	ParsedSeed      string `json:"parsedSeed,omitempty"`
	ParsedArraySize string `json:"parsedArraySize,omitempty"`
	ParsedReqNum    string `json:"parsedReqNum,omitempty"`
}

// MARKER_FOR_SIZE_UPDATE
// const ARRAY_SIZE = 3200000;

func init() {
	debug.SetGCPercent(-1) // Disable the garbage collector

	// Set GOGC, controls the garbage collector target percentage.
	if err := os.Setenv("GOGC", "off"); err != nil {
		panic(err)
	}

	// Set GOMEMLIMIT, an example environment variable. This is not standard in Go.
	// You'll need to implement its usage logic within your application.
	if err := os.Setenv("GOMEMLIMIT", "128M"); err != nil {
		panic(err)
	}
}

// Main is the function implementing the action
func Main(obj map[string]interface{}) map[string]interface{} {
	seed := 42               // default seed value
	ARRAY_SIZE := 10000      // default array size value
	REQ_NUM := math.MaxInt32 // default request number
	response := Response{
		// Existing fields
	}

	if query, ok := obj["__ow_query"].(string); ok && query != "" {
		response.ParsedSeed = "Found query, seed not parsed"
		response.ParsedArraySize = "Found query, arraysize not parsed"
		response.ParsedReqNum = "Found query, requestnumber not parsed"
		values, err := url.ParseQuery(query)
		if err == nil {
			if val, ok := values["seed"]; ok {
				seed, _ = strconv.Atoi(val[0])
			}

			if val, ok := values["arraysize"]; ok {
				ARRAY_SIZE, _ = strconv.Atoi(val[0])
			}

			if val, ok := values["requestnumber"]; ok {
				REQ_NUM, _ = strconv.Atoi(val[0])
			}
		}
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

	response.Sum = sum
	response.ExecutionTime = executionTime
	response.RequestNumber = REQ_NUM
	response.ArraySize = ARRAY_SIZE
	response.HeapAllocMemory = m.HeapAlloc
	response.GOGC = os.Getenv("GOGC")
	response.GOMEMLIMIT = os.Getenv("GOMEMLIMIT")
	response.NextGC = m.NextGC
	response.NumGC = m.NumGC

	responseMap := make(map[string]interface{})
	responseBytes, _ := json.Marshal(response)
	json.Unmarshal(responseBytes, &responseMap)

	return responseMap
}
