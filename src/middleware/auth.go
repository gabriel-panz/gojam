package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gabriel-panz/gojam/handler"
	"github.com/gabriel-panz/gojam/types"
	"github.com/gabriel-panz/gojam/utils"
)

func Authorize(logger *log.Logger, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				hf := handler.RedirectToSpotifyAuth(logger)
				hf(w, r)
				return
			}

			utils.HandleHttpError(err, logger, w)
			return
		}

		refresh_token, err := r.Cookie("refresh_token")

		if err != nil {
			if err == http.ErrNoCookie {
				hf := handler.RedirectToSpotifyAuth(logger)
				hf(w, r)
				return
			}

			utils.HandleHttpError(err, logger, w)
			return
		}

		ctxAuth := context.WithValue(r.Context(), types.AuthorizedKey, types.Auth{
			Token:        token.Value,
			RefreshToken: refresh_token.Value,
		})

		rAuthorized := r.WithContext(ctxAuth)

		h(w, rAuthorized)
	})
}
