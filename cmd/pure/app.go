package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amettod/life/game"
	"github.com/amettod/life/parse"
	"github.com/amettod/life/preset"
	"github.com/amettod/life/term"
	"github.com/amettod/life/theme"
)

const (
	unitCell = "  "
	unitHide = "@"
)

type app struct {
	game   game.Game
	preset preset.Preset
	term   term.Term
	theme  theme.Theme

	period time.Duration

	info []string
}

func newApp(w, h int, file string, p time.Duration) (*app, error) {
	g := game.New(w, h)
	if file != "" {
		s, err := parse.File(file)
		if err != nil {
			return nil, err
		}
		g.SetState(0, 0, s)
	}

	pr := preset.New()
	fse, err := parse.FilesEmbed(preset.EmbedFS, preset.EmbedDir)
	if err != nil {
		return nil, err
	}
	for k, v := range fse {
		pr.Append(k, v)
	}

	return &app{
		game:   g,
		preset: pr,
		term:   term.New(os.Stdout),
		theme:  theme.New(),
		period: p,
		info:   make([]string, h),
	}, err
}

func (a *app) height() int {
	return len(a.game.State())
}

func (a *app) width() int {
	if a.height() > 0 {
		return len(a.game.State()[0])
	}
	return 0
}

func (a *app) addInfo(x, y int, msg string) {
	if y < a.height() {
		s := fmt.Sprint(strings.Repeat(unitHide, x), msg)
		if w := a.width() * len(unitCell); len(s) > w {
			s = fmt.Sprint(s[:w-3], "...")
		}
		a.info[y] = s
	}
}

func (a *app) show() {
	for y := range a.game.State() {
		x := -1
		for _, cycle := range a.game.State()[y] {
			for _, cell := range unitCell {
				x++
				if x < len(a.info[y]) && string(a.info[y][x]) != unitHide {
					a.term.SetInfo(
						a.theme.Background(),
						a.theme.Foreground(),
						a.info[y][x],
					)
					continue
				}
				a.term.SetCell(
					a.theme.Color(cycle),
					cell,
				)
			}
		}
		a.info[y] = ""
		a.term.SetLn()
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
				a.theme.Next()
			case eventRandom:
				a.game.Random()
				cycle = 0
			case eventClear:
				a.game.Clear()
				cycle = 0
			case eventStep:
				a.game.Step()
				cycle++
			case eventSwitchPreset:
				a.preset.Next()
			case eventInsertPreset:
				a.game.Clear()
				a.game.SetState(0, 0, a.preset.State())
				cycle = 0
			case eventInfo:
				info = !info
			}
		case <-ticker.C:
			if stop && info {
				h := a.height()
				a.addInfo(0, 0, fmt.Sprintf("Cycle: %d", cycle))
				a.addInfo(0, h-4, "Press <key>+RET:")
				a.addInfo(0, h-3, fmt.Sprintf("<t>: switch theme, Current: \"%s\"", a.theme.Name()))
				a.addInfo(0, h-2, fmt.Sprintf("<p>: switch present, <i>: insert preset, Current: \"%s\"", a.preset.Name()))
				a.addInfo(0, h-1, "<SPC>: pause, <s>: next, <c>: clear, <r>: random, <h>: hide this message")
			}

			if !stop {
				a.game.Step()
				cycle++
			}
			a.show()
		}
	}
}
