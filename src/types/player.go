package types

type PlayerState int

const (
	PlayState PlayerState = iota
	PauseState
	NoDevice
)

type ShowType string

const (
	Playlist ShowType = "p"
	Track    ShowType = "t"
)
