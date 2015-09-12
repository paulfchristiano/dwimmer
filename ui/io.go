package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var at = termbox.ColorDefault

func Clear() {
	cursorx = 0
	cursory = 0
	termbox.Clear(at, at)
	termbox.Flush()
	/*
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	*/
}

var cursorx, cursory int

func SetCursor(x, y int) {
	cursorx = x
	cursory = y
	termbox.SetCursor(x, y)
}

func MoveCursor(dx, dy int) {
	SetCursor(cursorx+dx, cursory+dy)
}

func PrintCh(ch rune) {
	Print(string([]rune{ch}))
}

func Print(s string) {
	n := PrintNoMove(s)
	MoveCursor(n, 0)
}

func PrintNoMove(s string) int {
	x, y := cursorx, cursory
	for _, ch := range s {
		termbox.SetCell(x, y, ch, at, at)
		x++
	}
	return x - cursorx
}

func Println(s string) {
	Print(s)
	Newln()
}

func Newln() {
	SetCursor(0, cursory+1)
}

func ClearRestOfLine() {
	width, _ := termbox.Size()
	for i := cursorx; i < cursorx+width; i++ {
		termbox.SetCell(i, cursory, ' ', at, at)
	}
}

func GetLine(hints []string) string {
	Flush()
	startx, starty := cursorx, cursory
	offsetx, offsety := 0, 0
	input := ""
	hints = append(hints, "")
	hintIndex := len(hints) - 1
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			var ch rune
			switch {
			case ev.Key == termbox.KeyEnter:
				Newln()
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
			SetCursor(startx, starty)
			ClearRestOfLine()
			PrintNoMove(input)
			MoveCursor(offsetx, offsety)
			Flush()
		}

	}
}

func GetCh() rune {
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

func Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	SetCursor(0, 0)
	termbox.Flush()
}

func Close() {
	termbox.Close()
}

func Flush() {
	termbox.Flush()
}
