package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type MemoryEvent struct {
	EventType  string `json:"event_type"`
	Service    string `json:"service"`
	MemoryMB   uint64 `json:"memory_mb"`
	Threshold  uint64 `json:"threshold_mb"`
	Timestamp  string `json:"timestamp"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var event MemoryEvent

	decoder := json.NewDecoder(conn)

	err := decoder.Decode(&event)
	if err != nil {
		log.Println("Failed to decode event:", err)
		return
	}

	fmt.Println("===================================")
	fmt.Println("HIGH MEMORY USAGE ALERT RECEIVED")
	fmt.Println("Service:", event.Service)
	fmt.Println("Memory Usage:", event.MemoryMB, "MB")
	fmt.Println("Threshold:", event.Threshold, "MB")
	fmt.Println("Timestamp:", event.Timestamp)
	fmt.Println("===================================")
}

func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscriber listening on port 9090")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}