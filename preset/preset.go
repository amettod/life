package preset

import (
	"embed"
	"sort"
)

//go:embed presets/*
var EmbedFS embed.FS

const EmbedDir = "presets"

type Preset interface {
	// Append preset.
	Append(name string, state [][]int)
	// Name preset.
	Name() string
	// Next preset.
	Next()
	// State return current preset.
	State() [][]int
}

type preset struct {
	name  string
	state [][]int
}

type presets struct {
	current int
	store   []preset
}

func New() Preset {
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
	p.sort()
	return p
}

func (p *presets) sort() {
	sort.Slice(p.store, func(i, j int) bool {
		return p.store[i].name < p.store[j].name
	})
}

func (p *presets) Append(name string, state [][]int) {
	p.store = append(p.store, preset{
		name:  name,
		state: state,
	})
	p.sort()
}

func (p *presets) Name() string {
	return p.store[p.current].name
}

func (p *presets) Next() {
	p.current++
	if p.current > len(p.store)-1 {
		p.current = 0
	}
}

func (p *presets) State() [][]int {
	return p.store[p.current].state
}
