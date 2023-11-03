package parse

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_rle(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		want2   [][]int
		wantErr bool
	}{
		{
			name: "grin",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#N Grin",
							"#C A common parent of the block.",
							"#C https://www.conwaylife.com/wiki/index.php?title=Grin",
							"x = 4, y = 2, rule = B3/S23",
							"o2bo$b2o!",
						},
						"\n",
					),
				),
			},
			want:  4,
			want1: 2,
			want2: [][]int{
				{1, 0, 0, 1},
				{0, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "hat",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#N Hat",
							"#C A 9-cell still life.",
							"#C https://www.conwaylife.com/wiki/index.php?title=Hat",
							"x = 5, y = 4, rule = B3/S23",
							"2bo2b$bobob$bobob$2ob2o!",
						},
						"\n",
					),
				),
			},
			want:  5,
			want1: 4,
			want2: [][]int{
				{0, 0, 1, 0, 0},
				{0, 1, 0, 1, 0},
				{0, 1, 0, 1, 0},
				{1, 1, 0, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "hat split line",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#N Hat",
							"#C A 9-cell still life.",
							"#C https://www.conwaylife.com/wiki/index.php?title=Hat",
							"x = 5, y = 4, rule = B3/S23",
							"2bo2b",
							"$bobo",
							"b$bob",
							"ob$2o",
							"b2o!",
						},
						"\n",
					),
				),
			},
			want:  5,
			want1: 4,
			want2: [][]int{
				{0, 0, 1, 0, 0},
				{0, 1, 0, 1, 0},
				{0, 1, 0, 1, 0},
				{1, 1, 0, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "heart",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#N Heart",
							"#C A period 5 oscillator.",
							"#C www.conwaylife.com/wiki/Heart",
							"x = 11, y = 11, rule = b3/s23",
							"5bo5b$4bo2bo3b$bo2bo2bo3b$obobobob2ob$bo2bo2bo3b$4bo5bo$5b5ob2$7bo3b$",
							"6bobo2b$7bo!",
						},
						"\n",
					),
				),
			},
			want:  11,
			want1: 11,
			want2: [][]int{
				{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0},
				{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
				{1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0},
				{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0},
				{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
				{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0},
				{0},
				{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := rle(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("rle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("rle() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("rle() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_cells(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    [][]int
		wantErr bool
	}{
		{
			name: "4blocks",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"! 4blocks.cells",
							"! https://conwaylife.com/wiki/Density",
							"! https://www.conwaylife.com/patterns/4blocks.cells",
							"OO.OO",
							"OO.OO",
							".....",
							"OO.OO",
							"OO.OO",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{1, 1, 0, 1, 1},
				{1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0},
				{1, 1, 0, 1, 1},
				{1, 1, 0, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "4 boats",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"!Name: 4 boats",
							"!A period 2 oscillator made up of 4 boats.",
							"...O",
							"..O.O",
							".O.OO",
							"O.O..OO",
							".OO..O.O",
							"...OO.O",
							"...O.O",
							"....O",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{0, 0, 0, 1},
				{0, 0, 1, 0, 1},
				{0, 1, 0, 1, 1},
				{1, 0, 1, 0, 0, 1, 1},
				{0, 1, 1, 0, 0, 1, 0, 1},
				{0, 0, 0, 1, 1, 0, 1},
				{0, 0, 0, 1, 0, 1},
				{0, 0, 0, 0, 1},
			},
			wantErr: false,
		},
		{
			name: "schickengine",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"! 1x256schickengine.cells",
							"! https://conwaylife.com/wiki/One-cell-thick_pattern",
							"! https://www.conwaylife.com/patterns/1x256schickengine.cells",
							"OOOOO.OOOO..OOO..OOOOO.OOOO.OOOO..........................OOOOO",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{
					1, 1, 1, 1, 1, 0, 1, 1, 1, 1,
					0, 0, 1, 1, 1, 0, 0, 1, 1, 1,
					1, 1, 0, 1, 1, 1, 1, 0, 1, 1,
					1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
					1, 1, 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cells(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("cells() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cells() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_life(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    [][]int
		wantErr bool
	}{
		{
			name: "glider",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#Life 1.06",
							"0 -1",
							"1 0",
							"-1 1",
							"0 1",
							"1 1",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{0, 1, 0},
				{0, 0, 1},
				{1, 1, 1},
			},
			wantErr: false,
		},
		{
			name: "negative stone",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#Life 1.06",
							"-1 -1",
							"-2 -1",
							"-1 -2",
							"-2 -2",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{1, 1},
				{1, 1},
			},
			wantErr: false,
		},
		{
			name: "positive stone",
			args: args{
				strings.NewReader(
					strings.Join(
						[]string{
							"#Life 1.06",
							"1 1",
							"2 1",
							"1 2",
							"2 2",
						},
						"\n",
					),
				),
			},
			want: [][]int{
				{1, 1},
				{1, 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := life(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("life() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("life() = %v, want %v", got, tt.want)
			}
		})
	}
}
