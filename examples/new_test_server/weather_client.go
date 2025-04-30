package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// ToolParams represents the parameters for a tool request
type ToolParams struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments"`
}

// StateAlertRequest represents the input for the weather alerts tool
type StateAlertRequest struct {
	State string `json:"state"`
}

// Response represents a JSON-RPC response
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ErrorObject    `json:"error,omitempty"`
}

// ErrorObject represents a JSON-RPC error
type ErrorObject struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run weather_client.go <state>")
		fmt.Println("Example: go run weather_client.go CA")
		os.Exit(1)
	}

	state := strings.ToUpper(os.Args[1])
	fmt.Printf("Getting weather alerts for %s...\n\n", state)

	// Create the request
	req := Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/use",
		Params: ToolParams{
			Name: "get_weather_alerts",
			Arguments: StateAlertRequest{
				State: state,
			},
		},
	}

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "http://localhost:1234/sse", strings.NewReader(string(reqBody)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Parse the response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	// Check for errors
	if response.Error != nil {
		log.Fatalf("Error: %s", response.Error.Message)
	}

	// Pretty print the result
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, response.Result, "", "  "); err != nil {
		log.Fatalf("Failed to format JSON: %v", err)
	}

	fmt.Println(prettyJSON.String())
}
