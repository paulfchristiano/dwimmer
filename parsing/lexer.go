package parsing

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/term"
)

const eof = 0

//go:generate go tool yacc parser.y

type interpreter struct {
	names []string
}

func ParseAction(s string, names []string) *term.ActionC {
	interp := &interpreter{names: names}
	l := newLexer(s, interp, WANT_ACTION)
	yyParse(l)
	return l.getActionResult()
}

func ParseTerm(s string, names []string) term.C {
	interp := &interpreter{names: names}
	l := newLexer(s, interp, WANT_TERM)
	yyParse(l)
	return l.getTermResult()
}

type lexer struct {
	s            string
	first        int
	actionResult *term.ActionC
	termResult   term.C
	interpreter  *interpreter
	index        int
	last         int
}

func (r *lexer) parseWord(s string) ExprPart {
	if s[0] == '#' {
		s = s[1:]
	}
	for n, name := range r.interpreter.names {
		if name == s {
			return exprTerm{term.ReferenceC{n}}
		}
	}
	return exprText(s)
}

func newLexer(s string, interp *interpreter, first int) *lexer {
	return &lexer{s: s, interpreter: interp, first: first}
}

func (r *lexer) setActionResult(head string, e *Expr, n int) {
	var a term.ActionC
	var t term.C
	if e != nil {
		t = toC(e)
	}
	switch strings.ToLower(head) {
	case "return", "reply", "say", "respond", "answer":
		if t == nil {
			t = term.Cc(represent.Int(n))
		}
		a = term.ReturnC(t)
	case "ask", "inquire", "do":
		a = term.AskC(t)
	case "view", "check", "inspect", "see":
		if t == nil {
			t = term.Cc(represent.Int(n))
		}
		a = term.ViewC(t)
	case "replace", "rewrite", "change", "jump", "set":
		a = term.ReplaceC(t, n)
	case "replay", "redo", "repeat":
		a = term.ReplayC(n)
	case "correct", "fix", "debug":
		a = term.CorrectC(n)
	case "meta", "self", "here", "this":
		a = term.MetaC()
	case "close", "dismiss", "stop", "delete", "del", "remove":
		c := toC(e)
		switch c := c.(type) {
		case term.ReferenceC:
			a = term.DeleteC(c.Index)
		}
	default:
		return
	}
	r.actionResult = &a
}

func (r *lexer) setTransitiveActionResult(head string, object string, e *Expr) {
	var a term.ActionC
	switch strings.ToLower(head) {
	case "tell", "ask", "clarify", "push", "follow", "followup":
		switch o := r.parseWord(object).(type) {
		case exprTerm:
			c := toC(e)
			a = term.ClarifyC(o.val, c)
		default:
			return
		}
	}
	r.actionResult = &a
}

func (r *lexer) setTermResult(e *Expr) {
	r.termResult = toC(e)
}

func (r *lexer) getActionResult() *term.ActionC {
	return r.actionResult
}

func (r *lexer) getTermResult() term.C {
	return r.termResult
}

func (r *lexer) next() (next rune, ok bool) {
	if len(r.s) <= r.index {
		return
	}
	next, r.last = utf8.DecodeRuneInString(r.s[r.index:])
	ok = true
	r.index += r.last
	return
}

func (r *lexer) Error(s string) {
	//TODO do something intelligent with this
	//fmt.Printf("parse error: %s\n", s)
}

func (r *lexer) back() {
	if r.last == 0 {
		panic("backed up twice!")
	}
	r.index -= r.last
	r.last = 0
}

type state int

const (
	startS = state(iota)
	wordS
	numS
	proseS
	whiteS
)

func (r *lexer) Lex(lval *yySymType) int {
	if r.first != 0 {
		result := r.first
		r.first = 0
		return result
	}
	var b bytes.Buffer
	s := startS
	return r.fsm(lval, b, s)
}

func (r *lexer) backstop(lval *yySymType, b bytes.Buffer, s state) int {
	r.back()
	return r.stop(lval, b, s)
}

func (r *lexer) fsm(lval *yySymType, b bytes.Buffer, s state) int {
	c, ok := r.next()
	if !ok {
		return r.stop(lval, b, s)
	}
	switch s {
	case startS:
		switch {
		case alpha(c) || c == '_':
			s = wordS
			b.WriteRune(c)
		case c == '#':
			s = wordS
			b.WriteRune(c)
		case numeric(c):
			s = numS
			b.WriteRune(c)
		case symbol(c):
			return symbolToken(lval, c)
		case white(c):
			s = whiteS
			b.WriteRune(c)
		default:
			s = proseS
			b.WriteRune(c)
		}
	case wordS:
		switch {
		case alpha(c) || numeric(c) || c == '_':
			b.WriteRune(c)
		default:
			return r.backstop(lval, b, s)
		}
	case numS:
		switch {
		case numeric(c):
			b.WriteRune(c)
		case white(c) || symbol(c):
			return r.backstop(lval, b, s)
		default:
			s = proseS
			b.WriteRune(c)
		}
	case proseS:
		switch {
		case white(c) || symbol(c):
			return r.backstop(lval, b, s)
		default:
			b.WriteRune(c)
		}
	case whiteS:
		switch {
		case white(c):
			b.WriteRune(c)
		default:
			return r.backstop(lval, b, s)
		}
	default:
		panic("in an invalid state!")
	}
	return r.fsm(lval, b, s)
}

func symbolToken(lval *yySymType, c rune) int {
	lval.string = string(c)
	switch c {
	case '[', ']', '(', ')', '{', '}', ',', ':', '|', '"', '.', '@':
		return int(c)
	default:
		return SYMBOL
	}
}

func (r *lexer) stop(lval *yySymType, b bytes.Buffer, s state) int {
	switch s {
	case startS:
		return eof
	case whiteS:
		lval.string = b.String()
		return WHITE
	case proseS:
		lval.string = b.String()
		return PROSE
	case wordS:
		lval.string = b.String()
		switch lval.string {
		/*
			case "make", "Make", "MAKE", "hint", "Hint", "template", "m":
				return MAKE
		*/
		default:
			return WORD
		}
	case numS:
		val, _ := strconv.Atoi(b.String())
		lval.int = val
		return NUM
	default:
		panic("in an invalid state!")
	}
}

func white(c rune) bool {
	return unicode.IsSpace(c)
}

func numeric(c rune) bool {
	return unicode.IsNumber(c)
}

func alpha(c rune) bool {
	return unicode.IsLetter(c)
}

func symbol(c rune) bool {
	switch c {
	case '"', '[', ']', '(', ')', ':', ',', ';', '{', '}', '+':
		return true
	case '!', '/', '*', '-', '^', '&', '|', '?', '.', '@':
		return true
	default:
		return false
	}
}
