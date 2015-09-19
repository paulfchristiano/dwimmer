package term

import (
	"fmt"
	"reflect"

	"github.com/paulfchristiano/dwimmer/term/intern"
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

func MakeChannel(s *SettingT) Channel {
	return Channel{s.Copy()}
}

func (c Channel) Head() TemplateID {
	return channelHead
}

func (c Channel) String(ider intern.Packer) string {
	return fmt.Sprintf("->")
}

func (c Channel) Values() []T {
	return []T{}
}

func (c Channel) Instantiate() *SettingT {
	return c.Setting.Copy()
}

//Strings

func (s Str) Head() TemplateID {
	return strHead
}

func (s Str) String(ider intern.Packer) string {
	return fmt.Sprintf("\"%s\"", string(s))
}

func (s Str) Values() []T {
	return make([]T, 0)
}

//Integers

func (n Int) Head() TemplateID {
	return intHead
}

func (n Int) String(ider intern.Packer) string {
	return fmt.Sprintf("%d", int(n))
}

func (n Int) Values() []T {
	return make([]T, 0)
}

//Terms

func Quote(v T) T {
	return Quoted{v}
}

func (q Quoted) Head() TemplateID {
	return quotedHead
}

func (q Quoted) String(ider intern.Packer) string {
	return fmt.Sprintf("T(%s)", q.Value.String(ider))
}

func (q Quoted) Values() []T {
	return make([]T, 0)
}

//Wrapper

func Wrap(i interface{}) T {
	return Wrapper{reflect.ValueOf(i)}
}

func (w Wrapper) Head() TemplateID {
	return wrapperHead
}

func (w Wrapper) String(ider intern.Packer) string {
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

func (c ConstS) String(ider intern.Packer) string {
	return c.Val.String(ider)
}

func (s ConstS) Values() []S {
	vals := make([]S, len(s.Val.Values()))
	for i, val := range s.Val.Values() {
		vals[i] = ConstS{val}
	}
	return vals
}

func (s ConstS) Head() TemplateID {
	return s.Val.Head()
}

func (c ConstC) Instantiate(ts []T) T {
	return c.Val
}

func (c ConstC) Uninstantiate(names []string) S {
	return ConstS{c.Val}
}

func (c ConstC) String(ider intern.Packer) string {
	return c.Val.String(ider)
}

func (c ConstC) Values() []C {
	vals := make([]C, len(c.Val.Values()))
	for i, val := range c.Val.Values() {
		vals[i] = ConstC{val}
	}
	return vals
}

func (c ConstC) Head() TemplateID {
	return c.Val.Head()
}
