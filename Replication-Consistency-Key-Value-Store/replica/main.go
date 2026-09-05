package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	ID    string   `json:"id"`
	Port  string   `json:"port"`
	Peers []string `json:"peers"`
}

type Record struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Version   int    `json:"version"`
	UpdatedBy string `json:"updated_by"`
}

var (
	store   = make(map[string]Record)
	rwMutex sync.RWMutex

	config      Config
	consistency string
	delay       int
)

func main() {
	configPath := flag.String("config", "", "Path to config file (e.g., configs/replica1.json)")
	flag.StringVar(&consistency, "consistency", "eventual", "Consistency model: 'eventual' or 'strong'")
	flag.IntVar(&delay, "delay", 0, "Artificial delay for replication in ms")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Config file path is required. Use -config flag.")
	}

	// Read from config file
	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatalf("Failed to open config: %v", err)
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&config)

	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/sync", handleSync)

	log.Printf("Starting %s on port %s | Consistency: %s | Delay: %dms | Peers: %v\n", config.ID, config.Port, consistency, delay, config.Peers)
	http.ListenAndServe(":"+config.Port, nil)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]string
	json.NewDecoder(r.Body).Decode(&reqData)
	key, value := reqData["key"], reqData["value"]

	rwMutex.Lock()
	existingRecord, exists := store[key]
	newVersion := 1
	if exists {
		newVersion = existingRecord.Version + 1
	}

	newRecord := Record{Key: key, Value: value, Version: newVersion, UpdatedBy: config.ID}
	store[key] = newRecord
	rwMutex.Unlock()

	jsonData, _ := json.Marshal(newRecord)

	if consistency == "strong" {
		successCount := 1
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, peer := range config.Peers {
			if peer == "" { continue }
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				time.Sleep(time.Duration(delay) * time.Millisecond)
				client := http.Client{Timeout: 5 * time.Second}
				resp, err := client.Post(p+"/sync", "application/json", bytes.NewBuffer(jsonData))
				if err == nil && resp.StatusCode == http.StatusOK {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(peer)
		}
		wg.Wait()

		if successCount < 2 {
			http.Error(w, "Failed to reach quorum", http.StatusInternalServerError)
			return
		}
	} else {
		for _, peer := range config.Peers {
			if peer == "" { continue }
			go func(p string) {
				time.Sleep(time.Duration(delay) * time.Millisecond)
				http.Post(p+"/sync", "application/json", bytes.NewBuffer(jsonData))
			}(peer)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newRecord)
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	var incoming Record
	json.NewDecoder(r.Body).Decode(&incoming)

	rwMutex.Lock()
	defer rwMutex.Unlock()

	current, exists := store[incoming.Key]
	if !exists || incoming.Version > current.Version {
		store[incoming.Key] = incoming
	} else if incoming.Version == current.Version && incoming.UpdatedBy > current.UpdatedBy {
		store[incoming.Key] = incoming
	}
	w.WriteHeader(http.StatusOK)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	rwMutex.RLock()
	record, exists := store[key]
	rwMutex.RUnlock()

	if !exists {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}