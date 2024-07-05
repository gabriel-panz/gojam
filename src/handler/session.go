package handler

import (
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/session"
	"github.com/gabriel-panz/gojam/utils"
)

func CreateSession(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := utils.GetAuthorizedUser(r)
		ses, err := service.CreateSession(a.Token)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: ses.ID,
		})
	})
}

func JoinSession(logger *log.Logger, service session.SessionService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := utils.GetAuthorizedUser(r)
		sId := r.PathValue("session_id")
		err := service.JoinSession(sId, a.Token)

		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sId,
		})
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
