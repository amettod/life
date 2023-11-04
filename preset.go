package life

import (
	"embed"
	"path"
	"sort"
	"strings"
)

//go:embed presets/*
var embedFS embed.FS

const embedDir = "presets"

type preset struct {
	name  string
	state [][]int
}

type presets struct {
	current int
	store   []preset
}

func newPresets() (*presets, error) {
	p := &presets{
		store: []preset{
			{
				name: "cross",
				state: [][]int{
					{0, 0, 1, 0, 0},
					{0, 0, 1, 0, 0},
					{0, 0, 1, 0, 0},
					{1, 1, 0, 1, 1},
					{0, 0, 1, 0, 0},
					{0, 0, 1, 0, 0},
					{0, 0, 1, 0, 0},
				},
			},
			{
				name: "donut",
				state: [][]int{
					{0, 1, 0},
					{1, 0, 1},
					{0, 1, 0},
				},
			},
			{
				name: "quotes",
				state: [][]int{
					{0, 1, 1},
					{0, 0, 1},
					{1, 0, 0},
					{1, 1, 0},
				},
			},
			{
				name: "stone",
				state: [][]int{
					{1, 1},
					{1, 1},
				},
			},
		},
	}
	if err := p.load(); err != nil {
		return nil, err
	}
	p.sort()
	return p, nil
}

func (p *presets) sort() {
	sort.Slice(p.store, func(i, j int) bool {
		return p.store[i].name < p.store[j].name
	})
}

func (p *presets) load() error {
	files, err := embedFS.ReadDir(embedDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			s, err := parseFileEmbed(embedFS, path.Join(embedDir, f.Name()))
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			p.store = append(p.store, preset{name: name, state: s})
		}
	}
	return nil
}

// Name preset.
func (p *presets) Name() string {
	return p.store[p.current].name
}

// Next preset.
func (p *presets) Next() {
	p.current++
	if p.current > len(p.store)-1 {
		p.current = 0
	}
}

// State return current preset.
func (p *presets) State() [][]int {
	return p.store[p.current].state
}
