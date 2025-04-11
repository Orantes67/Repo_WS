package infrastructure

import (
	"sockets-go/application"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {
	// Initialize services
	wsService := application.NewWebsocketService()
	notificationService := application.NewNotificationService(wsService)

	// Initialize handlers
	wsHandler := NewWebsocketHandler(*wsService)
	productHandler := NewProductHandler(notificationService)
	orderHandler := NewOrderHandler(notificationService)

	// WebSocket routes
	ws_group := engine.Group("ws")
	ws_group.GET("handshake", wsHandler.Upgrade)

	// API routes
	api := engine.Group("api")
	
	// Product routes
	products := api.Group("products")
	products.POST("/:id/stock", productHandler.UpdateStock)
	
	// Order routes
	orders := api.Group("orders")
	orders.POST("/", orderHandler.CreateOrder)
	orders.POST("/:id/cancel", orderHandler.CancelOrder)
}