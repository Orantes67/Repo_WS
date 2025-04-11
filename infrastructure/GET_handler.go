package infrastructure

import (
	"log"
	"net/http"
	"sockets-go/application"

	"github.com/gin-gonic/gin"
)

// func WSHandler(ctx *gin.Context) {

// 	var upgrader websocket.Upgrader

// 	if websocket.IsWebSocketUpgrade(ctx.Request) {
// 		upgrader = websocket.Upgrader{
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true
// 			},
// 		}
// 	} else {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Cliente no aceptado"})
// 	}
// 	/* var upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	} */

// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

// 	if err != nil {
// 		log.Printf("Error de conexión: %v", err)
// 		return
// 	}

// 	defer conn.Close()

// 	log.Println("Cliente conectado")

// 	for {

// 		messageType, p, err := conn.ReadMessage()

// 		if err != nil {
// 			log.Printf("Error en la lectura: %v", err)
// 			break
// 		}

// 		log.Printf("Recibido: %s, de tipo %d", p, messageType)

// 		conn.WriteMessage(1, p)

// 		time.Sleep(16 * time.Millisecond)
// 	}

// 	log.Println("Cliente desconectado")

// }

type WebsocketHandler struct {
	wsService application.WebsocketService
}

func NewWebsocketHandler(
	appService application.WebsocketService,
) *WebsocketHandler {
	return &WebsocketHandler{
		wsService: appService,
	}
}

func (wsH *WebsocketHandler) Upgrade(ctx *gin.Context) {
	userID := ctx.Query("user_id")

	log.Println(userID)

	err := wsH.wsService.HandleConnection(ctx.Writer, ctx.Request, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to upgrade"})
	}
}
