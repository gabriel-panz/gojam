package api

import (
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/handler"
	"github.com/gabriel-panz/gojam/spotify"
)

func NewServer(addr string, clientId string) *http.Server {
	mux := http.NewServeMux()

	spotifyService := spotify.Service{
		Client: &http.Client{},
	}

	homeHandler := handler.HomeHandler{
		Spotify: &spotifyService,
		Handler: handler.Handler{
			Logger: log.Default(),
		},
	}

	addRoutes(mux, homeHandler, spotifyService)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return httpServer
}
