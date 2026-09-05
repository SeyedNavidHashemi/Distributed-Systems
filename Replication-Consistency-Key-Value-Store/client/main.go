package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Record struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Version   int    `json:"version"`
	UpdatedBy string `json:"updated_by"`
}

func main() {
	key := "test_key"
	expectedValue := "latency_test_1"

	// Measure PUT Latency
	startPut := time.Now()
	putData := []byte(fmt.Sprintf(`{"key":"%s", "value":"%s"}`, key, expectedValue))
	resp, err := http.Post("http://localhost:8081/put", "application/json", bytes.NewBuffer(putData))
	putLatency := time.Since(startPut)

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("PUT failed: %v\n", err)
		return
	}
	fmt.Printf("PUT Latency: %v\n", putLatency)

	// Measure GET Latency & Stale Reads & Convergence Time
	staleReads := 0
	startGet := time.Now()
	
	for {
		// Read from Replica 2 immediately after writing to Replica 1
		resp, err := http.Get("http://localhost:8082/get?key=" + key)
		if err != nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		var rec Record
		json.NewDecoder(resp.Body).Decode(&rec)
		resp.Body.Close()

		if rec.Value != expectedValue {
			staleReads++
			time.Sleep(50 * time.Millisecond) // Poll every 50ms
		} else {
			convergenceTime := time.Since(startGet)
			fmt.Printf("GET Latency (last successful): ~%v\n", time.Since(startGet)/time.Duration(staleReads+1))
			fmt.Printf("Stale Reads Count: %d\n", staleReads)
			fmt.Printf("Convergence Time: %v\n", convergenceTime)
			break
		}
	}
}