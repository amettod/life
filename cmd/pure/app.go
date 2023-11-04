package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amettod/life"
	"github.com/amettod/life/term"
)

const (
	unitCell = "  "
	unitHide = "@"
)

type app struct {
	*life.App

	term term.Term

	period time.Duration
	info   []string
}

func newApp(w, h int, file string, d time.Duration) (*app, error) {
	a, err := life.NewApp(w, h, file)
	if err != nil {
		return nil, err
	}

	return &app{
		App:    a,
		term:   term.New(os.Stdout),
		period: d,
		info:   make([]string, h),
	}, nil
}

func (a *app) setInfo(x, y int, msg string) {
	if y < a.Game.Height() {
		s := fmt.Sprint(strings.Repeat(unitHide, x), msg)
		if w := a.Game.Width() * len(unitCell); len(s) > w {
			s = fmt.Sprint(s[:w-3], "...")
		}
		a.info[y] = s
	}
}

func (a *app) show() {
	for y := range a.Game.State() {
		x := -1
		for _, cycle := range a.Game.State()[y] {
			for i := 0; i < len(unitCell); i++ {
				x++
				if x < len(a.info[y]) && string(a.info[y][x]) != unitHide {
					a.term.Write(
						a.Theme.Background(),
						a.Theme.Foreground(),
						string(a.info[y][x]),
					)
					continue
				}
				a.term.Write(
					a.Theme.Color(cycle),
					nil,
					string(unitCell[i]),
				)
			}
		}
		a.info[y] = ""
		a.term.WriteLn()
	}
	a.term.Print()
}

func (a *app) waitEvent(e chan<- event) {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		line := scan.Text()
		switch {
		case strings.Contains(line, "q"):
			e <- eventQuit
		case strings.Contains(line, " "):
			e <- eventPause
		case strings.Contains(line, "t"):
			e <- eventTheme
		case strings.Contains(line, "r"):
			e <- eventRandom
		case strings.Contains(line, "c"):
			e <- eventClear
		case strings.Contains(line, "s"):
			e <- eventStep
		case strings.Contains(line, "p"):
			e <- eventSwitchPreset
		case strings.Contains(line, "i"):
			e <- eventInsertPreset
		case strings.Contains(line, "h"):
			e <- eventInfo
		}
	}
}

func (a *app) doEvent(e <-chan event) {
	ticker := time.NewTicker(a.period * time.Millisecond)
	stop := true
	info := true
	cycle := 0
	for {
		select {
		case ev := <-e:
			switch ev {
			case eventQuit:
				os.Exit(0)
			case eventPause:
				stop = !stop
			case eventTheme:
				a.Theme.Next()
			case eventRandom:
				a.Game.Random()
				cycle = 0
			case eventClear:
				a.Game.Clear()
				cycle = 0
			case eventStep:
				a.Game.Step()
				cycle++
			case eventSwitchPreset:
				a.Preset.Next()
			case eventInsertPreset:
				a.Game.Clear()
				a.Game.SetState(0, 0, a.Preset.State())
				cycle = 0
			case eventInfo:
				info = !info
			}
		case <-ticker.C:
			if stop && info {
				h := a.Game.Height()
				a.setInfo(0, 0, fmt.Sprintf("Cycle: %d", cycle))
				a.setInfo(0, h-4, "Press <key>+RET:")
				a.setInfo(0, h-3, fmt.Sprintf("<t>: switch theme, Current: \"%s\"", a.Theme.Name()))
				a.setInfo(0, h-2, fmt.Sprintf("<p>: switch present, <i>: insert preset, Current: \"%s\"", a.Preset.Name()))
				a.setInfo(0, h-1, "<SPC>: pause, <s>: next, <c>: clear, <r>: random, <h>: hide this message")
			}

			if !stop {
				a.Game.Step()
				cycle++
			}
			a.show()
		}
	}
}
