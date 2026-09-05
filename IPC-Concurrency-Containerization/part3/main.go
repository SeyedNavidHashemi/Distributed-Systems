package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ComputeResponse struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}


func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


func computeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	// Read query parameters
	op := r.URL.Query().Get("op")
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	// Check missing parameters
	if op == "" || aStr == "" || bStr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing parameters: op, a, b are required"})
		return
	}

	// Convert to numbers
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "parameter a must be numeric"})
		return
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "parameter b must be numeric"})
		return
	}

	var result float64

	switch op {
	case "add":
		result = a + b
	case "sub":
		result = a - b
	case "mul":
		result = a * b
	case "div":
		if b == 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "division by zero"})
			return
		}
		result = a / b
	default:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("invalid operation: %s", op)})
		return
	}

	resp := ComputeResponse{
		Operation: op,
		A:         a,
		B:         b,
		Result:    result,
	}

	writeJSON(w, http.StatusOK, resp)
}

// Main + server setup
func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/compute", computeHandler)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
