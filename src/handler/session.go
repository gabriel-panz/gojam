package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/session"
	"github.com/gabriel-panz/gojam/ui/components"
	"github.com/gabriel-panz/gojam/ui/pages"
	"github.com/gabriel-panz/gojam/utils"
)

func Session(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sId := r.PathValue("session_id")
		ses, err := service.GetSession(sId)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sId,
		})

		sPage := pages.Session(ses)
		sPage.Render(r.Context(), w)
	})
}

// Creates and caches an empty session
func CreateSession(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := utils.GetAuthorizedUser(r)
		ses, err := service.CreateSession(a.Token)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		path := fmt.Sprintf("/session/%s", ses.ID)

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: ses.ID,
		})

		w.Header().Add("HX-Redirect", path)
	})
}

func LeaveSession(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := utils.GetAuthorizedUser(r)
		sId, err := r.Cookie("session_id")
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		err = service.LeaveSession(sId.Value, a.Token)

		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			MaxAge: 0,
		})

		sesCreate := components.SessionCreate()
		sesCreate.Render(r.Context(), w)
	})
}

func DeleteSession(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := utils.GetAuthorizedUser(r)
		sId, err := r.Cookie("session_id")
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		service.DeleteSession(sId.Value, a.Token)

		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			MaxAge: 0,
		})
	})
}
