package life

import "sort"

type RGB struct {
	r uint8
	g uint8
	b uint8
}

func NewRGB(r, g, b uint8) RGB {
	return RGB{
		r: r,
		g: g,
		b: b,
	}
}

func (c RGB) Color() (uint8, uint8, uint8) {
	return c.r, c.g, c.b
}

type theme struct {
	name       string
	background RGB
	foreground RGB
	alive      []RGB
	dead       []RGB
}

type themes struct {
	current int
	store   []theme
}

func newThemes() *themes {
	t := &themes{
		current: 0,
		store: []theme{
			{
				name:       "color16",
				background: NewRGB(255, 255, 255),
				foreground: NewRGB(0, 0, 0),
				alive: []RGB{
					NewRGB(0, 0, 0),
					NewRGB(0, 0, 95),
					NewRGB(0, 0, 135),
					NewRGB(0, 0, 175),
					NewRGB(0, 0, 215),
					NewRGB(0, 0, 255),
				},
				dead: []RGB{
					NewRGB(125, 255, 255),
					NewRGB(135, 255, 255),
					NewRGB(145, 255, 255),
					NewRGB(165, 255, 255),
					NewRGB(175, 255, 255),
					NewRGB(185, 255, 255),
					NewRGB(195, 255, 255),
					NewRGB(205, 255, 255),
					NewRGB(215, 255, 255),
				},
			},
			{
				name:       "orangeAndRed",
				background: NewRGB(255, 255, 255),
				foreground: NewRGB(0, 0, 0),
				alive: []RGB{
					NewRGB(255, 135, 0),
					NewRGB(255, 120, 0),
					NewRGB(255, 105, 0),
					NewRGB(255, 90, 0),
					NewRGB(255, 75, 0),
					NewRGB(255, 60, 0),
					NewRGB(255, 45, 0),
					NewRGB(255, 30, 0),
					NewRGB(255, 15, 0),
					NewRGB(255, 0, 0),
				},
				dead: []RGB{
					NewRGB(255, 155, 155),
					NewRGB(255, 165, 165),
					NewRGB(255, 175, 175),
					NewRGB(255, 185, 185),
					NewRGB(255, 195, 195),
					NewRGB(255, 205, 205),
					NewRGB(255, 215, 215),
					NewRGB(255, 225, 225),
					NewRGB(255, 235, 235),
					NewRGB(255, 245, 245),
				},
			},
			{
				name:       "whiteAndBlack",
				background: NewRGB(255, 255, 255),
				foreground: NewRGB(0, 0, 0),
				alive: []RGB{
					NewRGB(0, 0, 0),
				},
				dead: []RGB{
					NewRGB(255, 255, 255),
				},
			},
			{
				name:       "blackAndWhite",
				background: NewRGB(0, 0, 0),
				foreground: NewRGB(255, 255, 255),
				alive: []RGB{
					NewRGB(255, 255, 255),
				},
				dead: []RGB{
					NewRGB(0, 0, 0),
				},
			},
			{
				name:       "ocean",
				background: NewRGB(0, 0, 130),
				foreground: NewRGB(255, 255, 255),
				alive: []RGB{
					NewRGB(75, 75, 255),
					NewRGB(85, 85, 255),
					NewRGB(95, 95, 255),
					NewRGB(105, 105, 255),
					NewRGB(115, 115, 255),
					NewRGB(125, 125, 255),
					NewRGB(135, 135, 255),
					NewRGB(145, 145, 255),
					NewRGB(255, 255, 255),
				},
				dead: []RGB{
					NewRGB(0, 0, 70),
					NewRGB(0, 0, 80),
					NewRGB(0, 0, 90),
					NewRGB(0, 0, 100),
					NewRGB(0, 0, 110),
					NewRGB(0, 0, 120),
				},
			},
			{
				name:       "fire",
				background: NewRGB(130, 0, 0),
				foreground: NewRGB(255, 255, 0),
				alive: []RGB{
					NewRGB(255, 0, 0),
					NewRGB(255, 25, 0),
					NewRGB(255, 50, 0),
					NewRGB(255, 75, 0),
					NewRGB(255, 100, 0),
					NewRGB(255, 210, 0),
					NewRGB(255, 220, 0),
					NewRGB(255, 230, 0),
					NewRGB(255, 245, 0),
					NewRGB(255, 255, 0),
				},
				dead: []RGB{
					NewRGB(70, 0, 0),
					NewRGB(80, 0, 0),
					NewRGB(90, 0, 0),
					NewRGB(100, 0, 0),
					NewRGB(110, 0, 0),
					NewRGB(120, 0, 0),
				},
			},
			{
				name:       "matrix",
				background: NewRGB(0, 0, 0),
				foreground: NewRGB(0, 255, 0),
				alive: []RGB{
					NewRGB(0, 205, 0),
					NewRGB(0, 215, 0),
					NewRGB(0, 225, 0),
					NewRGB(0, 235, 0),
					NewRGB(0, 245, 0),
					NewRGB(0, 255, 0),
				},
				dead: []RGB{
					NewRGB(0, 70, 0),
					NewRGB(0, 60, 0),
					NewRGB(0, 50, 0),
					NewRGB(0, 40, 0),
					NewRGB(0, 30, 0),
					NewRGB(0, 20, 0),
					NewRGB(0, 10, 0),
				},
			},
		},
	}
	t.sort()
	return t
}

func (t *themes) theme() theme {
	return t.store[t.current]
}

func (t *themes) sort() {
	sort.Slice(t.store, func(i, j int) bool {
		return t.store[i].name < t.store[j].name
	})
}

func (t *themes) alive(cycle int) RGB {
	c := cycle - 1
	l := len(t.theme().alive)
	if c < l {
		return t.theme().alive[c]
	}
	return t.theme().alive[l-1]
}

func (t *themes) dead(cycle int) RGB {
	c := cycle*-1 - 1
	l := len(t.theme().dead)
	if c < l {
		return t.theme().dead[c]
	}
	return t.theme().background
}

// Background color.
func (t *themes) Background() RGB {
	return t.theme().background
}

// Foreground color.
func (t *themes) Foreground() RGB {
	return t.theme().foreground
}

// Color return depending on cycle.
func (t *themes) Color(cycle int) RGB {
	switch {
	case cycle > 0:
		return t.alive(cycle)
	case cycle < 0:
		return t.dead(cycle)
	default:
		return t.Background()
	}
}

// Name theme.
func (t *themes) Name() string {
	return t.theme().name
}

// Next theme.
func (t *themes) Next() {
	t.current++
	if t.current > len(t.store)-1 {
		t.current = 0
	}
}
