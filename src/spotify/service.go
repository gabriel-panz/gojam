package spotify

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/types"
	"github.com/gabriel-panz/gojam/utils"

	"github.com/google/go-querystring/query"
)

const (
	AccountsUrl = "https://accounts.spotify.com/"
	ApiUrl      = "https://api.spotify.com/v1/"
	playlists   = ApiUrl + "me/playlists"
	devices     = ApiUrl + "me/player/devices"
	play        = ApiUrl + "me/player/play"
	tokenUrl    = AccountsUrl + "api/token"
)

type Service struct {
	Client *http.Client
}

func GetAuthConsentUrl(codeChallange string) (string, error) {
	params := RedirectParams{
		ClientId:            os.Getenv("CLIENT_ID"),
		ResponseType:        "code",
		RedirectUri:         "http://192.168.0.15:9000/callback",
		Scope:               "user-read-private user-read-email playlist-read-private user-read-playback-state user-modify-playback-state",
		CodeChallengeMethod: "S256",
		CodeChallenge:       codeChallange,
	}

	v, err := query.Values(params)
	if err != nil {
		return "", err
	}

	return AccountsUrl + "authorize?" + v.Encode(), nil
}

func (s Service) GetProfileData(token string) (*Profile, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		res, err := utils.DecodeJsonBody[UnauthorizedRequest](resp)

		if err != nil {
			return nil, err
		}

		if res.Err.Message == "The access token expired" {
			return nil, errors.ErrExpiredToken
		}

		return nil, errors.ErrUnauthorized
	}

	p, err := utils.DecodeJsonBody[Profile](resp)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s Service) GetAccessToken(verifier string, code string) (*Token, error) {
	body := AccessTokenRequest{
		ClientId:     os.Getenv("CLIENT_ID"),
		GrantType:    "authorization_code",
		Code:         code,
		RedirectUri:  "http://192.168.0.15:9000/callback",
		CodeVerifier: verifier,
	}

	v, err := utils.EncodeQueryParameters(body)
	if err != nil {
		return nil, err
	}

	b := strings.NewReader(v)
	req, err := http.NewRequest(http.MethodPost, tokenUrl, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req)

	if err != nil {
		return nil, err
	}

	token, err := utils.DecodeJsonBody[Token](resp)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s Service) RefreshAccessToken(refresh string) (*Token, error) {
	body := RefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refresh,
		ClientId:     os.Getenv("CLIENT_ID"),
	}

	v, err := query.Values(body)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(v.Encode())

	req, err := http.NewRequest(http.MethodPost, tokenUrl, sr)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.ErrUnauthorized
	}

	token, err := utils.DecodeJsonBody[Token](resp)
	if err != nil {
		return nil, err
	}
	return token, nil
}

type OffsetQuery struct {
	Limit  int `url:"limit"`
	Offset int `url:"offset"`
}

type OffsetResponse[T any] struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []T    `json:"items"`
}

type BadRequest struct {
	Err              string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type UnauthorizedRequest struct {
	Err Error `json:"error"`
}
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (s Service) GetPlaylists(p *types.Pagination, token string) ([]Playlist, error) {
	q := OffsetQuery{
		Limit:  p.PageSize,
		Offset: p.GetOffset(),
	}

	v, err := query.Values(q)
	if err != nil {
		return nil, err
	}

	enc := v.Encode()
	req, err := http.NewRequest(http.MethodGet, playlists+"?"+enc, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		r, err := utils.DecodeJsonBody[BadRequest](resp)
		if err != nil {
			return nil, err
		}

		fmt.Println(r)
		return nil, errors.ErrBadRequest
	} else {
		r, err := utils.DecodeJsonBody[OffsetResponse[Playlist]](resp)
		if err != nil {
			return nil, err
		}
		return r.Items, nil
	}
}

func (s Service) Play(token string, deviceId string) error {
	req, err := http.NewRequest(http.MethodPut, play+"?device_id="+deviceId, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	} else {
		return errors.ErrBadRequest
	}
}

func (s Service) GetDevices(token string) (*DeviceResponse, error) {
	req, err := http.NewRequest(http.MethodGet, devices, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	ds, err := utils.DecodeJsonBody[DeviceResponse](resp)
	if err != nil {
		return nil, err
	}

	return ds, nil
}
