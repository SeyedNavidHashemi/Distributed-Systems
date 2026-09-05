# Replication-Consistency-Key-Value-Store

This assignment focuses on the practical concepts of **Replication** and **Consistency** in distributed systems.

The project implements a simple **Replicated Key-Value Store** consisting of multiple independent replicas. Each replica runs as a separate service and communicates with other replicas over the network.

## Project Structure

```text
HW3/
├── report.pdf
├── replica/
│   ├── main.go
│   └── README.md
├── client/
│   ├── main.go
│   └── README.md
├── configs/
│   ├── replica1.json
│   ├── replica2.json
│   └── replica3.json
└── results/
    ├── scenario1.txt
    ├── scenario2.txt
    ├── scenario3.txt
    └── scenario4.txt
```

## System Overview

The system contains at least three independent replicas:

```text
                 ┌───────────┐
                 │  Client   │
                 └─────┬─────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
     ┌─────────┐  ┌─────────┐  ┌─────────┐
     │Replica 1│◄►│Replica 2│◄►│Replica 3│
     └─────────┘  └─────────┘  └─────────┘
```

The replicas support `GET` and `PUT` operations and exchange updates through network communication. The project examines both **Eventual Consistency** and a simplified **Strong Consistency** model.

## Experiments

The project evaluates the system through four scenarios:

* Temporary inconsistency and stale reads
* Replica failure
* Concurrent write conflicts
* Network delay and its effect on convergence

The experiments and their results are stored in the `results/` directory.

## Technologies

* Go
* HTTP / JSON-RPC / gRPC / TCP
* Network-based replica communication

The implementation uses one of the communication methods permitted by the assignment.

## Documentation

Detailed instructions for each component are available in their corresponding `README.md` files.

The final analysis and experimental results are provided in `report.pdf`.

