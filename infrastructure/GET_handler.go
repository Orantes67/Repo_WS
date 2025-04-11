package infrastructure

import (
	"net/http"
	"sockets-go/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	notificationService *application.NotificationService
}

func NewProductHandler(notificationService *application.NotificationService) *ProductHandler {
	return &ProductHandler{
		notificationService: notificationService,
	}
}

func (ph *ProductHandler) UpdateStock(ctx *gin.Context) {
	productID := ctx.Param("id")
	stockStr := ctx.PostForm("stock")
	
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
		return
	}

	// Check if stock is low and notify
	if stock <= 5 {
		ph.notificationService.NotifyLowStock(productID, stock)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Stock updated successfully",
		"product_id": productID,
		"stock": stock,
	})
}

type OrderHandler struct {
	notificationService *application.NotificationService
}

func NewOrderHandler(notificationService *application.NotificationService) *OrderHandler {
	return &OrderHandler{
		notificationService: notificationService,
	}
}

// CreateOrder handles new order creation
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	orderID := ctx.PostForm("order_id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	
	amountStr := ctx.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Generate a products URL for this order
	productsURL := "/api/orders/" + orderID + "/products"

	// Send notification for new order
	oh.notificationService.NotifyNewOrder(orderID, amount, productsURL)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order_id": orderID,
		"amount": amount,
	})
}

// CancelOrder handles order cancellation
func (oh *OrderHandler) CancelOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	amountStr := ctx.PostForm("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	provider := ctx.PostForm("provider")

	// Send notification for canceled order
	oh.notificationService.NotifyCanceledOrder(orderID, amount, provider)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order canceled successfully",
		"order_id": orderID,
	})
}





