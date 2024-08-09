package main

import (
	"log"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
)

func SubscribeToNewHeads(ws *websocket.Conn) (string, error) {
	// Create a subscription request
	request := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_subscribe",
		Params:  []interface{}{"newHeads"},
		ID:      1,
	}

	// Send the subscription request
	if err := websocket.JSON.Send(ws, request); err != nil {
		return "", fmt.Errorf("failed to send subscription request: %v", err)
	}

	// Read the response to get the subscription ID
	var response RPCResponse
	if err := websocket.JSON.Receive(ws, &response); err != nil {
		return "", fmt.Errorf("failed to receive subscription response: %v", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("subscription error: %s", response.Error.Message)
	}

	var subscriptionID string
	if err := json.Unmarshal(response.Result, &subscriptionID); err != nil {
		return "", fmt.Errorf("failed to parse subscription ID: %v", err)
	}

	return subscriptionID, nil
}

func ListenForNewHeads(ws *websocket.Conn, subscriptionID string, parser *EthParser) {
	for {
		var message map[string]interface{}
		if err := websocket.JSON.Receive(ws, &message); err != nil {
			log.Printf("Failed to receive message: %v", err)
			break
		}

		// Check if the message is a subscription update
		if message["method"] == "eth_subscription" {
			params := message["params"].(map[string]interface{})
			if params["subscription"] == subscriptionID {
				result := params["result"].(map[string]interface{})
				blockNumberHex := result["number"].(string)
				go parser.GetTransactionsForBlock(blockNumberHex)
			}
		}
	}
}