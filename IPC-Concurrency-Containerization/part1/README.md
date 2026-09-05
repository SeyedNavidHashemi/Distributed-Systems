# **Part 1 - Distributed Calculator Using Named Pipes (Go)**

This folder contains the code and results for the experiment corresponding to the first part of the first homework of the distributed computing course at the University of Tehran,implements a simple distributed arithmetic system using two separate Go programs:

- `worker.go` — the server/worker process  
- `interface.go` — the client interface that sends requests and receives responses  

Both components communicate through Linux **named pipes (FIFOs)** using a simple line‑based request format and JSON‑encoded responses.

**Soroush Esfahanian 810101376** - **Navid Hashemi 810101549**

---

## 1. How to Run the Worker

Run the Worker first:

```bash
go run worker.go
```

The Worker will:

- Remove old pipes (`/tmp/dist_req.pipe` and `/tmp/dist_res.pipe`)  
- Create new named pipes  
- Wait for incoming requests  
- Perform arithmetic operations  
- Return responses in JSON format  
- Clean up both pipes when terminated (Ctrl+C)

---

## 2. How to Run the Interface

Run the Interface in another terminal:

```bash
go run interface.go
```

If the Worker is not running, the Interface immediately prints:

```
ERR worker_not_running: worker process is not running, please run it first
```

Once connected, you can enter commands such as:

```
> ADD 3 5
> SUB 10 4
> POW 2 8
```

---

## 3. Correct Execution Order

To use the system correctly:

1. Open Terminal #1  
   ```bash
   go run worker.go
   ```

2. Open Terminal #2  
   ```bash
   go run interface.go
   ```

3. Enter operations in the Interface terminal  

When finished, stop the Worker (Terminal #1) with **Ctrl+C** — this deletes the pipes automatically.

Running the Interface before the Worker will result in a `worker_not_running` error.

---

## 4. Example Inputs and Outputs

### Valid Examples

```
> ADD 4 5
OK 9

> SUB 10 3
OK 7

> MUL 6 7
OK 42

> DIV 12 3
OK 4

> POW 2 8
OK 256
```

### Error Examples

Unknown operation:
```
> HELLO 1 2
ERR unknown_command: unsupported operation
```

Incorrect argument count:
```
> ADD 5
ERR invalid_argument_count: expected exactly two arguments
```

Non-numeric input:
```
> ADD X 4
ERR non_numeric_argument: arguments must be integers
```

Division by zero:
```
> DIV 10 0
ERR division_by_zero: cannot divide by zero
```

Worker closed or crashed:
```
ERR pipe_read_error: connection lost or worker closed unexpectedly
```

---

## 5. Message Exchange Protocol (Short Explanation)

The system uses two named pipes for communication:

- **Request Pipe:** `/tmp/dist_req.pipe`  
  Sent from Interface → Worker  

- **Response Pipe:** `/tmp/dist_res.pipe`  
  Sent from Worker → Interface  

### Request Format (plain text)

```
COMMAND operand1 operand2
```

Example:
```
ADD 3 7
```

Commands supported: `ADD`, `SUB`, `MUL`, `DIV`, `POW`

### Response Format (JSON)

Success:
```json
{"status":"OK","result":10}
```

Error:
```json
{"status":"ERR","error":"division_by_zero"}
```

### Error Types Supported

- `unknown_command`
- `invalid_argument_count`
- `non_numeric_argument`
- `division_by_zero`
- `worker_not_running`
- `pipe_write_error`
- `pipe_read_error`
- `invalid_response_format`

---

## Summary

This project demonstrates:

- Inter-process communication (IPC) via FIFOs  
- JSON serialization/deserialization  
- Input validation on both client and server  
- Safe pipe creation, cleanup, and error handling  
- A minimal distributed arithmetic system
