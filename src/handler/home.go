package handler

import (
	"log"
	"time"

	"github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/spotify"
	"github.com/gabriel-panz/gojam/ui/components"
	"github.com/gabriel-panz/gojam/ui/pages"
	"github.com/gabriel-panz/gojam/utils"

	"net/http"
)

type HomeHandler struct {
	Spotify *spotify.Service
	Handler
}

const (
	cookieTokenName    = "token"
	cookieRefreshName  = "refresh_token"
	cookieVerifierName = "code_verifier"
	possible           = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	possibleLen        = len(possible)
)

func RedirectToSpotifyAuth(logger *log.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifier, err := utils.GenerateCodeVerifier(128)
		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}

		// saves an auth session
		http.SetCookie(w, &http.Cookie{
			Name:     cookieVerifierName,
			Value:    verifier,
			MaxAge:   int(2 * time.Minute),
			HttpOnly: true,
		})

		challenge := utils.GenerateCodeChallenge(verifier)

		url, err := spotify.GetAuthConsentUrl(challenge)

		if err != nil {
			utils.HandleHttpError(err, logger, w)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	})
}

func (h HomeHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.HandleHttpError(err, h.Logger, w)
		return
	}
	code := r.FormValue("code")

	verifier, err := r.Cookie(cookieVerifierName)
	if err != nil {
		if err == http.ErrNoCookie {
			utils.HandleHttpError(errors.ErrUnauthorized, h.Logger, w)
			return
		}
		utils.HandleHttpError(err, h.Logger, w)
		return
	}

	tokenResp, err := h.Spotify.GetAccessToken(verifier.Value, code)

	if err != nil {
		utils.HandleHttpError(err, h.Logger, w)
		return
	}

	SetAuthCookies(w, tokenResp)

	http.Redirect(w, r, "http://192.168.0.15:9000/", http.StatusSeeOther)
}

func (h HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	token := utils.GetAuthorizedUser(r).Token

	prof, err := h.Spotify.GetProfileData(token)

	if err != nil {
		switch err {
		case errors.ErrUnauthorized:
			hf := RedirectToSpotifyAuth(h.Logger)
			hf(w, r)
		case errors.ErrExpiredToken:
			_, err := refreshToken(*h.Spotify, w, r)
			if err != nil {
				if err == errors.ErrUnauthorized {
					hf := RedirectToSpotifyAuth(h.Logger)
					hf(w, r)
				}

				utils.HandleHttpError(err, h.Logger, w)
			} else {
				// tries again
				h.GetHomePage(w, r)
			}
		default:
			utils.HandleHttpError(err, h.Logger, w)
		}
		return
	}

	home := pages.Home(*prof)
	home.Render(r.Context(), w)
}

func Refresh(logger *log.Logger, spotify spotify.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := refreshToken(spotify, w, r)

		if err != nil {
			utils.HandleHttpError(err, logger, w)
		}

		authRefresh := components.AuthRefresh(token.ExpiresInSeconds)
		authRefresh.Render(r.Context(), w)
	})
}

func refreshToken(
	spotify spotify.Service,
	w http.ResponseWriter,
	r *http.Request,
) (*spotify.Token, error) {
	refresh := utils.GetAuthorizedUser(r).RefreshToken

	token, err := spotify.RefreshAccessToken(refresh)
	if err != nil {
		return nil, err
	}

	SetAuthCookies(w, token)

	return token, nil
}

func SetAuthCookies(w http.ResponseWriter, token *spotify.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieTokenName,
		Value:  token.AccessToken,
		MaxAge: int(token.ExpiresInSeconds * int(time.Second)),
	})

	http.SetCookie(w, &http.Cookie{
		Name:   cookieRefreshName,
		Value:  token.RefreshToken,
		MaxAge: int(token.ExpiresInSeconds * int(time.Second)),
	})
}
