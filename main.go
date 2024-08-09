package main

import (
    "log"
    "github.com/gofiber/fiber/v3"
		"fmt"
		"net/url"
		"golang.org/x/net/websocket"
)

const wsURL = "wss://nd-000-364-211.p2pify.com/5b8d22690a57f293b3a1ed8758014e35"

func main() {
	notificationService := NewNotificationService()
	parser := NewEthParser(notificationService)
	// Initialize a new Fiber app
	app := fiber.New()

	// Start the parser
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatalf("Failed to parse WebSocket URL: %v", err)
	}

	// Dial the WebSocket connection
	ws, err := websocket.Dial(u.String(), "", "http://localhost/")
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Subscribe to new block headers
	subscriptionID, err := SubscribeToNewHeads(ws)
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads: %v", err)
	}
	fmt.Printf("Subscribed with ID: %s\n", subscriptionID)

	// Listen for new headers
	go ListenForNewHeads(ws, subscriptionID, parser)

	// GetCurrentBlock
	app.Get("/get_current_block", func(c fiber.Ctx) error {
			res := parser.GetCurrentBlock()
			return c.Status(200).JSON(fiber.Map{"block_number": res})
	})

	// Subscribe Address
	app.Get("/subscribe/:address", func(c fiber.Ctx) error {
		// Send a string response to the client
		address := c.Params("address")
		res := parser.Subscribe(address)
		return c.Status(200).JSON(fiber.Map{"subscribed": res})
	})

	// Get Transactions
	app.Get("/get_transactions/:address", func(c fiber.Ctx) error {
		// Send a string response to the client
		address := c.Params("address")
		address = AddressToHex(address)
		transactions := parser.GetTransactions(address)
		return c.Status(200).JSON(fiber.Map{"transactions": transactions})
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
