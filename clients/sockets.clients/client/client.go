package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	go startClient()
	select {}
}

func startClient() {
	client := url.URL{
		Scheme: "ws",
		Host:   "localhost:4000",
		Path:   "/ws/handshake",
		//RawQuery: fmt.Sprintf("user_id"),
	}

	conn, _, err := websocket.DefaultDialer.Dial(client.String(), nil)

	defer conn.Close()

	if err != nil {
		fmt.Printf("error en la conexión %v", err)
	}

	for {
		message := fmt.Sprintf("Message from client add %v", time.Now().Format(time.RFC3339))

		err := conn.WriteMessage(1, []byte(message))

		if err != nil {
			fmt.Printf("Error en enviar mensaje")
		}

		fmt.Printf("send %s", message)
		time.Sleep(2 * time.Second)
	}
}
