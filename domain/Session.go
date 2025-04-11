package domain

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	conn         *websocket.Conn
	SessionID    string
	closeHandler func(sessionID string)
	sessions     map[string]*Session
	mutex        sync.Mutex
}

func NewSession(conn *websocket.Conn, userID string, sessions map[string]*Session) *Session {
	return &Session{
		conn:      conn,
		SessionID: userID,
		sessions:  sessions,
		mutex:     sync.Mutex{},
	}
}

func (s *Session) SetCloseHandler(handler func(sessionID string)) {
	s.closeHandler = handler
}

func (s *Session) StartHandling(removeSession func(sessionID string)) {
	s.closeHandler = removeSession
	s.readPump()
}

func (s *Session) readPump() {
	defer func() {
		s.conn.Close()
		if s.closeHandler != nil {
			s.closeHandler(s.SessionID)
		}
	}()

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
			// Simple echo for client messages
			message := fmt.Sprintf("Received message from client %s", s.SessionID)
			s.SendMessage(websocket.TextMessage, []byte(message))
		}

		time.Sleep(17 * time.Millisecond)
	}
}

// BroadcastToAll sends a message to all connected sessions
func (s *Session) BroadcastToAll(messageType int, payload []byte) {
	for _, session := range s.sessions {
		session.SendMessage(messageType, payload)
	}
}

// BroadcastNotification sends a notification to all connected clients
func (s *Session) BroadcastNotification(notification *Notification) {
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error marshalling notification: %v", err)
		return
	}
	
	s.BroadcastToAll(websocket.TextMessage, payload)
	log.Printf("Broadcast notification: %s", string(payload))
}

func (s *Session) SendMessage(messageType int, payloadByte []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	err := s.conn.WriteMessage(messageType, payloadByte)
	if err != nil {
		log.Printf("Error sending message to session %s: %v", s.SessionID, err)
	}
}