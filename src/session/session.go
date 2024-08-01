package session

import (
	"log"

	"github.com/gabriel-panz/gojam/types"
	"github.com/gorilla/websocket"
)

type Session struct {
	ID      string
	DjToken string
	Conns   map[string]*websocket.Conn
}

func NewSession(id string, creatorToken string) *Session {
	return &Session{
		DjToken: creatorToken,
		Conns:   make(map[string]*websocket.Conn),
		ID:      id,
	}
}

func (s *Session) ConnectToSession(token string, conn *websocket.Conn) {
	s.Conns[token] = conn
}

type SessionMessage struct {
	DeviceID string         `json:"device_id"`
	Uri      string         `json:"uri,omitempty"`
	Type     types.ShowType `json:"type,omitempty"`
}

// func (s *Session) readLoop(conn *websocket.Conn) {
// 	for {
// 		m := SessionMessage{}
// 		err := conn.ReadJSON(m)
// 		if err != nil {
// 			if err == io.EOF {
// 				delete(s.Conns, conn)
// 				return
// 			}
// 			log.Println("read error:", err)
// 			continue
// 		}

// 		for ws := range s.Conns {
// 			go Broadcast(m, ws)
// 		}
// 	}
// }

func Broadcast(m SessionMessage, ws *websocket.Conn) {
	log.Println(m)
	err := ws.WriteJSON(m)
	if err != nil {
		log.Println("write error:", err)
	}
}
