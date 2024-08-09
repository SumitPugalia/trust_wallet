package main

import (
	"fmt"
)

// NotificationService is a mock implementation of NotificationService.
type NotificationService struct{}

func NewNotificationService() *NotificationService{
	return &NotificationService{}
}

// Notify sends a notification (mock implementation).
func (s *NotificationService) Notify(address string, tx Transaction, incoming bool) {
	direction := "outgoing"
	if incoming {
		direction = "incoming"
	}
	
	fmt.Printf("Notification: %s transaction for address %s, tx: %+v", direction, address, tx)
}