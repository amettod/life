package life

import (
	"math"
	"reflect"
	"testing"
)

func Test_state_inside(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		s    state
		args args
		want bool
	}{
		{
			name: "in",
			s:    newState(1, 1),
			args: args{
				x: 0,
				y: 0,
			},
			want: true,
		},
		{
			name: "out x",
			s:    [][]int{},
			args: args{
				x: 1,
				y: 0,
			},
			want: false,
		},
		{
			name: "out y",
			s:    [][]int{},
			args: args{
				x: 0,
				y: 1,
			},
			want: false,
		},
		{
			name: "neg x",
			s:    [][]int{},
			args: args{
				x: -1,
				y: 0,
			},
			want: false,
		},
		{
			name: "neg y",
			s:    [][]int{},
			args: args{
				x: 0,
				y: -1,
			},
			want: false,
		},
		{
			name: "zero",
			s:    [][]int{},
			args: args{
				x: 0,
				y: 0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.inside(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("state.inside() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_state_boundless(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name  string
		s     state
		args  args
		want  int
		want1 int
	}{
		{
			name: "in",
			s:    newState(5, 5),
			args: args{
				x: 1,
				y: 1,
			},
			want:  1,
			want1: 1,
		},
		{
			name: "high",
			s:    newState(5, 5),
			args: args{
				x: 5,
				y: 5,
			},
			want:  0,
			want1: 0,
		},
		{
			name: "low",
			s:    newState(5, 5),
			args: args{
				x: -1,
				y: -1,
			},
			want:  4,
			want1: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.boundless(tt.args.x, tt.args.y)
			if got != tt.want {
				t.Errorf("state.boundless() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("state.boundless() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_state_cycleCalc(t *testing.T) {
	type args struct {
		x     int
		y     int
		alive bool
	}
	tests := []struct {
		name string
		s    state
		args args
		want state
	}{
		{
			name: "alive 0",
			s:    [][]int{{0}},
			args: args{
				x:     0,
				y:     0,
				alive: true,
			},
			want: [][]int{{1}},
		},
		{
			name: "alive 1",
			s:    [][]int{{1}},
			args: args{
				x:     0,
				y:     0,
				alive: true,
			},
			want: [][]int{{2}},
		},
		{
			name: "alive -2",
			s:    [][]int{{-2}},
			args: args{
				x:     0,
				y:     0,
				alive: true,
			},
			want: [][]int{{1}},
		},
		{
			name: "alive MaxInt",
			s:    [][]int{{math.MaxInt}},
			args: args{
				x:     0,
				y:     0,
				alive: true,
			},
			want: [][]int{{1}},
		},
		{
			name: "dead 0",
			s:    [][]int{{0}},
			args: args{
				x:     0,
				y:     0,
				alive: false,
			},
			want: [][]int{{0}},
		},
		{
			name: "dead 1",
			s:    [][]int{{2}},
			args: args{
				x:     0,
				y:     0,
				alive: false,
			},
			want: [][]int{{-1}},
		},
		{
			name: "dead MinInt",
			s:    [][]int{{math.MinInt}},
			args: args{
				x:     0,
				y:     0,
				alive: false,
			},
			want: [][]int{{0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.cycleCalc(tt.args.x, tt.args.y, tt.args.alive)
			got := tt.s
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("state = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_state_next(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		s    state
		args args
		want bool
	}{
		{
			name: "dead 1",
			s: [][]int{
				{1, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: false,
		},
		{
			name: "dead 2",
			s: [][]int{
				{1, 0, 0},
				{0, 0, 0},
				{0, 0, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: false,
		},
		{
			name: "died 4",
			s: [][]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: false,
		},
		{
			name: "died 7",
			s: [][]int{
				{1, 1, 1},
				{1, 1, 0},
				{1, 1, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: false,
		},
		{
			name: "alive 2",
			s: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{1, 0, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			name: "alive 3",
			s: [][]int{
				{1, 0, 0},
				{0, 1, 0},
				{1, 0, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			name: "birth",
			s: [][]int{
				{0, 1, 0},
				{0, 0, 0},
				{1, 0, 1},
			},
			args: args{
				x: 1,
				y: 1,
			},
			want: true,
		},
		{
			name: "boundless birth",
			s: [][]int{
				{1, 0, 0},
				{0, 0, 1},
				{0, 1, 0},
			},
			args: args{
				x: 0,
				y: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.next(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("state.next() = %v, want %v", got, tt.want)
			}
		})
	}
}
