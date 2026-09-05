# Web Server (VM1)

This service provides the web interface of the distributed system.

Responsibilities:
- Display login page
- Send authentication requests to Auth Server using JSON-RPC
- Display the image served by the File Server
- Monitor memory usage
- Publish high-memory alerts to the Subscriber

## Requirements

- Go installed
- Auth Server running
- File Server running
- Subscriber running (for Part 3)

## Run

```bash
go run main.go
```

The web server will start on:

```text
http://192.168.1.110:8080
```

## Test Login

Use one of the users defined in `auth-vm/users.json`.

## Memory Consumption Test

Increase memory usage:

```bash
curl "http://192.168.1.110:8080/consume-memory?mb=100"
```

Repeat multiple times until the memory threshold is exceeded.