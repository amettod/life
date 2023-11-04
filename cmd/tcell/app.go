package main

import (
	"fmt"
	"os"
	"time"

	"github.com/amettod/life"
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

type app struct {
	*life.App

	screen tcell.Screen

	period time.Duration
	rate   int
}

func newApp(file string, d time.Duration, rate int) (*app, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	w, h := s.Size()
	a, err := life.NewApp(w/rate, h, file)
	if err != nil {
		return nil, err
	}

	sd := tcell.StyleDefault.
		Background(rgbTo(a.Theme.Background())).
		Foreground(rgbTo(a.Theme.Foreground()))
	s.SetStyle(sd)
	s.EnableMouse()

	return &app{
		App:    a,
		screen: s,
		period: d,
		rate:   rate,
	}, nil
}

func (a *app) setInfo(x, y int, msg string) {
	sd := tcell.StyleDefault.
		Background(rgbTo(a.Theme.Background())).
		Foreground(rgbTo(a.Theme.Foreground()))
	for i, r := range msg {
		a.screen.SetContent(x+i, y, r, nil, sd)
	}
}

func (a *app) waitEvent(e chan<- event, p chan<- eventPoint) {
	for {
		switch ev := a.screen.PollEvent().(type) {
		case *tcell.EventResize:
			a.screen.Sync()
			e <- eventResize
		case *tcell.EventMouse:
			switch ev.Buttons() {
			case tcell.Button1:
				p <- eventShift.point(ev.Position())
			case tcell.Button2:
				p <- eventInsert.point(ev.Position())
			default:
				continue
			}
		case *tcell.EventKey:
			switch {
			case ev.Key() == tcell.KeyEnter:
				e <- eventStep
			case ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q':
				e <- eventQuit
			case ev.Rune() == ' ':
				e <- eventPause
			case ev.Rune() == 'c':
				e <- eventClear
			case ev.Rune() == 'p':
				e <- eventPreset
			case ev.Rune() == 'r':
				e <- eventRandom
			case ev.Rune() == 't':
				e <- eventTheme
			case ev.Rune() == 'h':
				e <- eventInfo
			}
		default:
			continue
		}
	}
}

func (a *app) doEvent(e <-chan event, p <-chan eventPoint) {
	ticker := time.NewTicker(a.period * time.Millisecond)
	cycle := 0
	stop := true
	info := true
	theme := a.Theme
	for {
		a.Theme = theme
		a.screen.Clear()
		for y, row := range a.Game.State() {
			for x, cycle := range row {
				a.screen.SetContent(x*a.rate, y, ' ', nil, tcell.StyleDefault.
					Background(rgbTo(a.Theme.Color(cycle))))
				a.screen.SetContent(x*a.rate+1, y, ' ', nil, tcell.StyleDefault.
					Background(rgbTo(a.Theme.Color(cycle))))
			}
		}
		select {
		case ev := <-e:
			switch ev {
			case eventRandom:
				a.Game.Random()
				cycle = 0
			case eventPause:
				stop = !stop
			case eventResize:
				w, h := a.screen.Size()
				a.Game.Resize(w/a.rate, h)
			case eventStep:
				cycle++
				a.Game.Step()
				a.screen.Show()
			case eventQuit:
				ticker.Stop()
				a.screen.Fini()
				os.Exit(0)
			case eventClear:
				a.Game.Clear()
				cycle = 0
				a.screen.Show()
			case eventInfo:
				info = !info
			case eventTheme:
				theme.Next()
			case eventPreset:
				a.Preset.Next()
			}
		case ep := <-p:
			switch ep.e {
			case eventShift:
				a.Game.Shift(ep.x/a.rate, ep.y)
			case eventInsert:
				a.Game.SetState(ep.x/a.rate, ep.y, a.Preset.State())
			}
		case <-ticker.C:
			if !stop {
				cycle++
				a.Game.Step()
			}
			if stop && info {
				_, h := a.screen.Size()
				a.setInfo(0, 0, fmt.Sprintf("Cycle: %d", cycle))
				a.setInfo(0, h-4, fmt.Sprintf("t: switch theme, Current: \"%s\"", a.Theme.Name()))
				a.setInfo(0, h-3, fmt.Sprintf("p: switch present, Current: \"%s\"", a.Preset.Name()))
				a.setInfo(0, h-2, "LeftClick: toggle state, RightClick: insert preset")
				a.setInfo(0, h-1, "SPC: pause, Enter: next, c: clear, r: random, h: hide this message")
			}
			a.screen.Show()
		}
	}
}

func rgbTo(rgb life.RGB) tcell.Color {
	r, g, b := rgb.Color()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
