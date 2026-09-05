# **Testing & Metrics Client**

This folder contains the automated client script designed to benchmark the Replicated Key-Value Store. It calculates the necessary metrics required for the final report's performance table, including latency, convergence time, and stale read counts.

**Soroush Esfahanian 810101376** - **Navid Hashemi 810101549**

## **Project Structure**
.
├── client/
│   ├── main.go             # Go script for automating Scenario 4 metrics
│   └── README.md           # This README file
├── results/
│   ├── scenario1.txt       # Captured logs for Temporary Inconsistency
│   ├── scenario2.txt       # Captured logs for Node Failures
│   ├── scenario3.txt       # Captured logs for Concurrent Conflicts
│   └── scenario4.txt       # Captured metrics from the client program

## **File Descriptions**
**client/main.go:**
An automated testing script that interacts with the running replicas over HTTP. It performs a specific sequence of `PUT` and `GET` operations across different nodes to measure:
*   **PUT Latency:** Time taken to receive an acknowledgment after writing data.
*   **Stale Reads:** The number of times old/missing data is read from a secondary replica before synchronization finishes.
*   **Convergence Time:** The total time taken for all nodes to become consistent.

**results/*.txt:**
Text files storing the raw output and command logs from the 4 testing scenarios outlined in the assignment.

## **How to Run the Client**
Before running the client, ensure that your three replicas are currently running (refer to `replica/README.md` for startup instructions).

1. Navigate to the root directory of the project (`HW3/`).
2. Run the client script:
go run client/main.go

3. The script will output the calculated metrics directly to your terminal. Record these metrics to fill out the final evaluation table in the PDF report.

**Note:** To thoroughly test Scenario 4, you should restart your replicas with different `-delay` values (e.g., 0, 500, 2000) and `-consistency` models before re-running this client script.
