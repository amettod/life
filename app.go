package life

type App struct {
	Game   *game
	Preset *presets
	Theme  *themes
}

func NewApp(w, h int, file string) (*App, error) {
	g := newGame(w, h)
	if file != "" {
		s, err := parseFile(file)
		if err != nil {
			return nil, err
		}
		g.SetState(0, 0, s)
	}

	p, err := newPresets()
	if err != nil {
		return nil, err
	}

	return &App{
		Game:   g,
		Preset: p,
		Theme:  newThemes(),
	}, nil
}
