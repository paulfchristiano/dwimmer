package term

import (
	"fmt"
	"reflect"
)

type Str string
type Int int
type Channel struct {
	Setting *SettingT
}
type Quoted struct {
	Value T
}
type Wrapper struct {
	Value reflect.Value
}

var (
	strHead     = Make("a string represented by a Go object")
	intHead     = Make("an integer represented by a Go object")
	quotedHead  = Make("a term represented by a Go object")
	wrapperHead = Make("a Go object")
	channelHead = Make("a channel to a state represented by a Go object")
)

//Channels

func (c Channel) Head() TemplateId {
	return channelHead
}

func (c Channel) String() string {
	return fmt.Sprintf("->")
}

func (c Channel) Values() []T {
	return []T{}
}

func (c Channel) Instantiate() *SettingT {
	return c.Setting.Copy()
}

//Strings

func (s Str) Head() TemplateId {
	return strHead
}

func (s Str) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (s Str) Values() []T {
	return make([]T, 0)
}

//Integers

func (n Int) Head() TemplateId {
	return intHead
}

func (n Int) String() string {
	return fmt.Sprintf("%d", int(n))
}

func (n Int) Values() []T {
	return make([]T, 0)
}

//Terms

func Quote(v T) T {
	return Quoted{v}
}

func (q Quoted) Head() TemplateId {
	return quotedHead
}

func (q Quoted) String() string {
	return fmt.Sprintf("T(%v)", q.Value)
}

func (q Quoted) Values() []T {
	return make([]T, 0)
}

//Wrapper

func Wrap(i interface{}) T {
	return Wrapper{reflect.ValueOf(i)}
}

func (w Wrapper) Head() TemplateId {
	return wrapperHead
}

func (w Wrapper) String() string {
	return fmt.Sprintf("GoObject(%v)", w.Value)
}

func (w Wrapper) Values() []T {
	return make([]T, 0)
}

//S and C Constants

type ConstS struct {
	Val T
}
type ConstC struct {
	Val T
}

func (c ConstS) Instantiate(names []string) C {
	return ConstC{c.Val}
}

func (c ConstS) String() string {
	return c.Val.String()
}

func (s ConstS) Values() []S {
	vals := make([]S, len(s.Val.Values()))
	for i, val := range s.Val.Values() {
		vals[i] = ConstS{val}
	}
	return vals
}

func (s ConstS) Head() TemplateId {
	return s.Val.Head()
}

func (c ConstC) Instantiate(ts []T) T {
	return c.Val
}

func (c ConstC) Uninstantiate(names []string) S {
	return ConstS{c.Val}
}

func (c ConstC) String() string {
	return c.Val.String()
}

func (c ConstC) Values() []C {
	vals := make([]C, len(c.Val.Values()))
	for i, val := range c.Val.Values() {
		vals[i] = ConstC{val}
	}
	return vals
}

func (c ConstC) Head() TemplateId {
	return c.Val.Head()
}
