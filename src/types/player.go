package types

type PlayerState int

const (
	PlayState PlayerState = iota
	PauseState
	NoDevice
)

type ShowType string

const (
	None     ShowType = ""
	Playlist ShowType = "p"
	Track    ShowType = "t"
)
