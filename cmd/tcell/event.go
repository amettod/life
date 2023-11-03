package main

type event uint

const (
	eventQuit event = iota
	eventPause
	eventRandom
	eventStep
	eventClear
	eventResize
	eventInfo
	eventPreset
	eventTheme
	eventShift
	eventInsert
)

type eventPoint struct {
	e    event
	x, y int
}

func (e event) point(x, y int) eventPoint {
	return eventPoint{
		e: e,
		x: x,
		y: y,
	}
}
