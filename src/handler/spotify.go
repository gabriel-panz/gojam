package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	gojamErrors "github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/session"
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

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		ul := components.PlaylistList(ps, p.PageIndex, sId)
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

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		devices := components.Devices(ds.Devices, sId)
		devices.Render(r.Context(), w)
	})
}

func Play(logger *log.Logger, spotify spotify.Service, sess session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dId := r.FormValue("deviceId")
		uri := r.FormValue("uri")
		t := r.FormValue("type")

		if dId == "" {
			utils.HandleHttpError(errors.New("no device was selected"), logger, w)
			return
		}

		if uri != "" {
			if t == "" {
				utils.HandleHttpError(errors.New("no type was selected"), logger, w)
				return
			}
		}

		token := utils.GetAuthorizedUser(r).Token

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		s, err := sess.GetSession(sId)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}
		if token != s.DjToken {
			return
		}

		// do this for every client in session
		for cToken, client := range s.Conns {
			m := session.SessionMessage{
				DeviceID: client.DeviceId,
				Uri:      uri,
				Type:     types.ShowType(t),
			}

			err := spotify.Play(cToken, client.DeviceId, uri, types.ShowType(t))
			if err != nil {
				utils.HandleHttpError(err, logger, w)
			}

			go session.Broadcast(m, client.Conn)
		}

		// broadcast this
		player := components.Player(dId, types.PauseState, sId)
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

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		player := components.Player(dId, types.PlayState, sId)
		player.Render(r.Context(), w)
	})
}

func Player(logger *log.Logger, spotify spotify.Service, sesService session.SessionService) http.HandlerFunc {
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

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		s, err := sesService.GetSession(sId)
		if err != nil {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		ses := s.Conns[utils.GetAuthorizedUser(r).Token]

		ses.DeviceId = dId

		player := components.Player(dId, types.PlayerState(pStateInt), sId)
		player.Render(r.Context(), w)
	})
}

func Search(logger *log.Logger, s spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		p, err := utils.GetPagination(r)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		if q == "" {
			utils.HandleHttpError(errors.New("please inform a search query"), logger, w)
		}

		token := utils.GetAuthorizedUser(r).Token

		params := spotify.SearchParams{
			Query: q,
			Type: []string{
				"playlist", "track",
			},
			Limit:  p.PageSize,
			Offset: p.PageIndex,
		}

		res, err := s.Search(token, params)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		sId := r.FormValue("session_id")
		if sId == "" {
			utils.HandleHttpError(gojamErrors.ErrNotFound, logger, w)
			return
		}

		results := components.SearchResults(*res, p.PageIndex, sId)
		results.Render(r.Context(), w)
	})
}
