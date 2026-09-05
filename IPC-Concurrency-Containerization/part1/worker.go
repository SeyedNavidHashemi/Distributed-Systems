package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const reqPipe = "/tmp/dist_req.pipe"
const resPipe = "/tmp/dist_res.pipe"

type Response struct {
	Status string `json:"status"`
	Result int    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func cleanup() {
	os.Remove(reqPipe)
	os.Remove(resPipe)
}

func compute(op string, a, b int) (int, string) {
	switch op {
	case "ADD":
		return a + b, ""
	case "SUB":
		return a - b, ""
	case "MUL":
		return a * b, ""
	case "DIV":
		if b == 0 {
			return 0, "division_by_zero"
		}
		return a / b, ""
	case "POW":
		result := 1
		for i := 0; i < b; i++ {
			result *= a
		}
		return result, ""
	default:
		return 0, "unknown_command"
	}
}

func main() {
	// Cleanup old pipes if they exist
	cleanup()

	// Ensure pipes are removed on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	// Create new named pipes
	if err := syscall.Mkfifo(reqPipe, 0666); err != nil {
		fmt.Println(`{"status":"ERR","error":"cannot_create_pipe"}`)
		os.Exit(1)
	}
	if err := syscall.Mkfifo(resPipe, 0666); err != nil {
		fmt.Println(`{"status":"ERR","error":"cannot_create_pipe"}`)
		os.Exit(1)
	}

	req, err := os.OpenFile(reqPipe, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		os.Exit(1)
	}
	defer req.Close()

	res, err := os.OpenFile(resPipe, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		os.Exit(1)
	}
	defer res.Close()

	reader := bufio.NewScanner(req)
	writer := bufio.NewWriter(res)

	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		parts := strings.Fields(line)
		if len(parts) != 3 {
			resp, _ := json.Marshal(Response{Status: "ERR", Error: "invalid_argument_count"})
			fmt.Fprintln(writer, string(resp))
			writer.Flush()
			continue
		}

		op := parts[0]
		a, errA := strconv.Atoi(parts[1])
		b, errB := strconv.Atoi(parts[2])

		if errA != nil || errB != nil {
			resp, _ := json.Marshal(Response{Status: "ERR", Error: "non_numeric_argument"})
			fmt.Fprintln(writer, string(resp))
			writer.Flush()
			continue
		}

		result, errCode := compute(op, a, b)
		if errCode != "" {
			resp, _ := json.Marshal(Response{Status: "ERR", Error: errCode})
			fmt.Fprintln(writer, string(resp))
			writer.Flush()
			continue
		}

		resp, _ := json.Marshal(Response{Status: "OK", Result: result})
		fmt.Fprintln(writer, string(resp))
		writer.Flush()
	}

	// If pipe reading ends unexpectedly → cleanup
	cleanup()
}
