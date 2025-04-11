package notification_client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	LowStockNotification    NotificationType = "low_stock"
	NewOrderNotification    NotificationType = "new_order"
	CancelOrderNotification NotificationType = "cancel_order"
)

// Notification represents a system notification
type Notification struct {
	Type        NotificationType `json:"type"`
	Message     string           `json:"message"`
	Timestamp   time.Time        `json:"timestamp"`
	EntityID    string           `json:"entity_id"`
	Amount      float64          `json:"amount,omitempty"`
	StockLevel  int              `json:"stock_level,omitempty"`
	Provider    string           `json:"provider,omitempty"`
	ProductsURL string           `json:"products_url,omitempty"`
}

func main() {
	// Create a channel to handle signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to the WebSocket server
	client := url.URL{
		Scheme:   "ws",
		Host:     "localhost:4000",
		Path:     "/ws/handshake",
		RawQuery: "user_id=notification-client",
	}

	fmt.Printf("Connecting to %s\n", client.String())
	conn, _, err := websocket.DefaultDialer.Dial(client.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a channel to receive messages from the server
	done := make(chan struct{})

	// Start a goroutine to read messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			
			// Try to parse the message as a notification
			var notification Notification
			if err := json.Unmarshal(message, &notification); err != nil {
				// Not a notification, just print the raw message
				log.Printf("Received message: %s", message)
				continue
			}
			
			// Handle different notification types
			switch notification.Type {
			case LowStockNotification:
				fmt.Printf("\n⚠️ LOW STOCK ALERT ⚠️\n")
				fmt.Printf("Product ID: %s\n", notification.EntityID)
				fmt.Printf("Current Stock: %d units\n", notification.StockLevel)
				fmt.Printf("Timestamp: %s\n\n", notification.Timestamp.Format(time.RFC1123))
				
			case NewOrderNotification:
				fmt.Printf("\n🛒 NEW ORDER CREATED 🛒\n")
				fmt.Printf("Order ID: %s\n", notification.EntityID)
				fmt.Printf("Total Amount: $%.2f\n", notification.Amount)
				fmt.Printf("Products: %s\n", notification.ProductsURL)
				fmt.Printf("Timestamp: %s\n\n", notification.Timestamp.Format(time.RFC1123))
				
			case CancelOrderNotification:
				fmt.Printf("\n❌ ORDER CANCELED ❌\n")
				fmt.Printf("Order ID: %s\n", notification.EntityID)
				fmt.Printf("Amount: $%.2f\n", notification.Amount)
				if notification.Provider != "" {
					fmt.Printf("Provider: %s\n", notification.Provider)
				}
				fmt.Printf("Timestamp: %s\n\n", notification.Timestamp.Format(time.RFC1123))
				
			default:
				fmt.Printf("Unknown notification type: %s\n", notification.Type)
			}
		}
	}()

	// Periodically send a heartbeat message
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Send a heartbeat message
			err := conn.WriteMessage(websocket.TextMessage, []byte("heartbeat"))
			if err != nil {
				log.Println("Error sending heartbeat:", err)
				return
			}
		case <-interrupt:
			// Close the connection gracefully
			log.Println("Interrupt received, closing connection...")
			err := conn.WriteMessage(websocket.CloseMessage, 
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}