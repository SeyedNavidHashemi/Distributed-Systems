package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
)

const reqPipe = "/tmp/dist_req.pipe"
const resPipe = "/tmp/dist_res.pipe"

type Response struct {
    Status string `json:"status"`
    Result int    `json:"result,omitempty"`
    Error  string `json:"error,omitempty"`
}

func pipeExists(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return (info.Mode() & os.ModeNamedPipe) != 0
}

func main() {

    // Check pipe existence BEFORE open
    if !pipeExists(reqPipe) || !pipeExists(resPipe) {
        fmt.Println("ERR worker_not_running: worker process is not running, please run it first")
        return
    }

    req, err := os.OpenFile(reqPipe, os.O_WRONLY, os.ModeNamedPipe)
    if err != nil {
        fmt.Println("ERR worker_not_running: worker process is not running")
        return
    }
    defer req.Close()

    res, err := os.OpenFile(resPipe, os.O_RDONLY, os.ModeNamedPipe)
    if err != nil {
        fmt.Println("ERR worker_not_running: worker process is not running")
        return
    }
    defer res.Close()

    reqWriter := bufio.NewWriter(req)
    resReader := bufio.NewScanner(res)

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            return
        }

        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        parts := strings.Fields(line)
        op := strings.ToUpper(parts[0])

        if op != "ADD" && op != "SUB" && op != "MUL" && op != "DIV" && op != "POW" {
            fmt.Println("ERR unknown_command: unsupported operation")
            continue
        }

		if len(parts) != 3 {
            fmt.Println("ERR invalid_argument_count: expected exactly two arguments")
            continue
        }

        a, errA := strconv.Atoi(parts[1])
        b, errB := strconv.Atoi(parts[2])
        if errA != nil || errB != nil {
            fmt.Println("ERR non_numeric_argument: arguments must be integers")
            continue
        }

        if op == "DIV" && b == 0 {
            fmt.Println("ERR division_by_zero: cannot divide by zero")
            continue
        }

        // Send request to worker
        _, err = fmt.Fprintf(reqWriter, "%s %d %d\n", op, a, b)
        if err != nil {
            fmt.Println("ERR pipe_write_error: failed to write to pipe")
            return
        }
        reqWriter.Flush()

        if !resReader.Scan() {
            fmt.Println("ERR pipe_read_error: connection lost or worker closed unexpectedly")
            return
        }

        raw := resReader.Text()
        if strings.TrimSpace(raw) == "" {
            fmt.Println("ERR invalid_response_format: empty response from worker")
            continue
        }

        var resp Response
        if err := json.Unmarshal([]byte(raw), &resp); err != nil {
            fmt.Println("ERR invalid_response_format: worker returned invalid JSON")
            continue
        }

        if resp.Status == "ERR" {
            fmt.Println("ERR", resp.Error)
        } else {
            fmt.Println("OK", resp.Result)
        }
    }
}
