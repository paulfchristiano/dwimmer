package meta

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	PutChar = term.Make("put the character [] at the location with [] spots to its left " +
		"and [] spots above it, and leave the cursor there")
	Clear   = term.Make("clear the screen")
	Size    = term.Make("what are the dimensions of the screen?")
	GetChar = term.Make("refresh the screen, wait until the user types something, " +
		"and return the next character that the user types")
	SetCursor = term.Make("put the cursor at the location with [] spots to its left " +
		"and [] spots above it")
	GetCursor = term.Make("where is the cursor?")
)

func init() {
	var s *term.SettingS
	s = dynamics.ExpectQuestion(term.InitS(), PutChar, "Q", "char", "x", "y")
	dynamics.AddNative(s, dynamics.Args3(nativePutChar), "char", "x", "y")

	s = dynamics.ExpectQuestion(term.InitS(), SetCursor, "Q", "x", "y")
	dynamics.AddNative(s, dynamics.Args2(nativeSetCursor), "x", "y")

	s = dynamics.ExpectQuestion(term.InitS(), GetCursor, "Q")
	dynamics.AddNative(s, dynamics.Args0(nativeGetCursor))

	s = dynamics.ExpectQuestion(term.InitS(), Clear, "Q")
	dynamics.AddNative(s, dynamics.Args0(nativeClear))

	s = dynamics.ExpectQuestion(term.InitS(), Size, "Q")
	dynamics.AddNative(s, dynamics.Args0(nativeSize))

	s = dynamics.ExpectQuestion(term.InitS(), GetChar, "Q")
	dynamics.AddNative(s, dynamics.Args0(nativeGetChar))
}

func nativeClear(d dynamics.Dwimmer, s *term.SettingT) term.T {
	d.Clear()
	return core.OK.T()
}

func nativeGetCursor(d dynamics.Dwimmer, s *term.SettingT) term.T {
	x, y := d.GetCursor()
	return Pos.T(XYPos.T(represent.Int(x), represent.Int(y)))
}

func nativeGetChar(d dynamics.Dwimmer, s *term.SettingT) term.T {
	c, key := d.GetCh()
	if c != 0 {
		return core.Answer.T(represent.Rune(c))
	}
	return KeyEntered.T(term.Int(int(key)))
}

var (
	KeyEntered = term.Make("the user did not enter a character, but entered a key with termbox Key code []")
)

func nativeSize(d dynamics.Dwimmer, s *term.SettingT) term.T {
	x, y := d.Size()
	return WidthHeight.T(represent.Int(x), represent.Int(y))
}

func nativePutChar(d dynamics.Dwimmer, s *term.SettingT, char, xt, yt term.T) term.T {
	c, err := represent.ToRune(d, char)
	if err != nil {
		return term.Make("asked to write character, but received [] " +
			"while converting to a character").T(err)
	}
	x, err := represent.ToInt(d, xt)
	if err != nil {
		return term.Make("asked to write character, but received [] "+
			"while converting coordinate [] to an integer").T(err, xt)
	}
	y, err := represent.ToInt(d, yt)
	if err != nil {
		return term.Make("asked to write character, but received [] "+
			"while converting coordinate [] to an integer").T(err, yt)
	}
	d.SetCursor(x, y)
	d.PrintCh(c)
	return core.OK.T()
}

func nativeSetCursor(d dynamics.Dwimmer, s *term.SettingT, xt, yt term.T) term.T {
	x, err := represent.ToInt(d, xt)
	if err != nil {
		return term.Make("asked to move cursor, but received [] "+
			"while converting coordinate [] to an integer").T(err, xt)
	}
	y, err := represent.ToInt(d, yt)
	if err != nil {
		return term.Make("asked to move cursor, but received [] "+
			"while converting coordinate [] to an integer").T(err, yt)
	}
	d.SetCursor(x, y)
	return core.OK.T()
}

var (
	WidthHeight = term.Make("the width of the output pane is [] and the height is []")
	XYPos       = term.Make("the position with [] spots to its left and [] spots above it")
	Pos         = term.Make("the cursor is at the position []")
)
