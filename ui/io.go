package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var at = termbox.ColorDefault

type UIImplementer interface {
	Println(string)
	Debug(string)
	Readln(string, ...[]string) string
	GetCh() rune
	Clear()
	InitUI()
	CloseUI()
}

type Term struct {
	initialized bool
	x, y        int
}

func (t *Term) Check() {
	if !t.initialized {
		panic("using an uninitialized terminal!")
	}

}

func (t *Term) Debug(s string) {
	t.Println(s)
	Flush()
}

func (t *Term) Clear() {
	t.Check()
	t.x = 0
	t.y = 0
	termbox.Clear(at, at)
	termbox.Clear(at, at)
	termbox.Flush()
}

func (t *Term) SetCursor(x, y int) {
	t.Check()
	t.x = x
	t.y = y
	termbox.SetCursor(x, y)
}

func (t *Term) MoveCursor(dx, dy int) {
	t.SetCursor(t.x+dx, t.y+dy)
}

func (t *Term) PrintCh(ch rune) {
	t.Print(string([]rune{ch}))
}

func (t *Term) Print(s string) {
	n := t.PrintNoMove(s)
	t.MoveCursor(n, 0)
}

func (t *Term) PrintNoMove(s string) int {
	t.Check()
	x, y := t.x, t.y
	for _, ch := range s {
		termbox.SetCell(x, y, ch, at, at)
		x++
	}
	return x - t.x
}

func (t *Term) Println(s string) {
	t.Print(s)
	t.Newln()
}

func (t *Term) Newln() {
	t.SetCursor(0, t.y+1)
}

func (t *Term) ClearRestOfLine() {
	t.Check()
	width, _ := termbox.Size()
	for i := t.x; i < t.x+width; i++ {
		termbox.SetCell(i, t.y, ' ', at, at)
	}
}

func (t *Term) Readln(s string, hintss ...[]string) string {
	hints := []string{}
	for _, hs := range hintss {
		hints = append(hints, hs...)
	}
	t.Print(s)
	return t.Getln(hints)
}

func (t *Term) Getln(hints []string) string {
	Flush()
	startx, starty := t.x, t.y
	offsetx, offsety := 0, 0
	input := ""
	hints = append(hints, "")
	hintIndex := len(hints) - 1
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			var ch rune
			switch {
			case ev.Key == termbox.KeyEnter:
				t.Newln()
				return input
			case ev.Key == termbox.KeyCtrlC:
				panic("interrupted")
			case ev.Key == termbox.KeyArrowUp:
				if hintIndex > 0 {
					if hintIndex == len(hints)-1 {
						hints[len(hints)-1] = input
					}
					hintIndex--
				}
				input = hints[hintIndex]
				if offsetx > len(input) {
					offsetx = len(input)
				}
			case ev.Key == termbox.KeyArrowDown:
				if hintIndex < len(hints)-1 {
					hintIndex++
				}
				input = hints[hintIndex]
				if offsetx > len(input) {
					offsetx = len(input)
				}
			case ev.Key == termbox.KeyDelete:
				if offsetx < len(input) {
					input = input[:offsetx] + input[offsetx+1:]
				}
			case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
				if offsetx > 0 {
					input = input[:offsetx-1] + input[offsetx:]
					offsetx--
				}
			case ev.Key == termbox.KeyArrowLeft:
				if offsetx > 0 {
					offsetx--
				}
			case ev.Key == termbox.KeyArrowRight:
				if offsetx < len(input) {
					offsetx++
				}
			case ev.Key == termbox.KeySpace:
				ch = ' '
			default:
				ch = ev.Ch
			}
			if ch != 0 {
				input = fmt.Sprintf("%s%c%s", input[:offsetx], ch, input[offsetx:])
				offsetx++
			}
			t.SetCursor(startx, starty)
			t.ClearRestOfLine()
			t.PrintNoMove(input)
			t.MoveCursor(offsetx, offsety)
			Flush()
		}

	}
}

func (t *Term) GetCh() rune {
	Flush()
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyCtrlC {
				panic("interrupted!")
			}
			return ev.Ch
		}
	}
}

func (t *Term) InitUI() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	t.initialized = true
	t.SetCursor(0, 0)
	Flush()
}

func (t *Term) CloseUI() {
	termbox.Close()
}

func Flush() {
	termbox.Flush()
}
