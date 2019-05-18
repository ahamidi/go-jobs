package jobs

type state int

const (
	New state = iota
	Running
	Paused
	Stopped
)
