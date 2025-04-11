package infrastructure

import (
	"sockets-go/application"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	wsService := application.NewWebsocketService()

	wsHandler := NewWebsocketHandler(*wsService)

	ws_group := engine.Group("ws")

	ws_group.GET("handshake", wsHandler.Upgrade)

}
