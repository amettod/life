package game

import "math/rand"

type Game interface {
	// Clear state.
	Clear()
	// Random fills no more than a quarter of the state.
	Random()
	// Resize state
	Resize(w, h int)
	// SetState to the origin x y.
	SetState(x, y int, s [][]int)
	// Shift cell state.
	Shift(x, y int)
	// State return.
	State() [][]int
	// Step to the next state.
	Step()
}

type game struct {
	s state
}

func New(w, h int) Game {
	return &game{
		s: newState(w, h),
	}
}

func (g *game) Clear() {
	g.s = newState(g.s.width(), g.s.height())
}

func (g *game) Random() {
	w := g.s.width()
	h := g.s.height()
	s := newState(w, h)
	for i := 0; i < w*h/4; i++ {
		s.init(rand.Intn(w), rand.Intn(h), true)
	}
	g.s = s
}

func (g *game) Resize(w, h int) {
	s := newState(w, h)
	for y := range g.s {
		for x, count := range g.s[y] {
			s.setCycle(x, y, count)
		}
	}
	g.s = s
}

func (g *game) SetState(x, y int, s [][]int) {
	for yy := range s {
		for xx := range s[yy] {
			g.s.init(x+xx, y+yy, state(s).alive(xx, yy))
		}
	}
}

func (g *game) Shift(x, y int) {
	g.s.cycleCalc(x, y, !g.s.alive(x, y))
}

func (g *game) State() [][]int {
	return g.s
}

func (g *game) Step() {
	s := newState(g.s.width(), g.s.height())
	for y := range g.s {
		for x := range g.s[y] {
			s.setCycle(x, y, g.s.cycle(x, y))
			s.cycleCalc(x, y, g.s.next(x, y))
		}
	}
	g.s = s
}
