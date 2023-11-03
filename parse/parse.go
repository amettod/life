package parse

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

func rle(r io.Reader) (int, int, [][]int, error) {
	var x, y int
	state := [][]int{}
	row := []int{}
	digits := ""
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "x") {
			_, err := fmt.Sscanf(line, "x = %d, y = %d,", &x, &y)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("parse: %w", err)
			}
			continue
		}
		for _, r := range line {
			if unicode.IsDigit(r) {
				digits += string(r)
				continue
			}
			count := 1
			if len(digits) > 0 {
				c, err := strconv.Atoi(digits)
				if err != nil {
					return 0, 0, nil, fmt.Errorf("parse: %w", err)
				}
				count = c
				digits = ""
			}
			if r == ' ' {
				continue
			}
			if r == 'o' || r == 'b' {
				v := 0
				if r == 'o' {
					v = 1
				}
				for i := 0; i < count; i++ {
					row = append(row, v)
				}
				continue
			}
			if r == '$' {
				if len(row) > 0 {
					state = append(state, row)
					row = []int{}
				}
				for i := 1; i < count; i++ {
					state = append(state, []int{0})
				}
				continue
			}
			if r == '!' {
				state = append(state, row)
				break
			}
		}
	}
	if err := scan.Err(); err != nil {
		return 0, 0, nil, fmt.Errorf("parse: %w", err)
	}
	return x, y, state, nil
}

func cells(r io.Reader) ([][]int, error) {
	state := [][]int{}
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "!") {
			continue
		}
		row := []int{}
		for _, r := range line {
			if r == '.' || r == 'O' {
				v := 0
				if r == 'O' {
					v = 1
				}
				row = append(row, v)
			}
		}
		state = append(state, row)
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return state, nil
}

func life(r io.Reader) ([][]int, error) {
	type point struct{ x, y int }
	var (
		maxP, minP point
		once       bool
	)
	points := []point{}
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		var p point
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		_, err := fmt.Sscanf(line, "%d %d", &p.x, &p.y)
		if err != nil {
			return nil, err
		}
		if !once {
			maxP = point{
				x: p.x,
				y: p.y,
			}
			minP = point{
				x: p.x,
				y: p.y,
			}
			once = true
		}
		maxP = point{
			x: max(maxP.x, p.x),
			y: max(maxP.y, p.y),
		}
		minP = point{
			x: min(minP.x, p.x),
			y: min(minP.y, p.y),
		}
		points = append(points, p)
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	state := make([][]int, maxP.y-minP.y+1)
	for i := range state {
		state[i] = make([]int, maxP.x-minP.x+1)
	}
	for _, p := range points {
		state[p.y-minP.y][p.x-minP.x] = 1
	}
	return state, nil
}

func parse(r io.Reader, name string) ([][]int, error) {
	switch path.Ext(name) {
	case ".rle":
		_, _, state, err := rle(r)
		return state, err
	case ".cells":
		return cells(r)
	case ".life":
		return life(r)
	default:
		return nil, fmt.Errorf("parse: file %s is unsupported", name)
	}
}

func File(name string) ([][]int, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f, name)
}

func FileEmbed(fs embed.FS, name string) ([][]int, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f, name)
}
