package application

import (
	"log"
	"net/http"
	"sockets-go/domain"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketHandlerService interface {
	HandleWebsocketConnection(w http.ResponseWriter, r *http.Request, sessionID string) error
}

type WebsocketService struct {
	upgrader      websocket.Upgrader
	clients       map[*websocket.Conn]bool // Nuevo mapa para clientes WebSocket
	clientsMutex  sync.RWMutex
	sessions      map[string]*domain.Session
	sessionsMutex sync.RWMutex
	nextSessionID int
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:       make(map[*websocket.Conn]bool),
		clientsMutex:  sync.RWMutex{},
		sessions:      make(map[string]*domain.Session),
		sessionsMutex: sync.RWMutex{},
		nextSessionID: 1,
	}
}

// RegisterClient registra un nuevo cliente WebSocket
func (ws *WebsocketService) RegisterClient(conn *websocket.Conn) *websocket.Conn {
	ws.clientsMutex.Lock()
	defer ws.clientsMutex.Unlock()
	ws.clients[conn] = true
	return conn
}

// UnregisterClient elimina un cliente WebSocket
func (ws *WebsocketService) UnregisterClient(conn *websocket.Conn) {
	ws.clientsMutex.Lock()
	defer ws.clientsMutex.Unlock()
	delete(ws.clients, conn)
	conn.Close()
}

// Broadcast envía un mensaje a todos los clientes conectados
func (ws *WebsocketService) Broadcast(message []byte) {
	ws.clientsMutex.RLock()
	defer ws.clientsMutex.RUnlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
		}
	}
}

// El resto de los métodos permanecen igual...
func (ws *WebsocketService) HandleConnection(
	w http.ResponseWriter, r *http.Request, sessionID string,
) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	log.Println("New connection with session ID:", sessionID)

	if sessionID == "" {
		sessionID = ws.generateSessionID()
		log.Println("Generated new session ID:", sessionID)
	}

	defer conn.Close()

	session := domain.NewSession(conn, sessionID, ws.sessions)

	ws.addSession(sessionID, session)

	session.StartHandling(ws.removeSession)

	return nil
}

func (ws *WebsocketService) generateSessionID() string {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	id := ws.nextSessionID
	ws.nextSessionID++
	return strconv.Itoa(id)
}

func (ws *WebsocketService) addSession(sessionID string, session *domain.Session) {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	ws.sessions[sessionID] = session
	log.Printf("Session %s added, total sessions: %d", sessionID, len(ws.sessions))
}

func (ws *WebsocketService) removeSession(sessionID string) {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	delete(ws.sessions, sessionID)
	log.Printf("Session %s removed, remaining sessions: %d", sessionID, len(ws.sessions))
}

// GetSessions returns the current active sessions
func (ws *WebsocketService) GetSessions() map[string]*domain.Session {
	ws.sessionsMutex.RLock()
	defer ws.sessionsMutex.RUnlock()
	return ws.sessions
}