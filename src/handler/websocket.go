package handler

import (
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/session"
	"github.com/gabriel-panz/gojam/utils"
	"github.com/gorilla/websocket"
)

func Ws(ss session.SessionService, up websocket.Upgrader, logger *log.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("session_id")

		if id == "" {
			log.Printf("No connection provided.\n")
			w.WriteHeader(404)
			return
		}

		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		s, err := ss.GetSession(id)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}

		a := utils.GetAuthorizedUser(r)

		s.ConnectToSession(a.Token, ws)
	})
}
