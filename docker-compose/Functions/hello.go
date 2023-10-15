package main

import (
	"log"
	"math/rand"
	"runtime"
)

type Response struct {
	Sum               int64  `json:"sum"`
	HeapAllocMemory   uint64 `json:"heapAllocMemory"`
	HeapSysMemory     uint64 `json:"heapSysMemory"`
	HeapIdleMemory    uint64 `json:"heapIdleMemory"`
	HeapInuseMemory   uint64 `json:"heapInuseMemory"`
	HeapReleasedMemory uint64 `json:"heapReleasedMemory"`
	HeapObjects       uint64 `json:"heapObjects"`
}

// MARKER_FOR_SIZE_UPDATE
const ARRAY_SIZE = 5000000;

// Main is the function implementing the action
func Main(obj map[string]interface{}) map[string]interface{} {
	seedValue := int64(42) // default seed value

	if seed, exists := obj["seed"].(float64); exists {
		seedValue = int64(seed)
	}

	r := rand.New(rand.NewSource(seedValue))
	arr := make([]int, ARRAY_SIZE)
	var sum int64 = 0

	for i := 0; i < len(arr); i++ {
		arr[i] = r.Intn(ARRAY_SIZE) // populate array with random integers between 0 and (ARRAY_IZE - 1)
		sum += int64(arr[i])
	}

	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	response := make(map[string]interface{})

	response["sum"] = sum
	response["heapAllocMemory"] = m.Alloc
	response["heapSysMemory"] = m.HeapSys
	response["heapIdleMemory"] = m.HeapIdle
	response["heapInuseMemory"] = m.HeapInuse
	response["heapReleasedMemory"] = m.HeapReleased
	response["heapObjects"] = m.HeapObjects

	// log in stdout
	log.Printf("Seed=%d\n", seedValue)

	return response
}
