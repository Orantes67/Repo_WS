package domain

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	conn         *websocket.Conn
	SessionID    string
	closeHandler func(sessionID string)
	sessions     map[string]*Session
}

func NewSession(conn *websocket.Conn, userID string, sessions map[string]*Session) *Session {
	return &Session{
		conn:      conn,
		SessionID: userID,
		sessions:  sessions,
	}
}

func (s *Session) SetCloseHandler(handler func(sessionID string)) {
	s.closeHandler = handler
}

func (s *Session) StartHandling(removeSession func(sessionID string)) {
	//log.Println(s.SessionID)
	s.closeHandler = removeSession
	s.readPump()
	//s.writePump()
}

func (s *Session) readPump() {
	defer func() {
		s.conn.Close()
		if s.closeHandler != nil {
			s.closeHandler(s.SessionID)
		}
	}()

	//go s.writePump()

	for {

		messageType, _, err := s.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("Error %v", err)
				break
			}
		}

		if messageType != -1 {
			//log.Printf("Recibido: %s, de tipo %d desde seesion %s", p, messageType, s.SessionID)

			//s.broadcast(1, p)

			message := fmt.Sprintf("message from client %s", s.SessionID)

			session := s.sessions["consumer"]

			session.SendMessage(messageType, []byte(message))

		}

		time.Sleep(17 * time.Millisecond)
	}

	//select {}
}

/* func (s *Session) writePump() {
	defer func() {
		s.conn.Close()
	}()

	for {

		messageType := websocket.TextMessage
		message := []byte("Message from server")

		err := s.conn.WriteMessage(messageType, message)

		if err != nil {
			log.Println("Write error: ", err)
			break
		}

		time.Sleep(10 * time.Second)
	}

	select {}
} */

func (s *Session) broadcast(messageType int, payloadbyte []byte) {
	err := s.conn.WriteMessage(messageType, payloadbyte)
	if err != nil {
		log.Println("Broadcast error: ", err)
	}
}

func (s *Session) SendMessage(messageType int, payloadByte []byte) {
	err := s.conn.WriteMessage(messageType, payloadByte)

	if err != nil {
		log.Printf("Error en SendMessage de la sesión %s: %v", s.SessionID, err)
	}
}
