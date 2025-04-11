package application

import (
	"log"
	"sockets-go/domain"
	"sync"
)

// NotificationService handles system notifications
type NotificationService struct {
	websocketService *WebsocketService
	mutex            sync.RWMutex
}

// NewNotificationService creates a new notification service
func NewNotificationService(wsService *WebsocketService) *NotificationService {
	return &NotificationService{
		websocketService: wsService,
		mutex:            sync.RWMutex{},
	}
}

// NotifyLowStock sends a notification when a product's stock is low
func (ns *NotificationService) NotifyLowStock(productID string, stockLevel int) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	if stockLevel <= 5 {
		notification := domain.NewLowStockNotification(productID, stockLevel)
		ns.broadcastNotification(notification)
		log.Printf("Low stock notification for product %s with stock level %d", productID, stockLevel)
	}
}

// NotifyNewOrder sends a notification when a new order is created
func (ns *NotificationService) NotifyNewOrder(orderID string, amount float64, productsURL string) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	notification := domain.OrderNotification(orderID, amount, productsURL)
	ns.broadcastNotification(notification)
	log.Printf("New order notification for order %s with amount %.2f", orderID, amount)
}

// NotifyCanceledOrder sends a notification when an order is canceled
func (ns *NotificationService) NotifyCanceledOrder(orderID string, amount float64, provider string) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	notification := domain.NewCancelOrderNotification(orderID, amount, provider)
	ns.broadcastNotification(notification)
	log.Printf("Order canceled notification for order %s with amount %.2f and provider %s", 
		orderID, amount, provider)
}

// broadcastNotification sends a notification to all connected clients
func (ns *NotificationService) broadcastNotification(notification *domain.Notification) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	// Get a session to broadcast (any session can broadcast to all)
	if len(ns.websocketService.sessions) > 0 {
		for _, session := range ns.websocketService.sessions {
			session.BroadcastNotification(notification)
			break
		}
	} else {
		log.Println("No active sessions to broadcast notification")
	}
}