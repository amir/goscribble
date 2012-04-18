package mpdclient

import (
	"strconv"
	"strings"
)

type PlayState uint8

const (
	Playing PlayState = iota
	Paused
	Stopped
)

type Status struct {
	State       PlayState
	TotalTime   int
	ElapsedTime int
}

func readStatus(values Values) *Status {
	status := new(Status)

	switch values["state"] {
	case "play":
		status.State = Playing
	case "pause":
		status.State = Paused
	case "stop":
		status.State = Stopped
	}

	if status.State != Stopped {
		s := values["time"]
		z := strings.Index(s, ":")
		elapsedTime, err := strconv.Atoi(s[0:z])
		if err != nil {
			elapsedTime = 0
		}
		totalTime, err := strconv.Atoi(s[z+1:])
		if err != nil {
			totalTime = 0
		}
		status.TotalTime = totalTime
		status.ElapsedTime = elapsedTime
	}

	return status
}
