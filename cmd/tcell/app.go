package main

import (
	"fmt"
	"os"
	"time"

	"github.com/amettod/life/internal/game"
	"github.com/amettod/life/internal/parse"
	"github.com/amettod/life/internal/preset"
	"github.com/amettod/life/internal/theme"
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

type app struct {
	game   game.Game
	preset preset.Preset
	theme  theme.Theme

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
	t := theme.New()
	sd := tcell.StyleDefault.
		Background(toTcellColor(t.Background())).
		Foreground(toTcellColor(t.Foreground()))
	s.SetStyle(sd)
	s.EnableMouse()

	w, h := s.Size()
	g := game.New(w/rate, h)
	if file != "" {
		s, err := parse.File(file)
		if err != nil {
			return nil, err
		}
		g.SetState(0, 0, s)
	}

	p := preset.New()
	states, err := parse.FilesEmbed(preset.EmbedFS, preset.EmbedDir)
	if err != nil {
		return nil, err
	}
	for name, state := range states {
		p.Append(name, state)
	}

	return &app{
		game:   g,
		preset: p,
		theme:  t,

		screen: s,

		period: d,
		rate:   rate,
	}, nil
}

func (a *app) info(x, y int, msg string) {
	sd := tcell.StyleDefault.
		Background(toTcellColor(a.theme.Background())).
		Foreground(toTcellColor(a.theme.Foreground()))
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
	theme := a.theme
	for {
		a.theme = theme
		a.screen.Clear()
		for y, row := range a.game.State() {
			for x, cycle := range row {
				a.screen.SetContent(x*a.rate, y, ' ', nil, tcell.StyleDefault.
					Background(toTcellColor(a.theme.Color(cycle))))
				a.screen.SetContent(x*a.rate+1, y, ' ', nil, tcell.StyleDefault.
					Background(toTcellColor(a.theme.Color(cycle))))
			}
		}
		select {
		case ev := <-e:
			switch ev {
			case eventRandom:
				a.game.Random()
				cycle = 0
			case eventPause:
				stop = !stop
			case eventResize:
				w, h := a.screen.Size()
				a.game.Resize(w/a.rate, h)
			case eventStep:
				cycle++
				a.game.Step()
				a.screen.Show()
			case eventQuit:
				ticker.Stop()
				a.screen.Fini()
				os.Exit(0)
			case eventClear:
				a.game.Clear()
				cycle = 0
				a.screen.Show()
			case eventInfo:
				info = !info
			case eventTheme:
				theme.Next()
			case eventPreset:
				a.preset.Next()
			}
		case ep := <-p:
			switch ep.e {
			case eventShift:
				a.game.Shift(ep.x/a.rate, ep.y)
			case eventInsert:
				a.game.SetState(ep.x/a.rate, ep.y, a.preset.State())
			}
		case <-ticker.C:
			if !stop {
				cycle++
				a.game.Step()
			}
			if stop && info {
				_, h := a.screen.Size()
				a.info(0, 0, fmt.Sprintf("Cycle: %d", cycle))
				a.info(0, h-4, fmt.Sprintf("t: switch theme, Current: \"%s\"", a.theme.Name()))
				a.info(0, h-3, fmt.Sprintf("p: switch present, Current: \"%s\"", a.preset.Name()))
				a.info(0, h-2, "LeftClick: toggle state, RightClick: insert preset")
				a.info(0, h-1, "SPC: pause, Enter: next, c: clear, r: random, h: hide this message")
			}
			a.screen.Show()
		}
	}
}

func toTcellColor(rgb theme.RGB) tcell.Color {
	r, g, b := rgb.Color()
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
