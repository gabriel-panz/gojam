package spotify

type Profile struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Id          string `json:"id"`
	Uri         string `json:"uri"`
	Href        string `json:"href"`
	Product     string `json:"product"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
}

type Playlist struct {
	Collaborative bool         `json:"collaborative"`
	Description   string       `json:"description"`
	ExternalUrls  ExternalUrls `json:"external_urls"`
	Href          string       `json:"href"`
	ID            string       `json:"id"`
	Images        []Images     `json:"images"`
	Name          string       `json:"name"`
	Owner         Owner        `json:"owner"`
	Public        bool         `json:"public"`
	SnapshotID    string       `json:"snapshot_id"`
	Tracks        Tracks       `json:"tracks"`
	Type          string       `json:"type"`
	URI           string       `json:"uri"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Images struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
type Followers struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}
type Owner struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Followers    Followers    `json:"followers"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
	DisplayName  string       `json:"display_name"`
}

type Tracks struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsPlayable bool   `json:"is_playable"`
	Href       string `json:"href"`
	URI        string `json:"uri"`
	Total      int    `json:"total"`
}

type DeviceResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
	SupportsVolume   bool   `json:"supports_volume"`
}

type RedirectParams struct {
	ClientId            string `url:"client_id"`
	ResponseType        string `url:"response_type"`
	RedirectUri         string `url:"redirect_uri"`
	Scope               string `url:"scope"`
	CodeChallengeMethod string `url:"code_challenge_method"`
	CodeChallenge       string `url:"code_challenge"`
}

type AccessTokenRequest struct {
	ClientId     string `url:"client_id"`
	GrantType    string `url:"grant_type"`
	Code         string `url:"code"`
	RedirectUri  string `url:"redirect_uri"`
	CodeVerifier string `url:"code_verifier"`
}

type RefreshTokenRequest struct {
	ClientId     string `url:"client_id"`
	GrantType    string `url:"grant_type"`
	RefreshToken string `url:"redirect_uri"`
}

type PlayRequest struct {
	// Optional. Spotify URI of the context to play. Valid contexts are albums, artists & playlists.
	ContextURI string `json:"context_uri,omitempty"`
	// Optional. A JSON array of the Spotify track URIs to play. For example: {"uris": ["spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M"]}
	Uris []string `json:"uris,omitempty"`
	// Optional. Indicates from where in the context playback should start. Only available when context_uri corresponds to an album or playlist object "position" is zero based and can’t be negative.
	Offset     *Offset `json:"offset,omitempty"`
	PositionMs int     `json:"position_ms,omitempty"`
}

type Offset struct {
	Position int `json:"position"`
}
