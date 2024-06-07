package main

import (
	"container/list"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
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

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func main() {
	ln, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	server := &http.Server{Addr: serverPort, Handler: nil}
	http.HandleFunc("/GoNative", jsonHandler)
	http.HandleFunc("/ImageProcess", ImageProcessor)
	log.Println("Server listening on http://localhost" + serverPort)

	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func ImageProcessor(w http.ResponseWriter, r *http.Request) {
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

	jsonResponse, err := ImageLogic(seed, ARRAY_SIZE, REQ_NUM)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	// log.Printf("Request processed in %v\n", executionTime)
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

func ImageLogic(seed int, ARRAY_SIZE int, REQ_NUM int) ([]byte, error) {
	start := time.Now().UnixMicro()

	rand.Seed(int64(seed))
	
	// Load an example image
	fileNames := []string{"Resources/img1.jpg", "Resources/img2.jpg"}
	selectedFile := fileNames[rand.Intn(len(fileNames))]

	file, err := os.Open(selectedFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Add random seed to every pixel
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()
			r = clamp(r + uint32(rand.Intn(256)))
			g = clamp(g + uint32(rand.Intn(256)))
			b = clamp(b + uint32(rand.Intn(256)))
			newColor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			newImg.Set(x, y, newColor)
		}
	}

	// Resize the image
	newImg = resize(newImg, ARRAY_SIZE)

	// Sum all pixel values
	sum := sumPixels(newImg)

	// Flip horizontally
	newImg = flipHorizontally(newImg)
	sum += sumPixels(newImg)

	// Rotate 90 degrees
	newImg = rotate(newImg, 90)
	sum += sumPixels(newImg)

	executionTime := time.Now().UnixMicro() - start

	response := map[string]interface{}{
		"sum":           sum,
		"executionTime": executionTime, // Include raw execution time in microseconds
		"requestNumber": REQ_NUM,
		"arraysize":     ARRAY_SIZE,
	}

	gogcValue := os.Getenv("GOGC")
	gomemlimitValue := os.Getenv("GOMEMLIMIT")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response["heapAlloc"] = m.HeapAlloc
	// response["heapSys"] = m.HeapSys
	response["heapIdle"] = m.HeapIdle
	// response["heapInuse"] = m.HeapInuse
	response["NextGC"] = m.NextGC
	response["NumGC"] = m.NumGC
	response["GOGC"] = gogcValue
	response["GOMEMLIMIT"] = gomemlimitValue
	jsonResponse, err := json.Marshal(response)
	return jsonResponse, err
}

// Implement a basic nearest-neighbor resizing algorithm
func resize(img image.Image, newSize int) *image.RGBA {
    srcBounds := img.Bounds()
    dstBounds := image.Rect(0, 0, newSize, newSize)
    newImg := image.NewRGBA(dstBounds)

    xRatio := float64(srcBounds.Dx()) / float64(newSize)
    yRatio := float64(srcBounds.Dy()) / float64(newSize)

    for y := 0; y < newSize; y++ {
        for x := 0; x < newSize; x++ {
            srcX := int(float64(x) * xRatio)
            srcY := int(float64(y) * yRatio)
            newImg.Set(x, y, img.At(srcX, srcY))
        }
    }

    return newImg
}

func sumPixels(img image.Image) int64 {
	var sum int64 = 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			sum += int64(r + g + b)
		}
	}
	return sum
}

func flipHorizontally(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(bounds.Max.X-x, y, img.At(x, y))
		}
	}
	return flipped
}

func rotate(img image.Image, angle int) *image.RGBA {
	bounds := img.Bounds()
	rotated := image.NewRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			rotated.Set(bounds.Dx()-x-1, y, img.At(x, y))
		}
	}
	return rotated
}

func clamp(value uint32) uint32 {
	if value > 255 {
		return 255
	}
	return value
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
		"arraysize":     ARRAY_SIZE,
	}

	gogcValue := os.Getenv("GOGC")
	gomemlimitValue := os.Getenv("GOMEMLIMIT")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response["heapAlloc"] = m.HeapAlloc
	// response["heapSys"] = m.HeapSys
	response["heapIdle"] = m.HeapIdle
	// response["heapInuse"] = m.HeapInuse
	response["NextGC"] = m.NextGC
	response["NumGC"] = m.NumGC
	response["GOGC"] = gogcValue
	response["GOMEMLIMIT"] = gomemlimitValue
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
