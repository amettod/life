package life

import "math"

type state [][]int

func newState(w, h int) state {
	s := make([][]int, h)
	for i := range s {
		s[i] = make([]int, w)
	}
	return s
}

func (s state) height() int {
	return len(s)
}

func (s state) width() int {
	if s.height() > 0 {
		return len(s[0])
	}
	return 0
}

func (s state) inside(x, y int) bool {
	return x >= 0 && y >= 0 && x < s.width() && y < s.height()
}

func (s state) cycle(x, y int) int {
	if s.inside(x, y) {
		return s[y][x]
	}
	return 0
}

func (s state) setCycle(x, y int, c int) {
	if s.inside(x, y) {
		s[y][x] = c
	}
}

func (s state) init(x, y int, alive bool) {
	if alive {
		s.setCycle(x, y, 1)
	}
}

func (s state) alive(x, y int) bool {
	return s.cycle(x, y) > 0
}

func (s state) boundless(x, y int) (int, int) {
	if !s.inside(x, y) {
		x += s.width()
		x %= s.width()
		y += s.height()
		y %= s.height()
	}
	return x, y
}

func (s state) cycleCalc(x, y int, alive bool) {
	x, y = s.boundless(x, y)
	switch {
	case alive && s[y][x] > 0 && s[y][x] != math.MaxInt:
		s[y][x]++
	case alive:
		s[y][x] = 1
	case !alive && s[y][x] < 0 && s[y][x] != math.MinInt:
		s[y][x]--
	case !alive && s[y][x] > 0:
		s[y][x] = -1
	default:
		s[y][x] = 0
	}
}

func (s state) next(x, y int) bool {
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && s.alive(s.boundless(x+j, y+i)) {
				alive++
			}
		}
	}
	return alive == 3 || alive == 2 && s.alive(x, y)
}
