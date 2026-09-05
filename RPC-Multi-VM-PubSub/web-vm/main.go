package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/rpc/jsonrpc"
	"runtime"
	"strconv"
	"time"
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Success bool
	Message string
}

type MemoryEvent struct {
	EventType  string `json:"event_type"`
	Service    string `json:"service"`
	MemoryMB   uint64 `json:"memory_mb"`
	Threshold  uint64 `json:"threshold_mb"`
	Timestamp  string `json:"timestamp"`
}

var memoryHog [][]byte

const memoryThresholdMB = 300

const subscriberAddress = "192.168.1.100:9090"

func loginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	conn, err := jsonrpc.Dial("tcp", "192.168.1.111:9000")
	if err != nil {
		http.Error(w, "Cannot connect to auth server", 500)
		return
	}

	req := LoginRequest{
		Username: username,
		Password: password,
	}

	var res LoginResponse

	err = conn.Call("AuthService.Login", req, &res)
	if err != nil {
		http.Error(w, "RPC call failed", 500)
		return
	}

	if !res.Success {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/success.html"))

	data := map[string]string{
		"Username": username,
		"ImageURL": "http://192.168.1.108:8081/files/image.jpg",
	}

	tmpl.Execute(w, data)
}

func consumeMemoryHandler(w http.ResponseWriter, r *http.Request) {
	mbStr := r.URL.Query().Get("mb")

	if mbStr == "" {
		mbStr = "100"
	}

	mb, err := strconv.Atoi(mbStr)
	if err != nil {
		http.Error(w, "Invalid mb parameter", 400)
		return
	}

	block := make([]byte, mb*1024*1024)

	for i := range block {
		block[i] = 1
	}

	memoryHog = append(memoryHog, block)

	fmt.Fprintf(w,
		"Allocated %d MB memory\nCurrent allocations: %d\n",
		mb,
		len(memoryHog),
	)
}

func publishEvent(event MemoryEvent) {
	conn, err := net.Dial("tcp", subscriberAddress)
	if err != nil {
		log.Println("Failed to connect to subscriber:", err)
		return
	}

	defer conn.Close()

	encoder := json.NewEncoder(conn)

	err = encoder.Encode(event)
	if err != nil {
		log.Println("Failed to publish event:", err)
		return
	}

	log.Println("Memory event published")
}

func monitorMemory() {
	for {
		var mem runtime.MemStats

		runtime.ReadMemStats(&mem)

		memoryMB := mem.Alloc / 1024 / 1024

		log.Printf("Current memory usage: %d MB\n", memoryMB)

		if memoryMB > memoryThresholdMB {
			event := MemoryEvent{
				EventType: "HIGH_MEMORY_USAGE",
				Service: "web-server",
				MemoryMB: memoryMB,
				Threshold: memoryThresholdMB,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			publishEvent(event)
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	go monitorMemory()

	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consume-memory", consumeMemoryHandler)

	log.Println("Web server running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}