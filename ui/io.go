package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var at = termbox.ColorDefault

type UIImplementer interface {
	Println(string)
	Debug(string)
	Readln(string, []string, map[rune]string) string
	GetCh() rune
	CheckCh() (rune, bool)
	SetCursor(int, int)
	PrintCh(rune)
	Clear()
	InitUI()
	CloseUI()
	Size() (int, int)
}

type Term struct {
	initialized bool
	input       chan termbox.Event
	x, y        int
}

func (t *Term) Check() {
	if !t.initialized {
		panic("using an uninitialized terminal!")
	}

}

type dummy struct{}

func Dummy() *dummy {
	return &dummy{}
}

func (d *dummy) Println(s string) { fmt.Println(s) }
func (d *dummy) Debug(s string)   { fmt.Println(s) }
func (d *dummy) GetCh() rune      { panic("asked for input") }
func (d *dummy) Clear()           {}
func (d *dummy) InitUI()          {}
func (d *dummy) CloseUI()         {}

func (d *dummy) Size() (int, int)    { return 0, 0 }
func (d *dummy) PrintCh(c rune)      {}
func (d *dummy) MoveCursor(x, y int) {}

func (d *dummy) Readln(s string, hints []string, tools map[rune]string) string {
	panic("asekd for input")
}
func (d *dummy) CheckCh() (rune, bool) { return 0, false }

func (t *Term) Debug(s string) {
	t.Println(s)
	Flush()
	_, height := termbox.Size()
	if t.y > height {
		t.GetCh()
		t.Clear()
	}
}

func (t *Term) Clear() {
	t.Check()
	t.x = 0
	t.y = 0
	termbox.Clear(at, at)
	termbox.Clear(at, at)
}

func (t *Term) Size() (int, int) {
	return termbox.Size()
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

func (t *Term) Readln(s string, hints []string, tools map[rune]string) string {
	t.Print(s)
	return t.Getln(hints, tools)
}

func (t *Term) Getln(hints []string, tools map[rune]string) string {
	Flush()
	startx, starty := t.x, t.y
	offsetx, offsety := 0, 0
	position := 0
	screenWidth, _ := t.Size()
	width := screenWidth - startx - 1
	input := ""
	hints = append(hints, "")
	hintIndex := len(hints) - 1

	//Helpers

	fromPosition := func(p int) (x, y int) {
		x = p % width
		y = p / width
		return
	}
	resetOffsets := func() {
		offsetx, offsety = fromPosition(position)
	}
	setPos := func(p int) {
		if p < 0 {
			p = 0
		}
		if p > len(input) {
			p = len(input)
		}
		position = p
		resetOffsets()
	}
	incPos := func() {
		setPos(position + 1)
	}
	decPos := func() {
		setPos(position - 1)
	}
	addChar := func(c rune) {
		input = fmt.Sprintf("%s%c%s", input[:position], c, input[position:])
		incPos()
	}
	splitLines := func(s string) (result []string) {
		for len(s) > 0 {
			if len(s) < width {
				result = append(result, s)
				s = ""
			} else {
				result = append(result, s[:width])
				s = s[width:]
			}
		}
		return result
	}
	var shownLines []string
	refresh := func() {
		for i := range shownLines {
			t.SetCursor(startx, starty+i)
			t.ClearRestOfLine()
		}
		shownLines = splitLines(input)
		for i, line := range shownLines {
			t.SetCursor(startx, starty+i)
			t.Print(line)
		}
		resetOffsets()
		t.SetCursor(startx+offsetx, starty+offsety)
		Flush()
	}
	scan := func(d int, c rune, offset int) int {
		try := position
		for try >= 0 && try <= len(input) {
			try += d
			if try >= 0 && try < len(input) {
				if int(input[try]) == int(c) && try+offset >= 0 && try+offset <= len(input) {
					return try + offset
				}
			}
		}
		if try < 0 {
			return 0
		}
		if try > len(input) {
			return len(input)
		}
		return try
	}
	exceeds := func(d int, down, up rune, offset int) int {
		stack := 0
		try := position + offset
		for stack >= 0 {
			try += d
			if try < 0 {
				try = -1
				break
			}
			if try >= len(input) {
				try = len(input)
				break
			}
			if rune(input[try]) == down {
				stack--
			}
			if rune(input[try]) == up {
				stack++
			}
		}
		return try
	}

	//Main loop

	for {
		if ev := <-t.input; ev.Type == termbox.EventKey {
			var ch rune
			switch {
			case ev.Key == termbox.KeyEnter:
				t.Newln()
				return input
			case ev.Key == termbox.KeyCtrlC:
				panic("interrupted")
			case ev.Key == termbox.KeyCtrlF:
				selector := t.GetCh()
				tool, used := tools[selector]
				if used {
					for _, c := range tool {
						addChar(c)
					}
				}
			case ev.Key == termbox.KeyCtrlH:
				setPos(scan(-1, ']', 0))
			case ev.Key == termbox.KeyCtrlL:
				setPos(scan(1, '[', 1))
			case ev.Key == termbox.KeyCtrlS:
				panic("meta")
			case ev.Key == termbox.KeyCtrlE:
				left := exceeds(-1, '[', ']', 0)
				if left < len(input) {
					right := exceeds(1, ']', '[', -1)
					setPos(left + 1)
					input = fmt.Sprintf("%s%s", input[:left+1], input[right:])
				}
			case ev.Key == termbox.KeyArrowUp:
				if hintIndex > 0 {
					if hintIndex == len(hints)-1 {
						hints[len(hints)-1] = input
					}
					hintIndex--
				}
				input = hints[hintIndex]
				setPos(position)
			case ev.Key == termbox.KeyArrowDown:
				if hintIndex < len(hints)-1 {
					hintIndex++
				}
				input = hints[hintIndex]
				setPos(position)
			case ev.Key == termbox.KeyDelete:
				if position < len(input) {
					input = input[:position] + input[position+1:]
				}
			case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
				if position > 0 {
					input = input[:position-1] + input[position:]
				}
				decPos()
			case ev.Key == termbox.KeyArrowLeft:
				decPos()
			case ev.Key == termbox.KeyArrowRight:
				incPos()
			case ev.Key == termbox.KeySpace:
				ch = ' '
			default:
				ch = ev.Ch
			}
			if ch != 0 {
				addChar(ch)
			}
			refresh()
		}
	}
}

func (t *Term) GetCh() rune {
	Flush()
	for {
		if ev := <-t.input; ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyCtrlC {
				panic("interrupted")
			}
			return ev.Ch
		}
	}
}

func getInput(ch chan termbox.Event) {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyCtrlD {
			panic("interrupted with EOF")
		}
		ch <- ev
	}
}

func (t *Term) InitUI() {
	t.input = make(chan termbox.Event)
	go func() {
		defer termbox.Close()
		getInput(t.input)
	}()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	t.initialized = true
	t.SetCursor(0, 0)
	Flush()
}

func (t *Term) CheckCh() (rune, bool) {
	select {
	case ev := <-t.input:
		if ev.Type == termbox.EventKey && ev.Ch != 0 {
			return ev.Ch, true
		}
	default:
	}
	return 0, false
}

func (t *Term) CloseUI() {
	termbox.Close()
}

func Flush() {
	termbox.Flush()
}
