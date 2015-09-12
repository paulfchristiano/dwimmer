package parsing

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/term"
)

type Expr struct {
	head      []string
	args      []term.C
	singleton bool
	keys      []string
	varIndex  int
}

type ExprPart interface {
	isExprPart()
	String() string
}

func exprFromTerm(t term.C) *Expr {
	result := EmptyExpr()
	result.append(exprTerm{t})
	return result
}

type exprText string
type exprTerm struct {
	val term.C
}

func (c *Expr) isExprPart() {
}

func (c exprTerm) isExprPart() {
}

func (c exprTerm) String() string {
	return c.val.String()
}

func (c exprText) isExprPart() {
}

func (c exprText) String() string {
	return string(c)
}

func (c *Expr) String() string {
	return fmt.Sprintf("E(%v)", toC(c))
}

func (c *Expr) addArg(t term.C) {
	if len(c.head) == 1 && c.head[0] == "" {
		c.singleton = true
	}
	c.head = append(c.head, "")
	c.args = append(c.args, t)
}

func (c *Expr) append(x ExprPart) {
	c.singleton = false
	switch x := x.(type) {
	case (*Expr):
		c.addArg(toC(x))
	case (exprText):
		z := &c.head[len(c.head)-1]
		*z = *z + string(x)
	case (exprTerm):
		c.addArg(x.val)
	default:
		panic("tried to add an ExprPart of an unknown type")
	}
}

func toC(c *Expr) term.C {
	if c.singleton {
		return c.args[0]
	}
	return term.Make(c.head...).C(c.args...)
}

func EmptyExpr() *Expr {
	return &Expr{
		head: []string{""},
		args: make([]term.C, 0),
	}
}

const vars = "xyzwijklmntuv"

func (c *Expr) nextName() string {
	i := c.varIndex
	c.varIndex++
	return string(vars[i])
}
