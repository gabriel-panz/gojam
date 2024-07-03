package handler

import (
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/spotify"
	"github.com/gabriel-panz/gojam/ui/components"
	"github.com/gabriel-panz/gojam/utils"
)

func ListPlaylists(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := utils.GetPagination(r)

		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}

		// Get user playlists
		ps, err := spotify.GetPlaylists(p, utils.GetAuthorizedUser(r).Token)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}

		ul := components.PlaylistList(ps, p.PageIndex)
		ul.Render(r.Context(), w)
	})
}

func Play(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetAuthorizedUser(r).Token
		err := spotify.Play(token)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}
	})
}
