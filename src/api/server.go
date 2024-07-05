package api

import (
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/handler"
	"github.com/gabriel-panz/gojam/session"
	"github.com/gabriel-panz/gojam/spotify"
	"github.com/patrickmn/go-cache"
)

func NewServer(addr string, clientId string, c *cache.Cache) *http.Server {
	mux := http.NewServeMux()

	spotifyService := spotify.Service{
		Client: &http.Client{},
	}

	sessionService := session.SessionService{
		Cache: c,
	}

	homeHandler := handler.HomeHandler{
		Spotify: &spotifyService,
		Handler: handler.Handler{
			Logger: log.Default(),
		},
	}

	addRoutes(mux, homeHandler, spotifyService, sessionService)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return httpServer
}
