package domain


import (
	"encoding/json"
	"time"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	LowStockNotification   NotificationType = "low_stock"
	NewOrderNotification   NotificationType = "new_order"
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

// NewLowStockNotification creates a notification for low stock
func NewLowStockNotification(productID string, stockLevel int) *Notification {
	return &Notification{
		Type:       LowStockNotification,
		Message:    "Product is running low on stock",
		Timestamp:  time.Now(),
		EntityID:   productID,
		StockLevel: stockLevel,
	}
}

// NewOrderNotification creates a notification for a new order
func OrderNotification(orderID string, amount float64, productsURL string) *Notification {
	return &Notification{
		Type:        NewOrderNotification,
		Message:     "New order has been created",
		Timestamp:   time.Now(),
		EntityID:    orderID,
		Amount:      amount,
		ProductsURL: productsURL,
	}
}

// NewCancelOrderNotification creates a notification for a canceled order
func NewCancelOrderNotification(orderID string, amount float64, provider string) *Notification {
	return &Notification{
		Type:      CancelOrderNotification,
		Message:   "Order has been canceled",
		Timestamp: time.Now(),
		EntityID:  orderID,
		Amount:    amount,
		Provider:  provider,
	}
}

// ToJSON converts the notification to JSON bytes
func (n *Notification) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}