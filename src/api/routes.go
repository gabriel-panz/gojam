package api

import (
	"log"
	"net/http"

	h "github.com/gabriel-panz/gojam/handler"
	m "github.com/gabriel-panz/gojam/middleware"
	"github.com/gabriel-panz/gojam/session"
	"github.com/gabriel-panz/gojam/spotify"
	"github.com/gorilla/websocket"
)

func addRoutes(
	mux *http.ServeMux,
	homeHandler h.HomeHandler,
	spotifyService spotify.Service,
	sessionService session.SessionService,
	up websocket.Upgrader,
) {
	logger := log.Default()

	mux.HandleFunc("GET /", m.Authorize(logger, homeHandler.GetHomePage))
	mux.HandleFunc("GET /callback", homeHandler.AuthCallback)

	// Auth
	mux.HandleFunc("GET /auth/refresh", m.Authorize(logger, h.Refresh(logger, spotifyService)))

	// SpotifyUserData
	mux.HandleFunc("GET /user/playlists", m.Authorize(logger, h.ListPlaylists(logger, spotifyService)))
	mux.HandleFunc("GET /user/devices", m.Authorize(logger, h.ListDevices(logger, spotifyService)))

	// SpotifyPlayer
	mux.HandleFunc("GET /player", m.Authorize(logger, h.Player(logger, spotifyService)))
	mux.HandleFunc("PUT /player/play", m.Authorize(logger, h.Play(logger, spotifyService, sessionService)))
	mux.HandleFunc("PUT /player/pause", m.Authorize(logger, h.Pause(logger, spotifyService)))

	// Search
	mux.HandleFunc("GET /search", m.Authorize(logger, h.Search(logger, spotifyService)))

	// Session
	mux.HandleFunc("POST /session/", m.Authorize(logger, h.CreateSession(logger, sessionService)))
	mux.HandleFunc("GET /session/{session_id}", m.Authorize(logger, h.Session(logger, sessionService)))
	mux.HandleFunc("PUT /session/leave", m.Authorize(logger, h.LeaveSession(logger, sessionService)))

	// WS
	mux.HandleFunc("GET /ws/{session_id}", m.Authorize(logger, h.Ws(sessionService, up, logger)))
}
