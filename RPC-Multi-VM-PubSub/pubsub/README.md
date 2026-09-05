# Subscriber Service (VM4)

This service receives high-memory alerts published by the Web Server.

Responsibilities:
- Listen for events from the Publisher
- Display memory alerts

## Requirements

- Go installed

## Run

```bash
go run subscriber.go
```

The subscriber ip is **192.168.1.100** listens on:

```text
TCP :9090
```

## Expected Output

```text
===================================
HIGH MEMORY USAGE ALERT RECEIVED
Service: web-server
Memory Usage: 345 MB
Threshold: 300 MB
Timestamp: ...
===================================
```

## Test

Start the Subscriber first.

Then run the Web Server and trigger memory allocation:

```bash
curl "http://VM1_IP:8080/consume-memory?mb=400"
```

After the memory threshold is exceeded, an alert should appear in the Subscriber console.