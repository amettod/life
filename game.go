package life

import "math/rand"

type game struct {
	s state
}

func newGame(w, h int) *game {
	return &game{
		s: newState(w, h),
	}
}

// Clear state.
func (g *game) Clear() {
	g.s = newState(g.s.width(), g.s.height())
}

// Random fills no more than a quarter of the state.
func (g *game) Random() {
	w := g.s.width()
	h := g.s.height()
	s := newState(w, h)
	for i := 0; i < w*h/4; i++ {
		s.init(rand.Intn(w), rand.Intn(h), true)
	}
	g.s = s
}

// Resize state
func (g *game) Resize(w, h int) {
	s := newState(w, h)
	for y := range g.s {
		for x, count := range g.s[y] {
			s.setCycle(x, y, count)
		}
	}
	g.s = s
}

// SetState to the origin x y.
func (g *game) SetState(x, y int, s [][]int) {
	for yy := range s {
		for xx := range s[yy] {
			g.s.init(x+xx, y+yy, state(s).alive(xx, yy))
		}
	}
}

// Shift cell state.
func (g *game) Shift(x, y int) {
	g.s.cycleCalc(x, y, !g.s.alive(x, y))
}

// State return.
func (g *game) State() [][]int {
	return g.s
}

// Step to the next state.
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

// Height state return.
func (g *game) Height() int {
	return g.s.height()
}

// Width state return.
func (g *game) Width() int {
	return g.s.width()
}
