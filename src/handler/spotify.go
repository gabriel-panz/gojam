package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gabriel-panz/gojam/spotify"
	"github.com/gabriel-panz/gojam/types"
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

func ListDevices(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := utils.GetAuthorizedUser(r)

		ds, err := spotify.GetDevices(u.Token)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		devices := components.Devices(ds.Devices)
		devices.Render(r.Context(), w)
	})
}

func Play(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dId := r.URL.Query().Get("deviceId")

		if dId == "" {
			utils.HandleHttpError(errors.New("no device was selected"), logger, w)
		}

		token := utils.GetAuthorizedUser(r).Token
		err := spotify.Play(token, dId)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		player := components.Player(dId, types.PauseState)
		player.Render(r.Context(), w)
	})
}

func Pause(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dId := r.URL.Query().Get("deviceId")

		if dId == "" {
			utils.HandleHttpError(errors.New("no device was selected"), logger, w)
		}

		token := utils.GetAuthorizedUser(r).Token
		err := spotify.Pause(token, dId)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		player := components.Player(dId, types.PlayState)
		player.Render(r.Context(), w)
	})
}

func Player(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dId := r.URL.Query().Get("deviceId")
		pState := r.URL.Query().Get("state")

		if pState == "" {
			utils.HandleHttpError(errors.New("undefined player state"), logger, w)
		}

		pStateInt, err := strconv.ParseInt(pState, 10, 32)
		if err != nil {
			utils.HandleHttpError(errors.New("player state must be integer"), logger, w)
		}

		player := components.Player(dId, types.PlayerState(pStateInt))
		player.Render(r.Context(), w)
	})
}
