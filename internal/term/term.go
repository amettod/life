package term

import (
	"fmt"
	"io"
	"strings"
)

type rgb interface {
	Color() (uint8, uint8, uint8)
}

type Term interface {
	SetCell(background rgb, r rune)
	SetInfo(background, foreground rgb, b byte)
	SetLn()
	Print()
}

const (
	escUp            = "\x1b[1A"
	escClearLines    = "\x1b[2K"
	escStartOfLine   = "\x1b[1G"
	escBackgroundRGB = "\x1b[48;2;%d;%d;%dm"
	escForegroundRGB = "\x1b[38;2;%d;%d;%dm"
	escReset         = "\x1b[0m"
)

type term struct {
	w     io.Writer
	s     strings.Builder
	lines int
}

func New(w io.Writer) Term {
	return &term{
		w:     w,
		s:     strings.Builder{},
		lines: 0,
	}
}

func (t *term) Print() {
	t.clear()
	content := t.s.String()
	t.s.Reset()
	fmt.Fprint(t.w, content)
	t.lines = strings.Count(content, "\n")
}

func (t *term) clear() {
	if t.lines == 0 {
		fmt.Fprint(t.w, escStartOfLine)
		fmt.Fprint(t.w, escClearLines)
		return
	}
	for i := 0; i < t.lines; i++ {
		fmt.Fprint(t.w, escUp)
		fmt.Fprint(t.w, escClearLines)
	}
}

func (t *term) SetInfo(background, foreground rgb, b byte) {
	fmt.Fprint(
		&t.s,
		backgroundRGB(background.Color()),
		foregroundRGB(foreground.Color()),
		string(b),
		escReset,
	)
}

func (t *term) SetCell(background rgb, r rune) {
	fmt.Fprint(
		&t.s,
		backgroundRGB(background.Color()),
		string(r),
		escReset,
	)
}

func (t *term) SetLn() {
	fmt.Fprintln(&t.s)
}

func backgroundRGB(r, g, b uint8) string {
	return fmt.Sprintf(escBackgroundRGB, r, g, b)
}

func foregroundRGB(r, g, b uint8) string {
	return fmt.Sprintf(escForegroundRGB, r, g, b)

}
