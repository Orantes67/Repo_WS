package infrastructure

import (
	"sockets-go/application"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketHandler struct {
	wsService *application.WebsocketService
}

func NewWebsocketHandler(wsService application.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{
		wsService: &wsService,
	}
}

func (wh *WebsocketHandler) Upgrade(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Register the new client connection
	wh.wsService.RegisterClient(conn)
	defer wh.wsService.UnregisterClient(conn)

	// Start listening for messages from this client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", string(p))

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}