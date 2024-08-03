package session

import (
	"log"

	"github.com/gabriel-panz/gojam/types"
	"github.com/gorilla/websocket"
)

type Session struct {
	ID      string
	DjToken string
	// map where the key is the [connection token]
	Conns map[string]*SessionClient
}

type SessionClient struct {
	DeviceId string
	Conn     *websocket.Conn
}

func NewSession(id string, creatorToken string) *Session {
	return &Session{
		DjToken: creatorToken,
		Conns:   make(map[string]*SessionClient),
		ID:      id,
	}
}

func (s *Session) ConnectToSession(token string, conn *websocket.Conn) {
	s.Conns[token] = &SessionClient{
		DeviceId: "",
		Conn:     conn,
	}
}

type SessionMessage struct {
	DeviceID string         `json:"device_id"`
	Uri      string         `json:"uri,omitempty"`
	Type     types.ShowType `json:"type,omitempty"`
}

func Broadcast(m SessionMessage, ws *websocket.Conn) {
	log.Println(m)
	err := ws.WriteJSON(m)
	if err != nil {
		log.Println("write error:", err)
	}
}
