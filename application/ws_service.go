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
		sessions:      make(map[string]*domain.Session),
		sessionsMutex: sync.RWMutex{},
		nextSessionID: 1,
	}
}

func (ws *WebsocketService) HandleConnection(
	w http.ResponseWriter, r *http.Request, sessionID string,
) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	log.Println(sessionID)

	if sessionID == "" {
		log.Println("Condición correcta")
		sessionID = ws.generateSessionID()
	}

	log.Println(sessionID)

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
}

func (ws *WebsocketService) removeSession(sessionID string) {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	delete(ws.sessions, sessionID)

	log.Printf("Session %s Removed", sessionID)
}
