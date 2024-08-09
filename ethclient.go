package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON-RPC request structure
type RPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// JSON-RPC response structure for block number
type RPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage 		`json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

// JSON-RPC error structure
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendRPCRequest sends a JSON-RPC request to the Ethereum node
func SendRPCRequest(url string, method string, params interface{}) (*RPCResponse, error) {
	// Create the RPC request object
	request := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	// Serialize the request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Send the HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("\n Resp: %+v", resp.Body)

	// Parse the response
	var rpcResponse RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if the response contains an error
	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	return &rpcResponse, nil
}