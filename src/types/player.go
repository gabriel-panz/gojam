package types

type PlayerState int

const (
	PlayState PlayerState = iota
	PauseState
	NoDevice
)
