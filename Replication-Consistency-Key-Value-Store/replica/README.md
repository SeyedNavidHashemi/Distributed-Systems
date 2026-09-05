# **Replicated Key-Value Store (Replica Service)**

This folder contains the core implementation of the Replicated Key-Value Store for the third homework of the Distributed Computing course at the University of Tehran. It handles data replication, conflict resolution, and supports both Eventual and Simplified Strong consistency models.

**Soroush Esfahanian** - **Navid Hashemi**

## **Project Structure**
.
├── configs/
│   ├── replica1.json       # Port and peer configurations for Replica 1
│   ├── replica2.json       # Port and peer configurations for Replica 2
│   └── replica3.json       # Port and peer configurations for Replica 3
├── replica/
│   ├── main.go             # Core Go program for running the replica service
│   └── README.md           # This README file

## **File Descriptions**
**main.go:**
The main Go program that acts as an independent network service. It provides HTTP endpoints for `PUT` and `GET` requests from clients, and a `/sync` endpoint for internal replication between nodes. It supports dynamic configuration for network delays and consistency models via command-line flags.

**configs/*.json:**
JSON configuration files that supply each replica with its unique ID, local port, and the URLs of its peers. This ensures replicas do not share internal state and communicate strictly over the network.

## **How to Run the Replicas**
To test the system, you need to run at least three separate instances of `main.go`, each in its own terminal.

1. Navigate to the root directory of the project (`HW3/`).
2. Open three separate terminal windows.
3. Run the following commands to start the cluster in the default **Eventual Consistency** mode:

**Terminal 1:**
go run replica/main.go -config=configs/replica1.json

**Terminal 2:**
go run replica/main.go -config=configs/replica2.json

**Terminal 3:**
go run replica/main.go -config=configs/replica3.json

### **Advanced Configurations (Flags)**
You can alter the behavior of the replicas using command-line flags:

*   **`-consistency`**: Switch between consistency models (default: `eventual`).
    *   *Example:* `-consistency=strong`
*   **`-delay`**: Inject artificial network delay in milliseconds to observe convergence and stale reads.
    *   *Example:* `-delay=500`

**Example - Running with Strong Consistency and 500ms Delay:**
go run replica/main.go -config=configs/replica1.json -consistency=strong -delay=500
