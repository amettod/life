package term

import (
	"fmt"
	"io"
	"strings"
)

const (
	escUp            = "\x1b[1A"
	escClearLines    = "\x1b[2K"
	escStartOfLine   = "\x1b[1G"
	escBackgroundRGB = "\x1b[48;2;%d;%d;%dm"
	escForegroundRGB = "\x1b[38;2;%d;%d;%dm"
	escReset         = "\x1b[0m"
)

type rgb interface {
	Color() (uint8, uint8, uint8)
}

type Term interface {
	// Print to io.Writer.
	Print()
	// Write color the background and foreground.
	Write(background, foreground rgb, s string)
	// Write new line.
	WriteLn()
}

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

func (t *term) Print() {
	t.clear()
	content := t.s.String()
	t.s.Reset()
	fmt.Fprint(t.w, content)
	t.lines = strings.Count(content, "\n")
}

func (t *term) Write(background, foreground rgb, s string) {
	var b, f string
	if background != nil {
		b = rgbTo(escBackgroundRGB, background)
	}
	if foreground != nil {
		f = rgbTo(escForegroundRGB, foreground)
	}
	fmt.Fprint(&t.s, b, f, s, escReset)
}

func (t *term) WriteLn() {
	fmt.Fprintln(&t.s)
}

func rgbTo(format string, c rgb) string {
	r, g, b := c.Color()
	return fmt.Sprintf(format, r, g, b)
}
