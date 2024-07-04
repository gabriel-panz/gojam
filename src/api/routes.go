package api

import (
	"log"
	"net/http"

	h "github.com/gabriel-panz/gojam/handler"
	m "github.com/gabriel-panz/gojam/middleware"
	"github.com/gabriel-panz/gojam/spotify"
)

func addRoutes(
	mux *http.ServeMux,
	homeHandler h.HomeHandler,
	spotifyService spotify.Service,
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
	mux.HandleFunc("PUT /player/play", m.Authorize(logger, h.Play(logger, spotifyService)))
	mux.HandleFunc("PUT /player/pause", m.Authorize(logger, h.Pause(logger, spotifyService)))

	// Search
	mux.HandleFunc("GET /search", m.Authorize(logger, h.Search(logger, spotifyService)))
}
