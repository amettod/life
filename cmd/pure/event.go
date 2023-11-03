package main

type event uint

const (
	eventQuit event = iota
	eventPause
	eventTheme
	eventRandom
	eventClear
	eventStep
	eventSwitchPreset
	eventInsertPreset
	eventInfo
)
