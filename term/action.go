package term

import (
	"bytes"
	"fmt"
)

type Action int

const (
	Return = Action(iota)
	Ask
	View
	Replace
	Clarify
	Correct
	Delete
	Meta
	Replay
)

func (a Action) String() string {
	switch a {
	case Return:
		return "reply"
	case Ask:
		return "ask"
	case View:
		return "view"
	case Replace:
		return "replace"
	case Replay:
		return "replay"
	case Clarify:
		return "ask@"
	case Correct:
		return "correct"
	case Delete:
		return "del"
	case Meta:
		return "meta"
	}
	panic("unknown type of action")
}

type ActionT struct {
	Act     Action
	Args    []T
	IntArgs []int
}

type ActionC struct {
	Act     Action
	Args    []C
	IntArgs []int
}

type ActionS struct {
	Act     Action
	Args    []S
	IntArgs []int
}

func (a ActionC) IsValid() bool {
	var allowed [][2]int
	switch a.Act {
	case Return, Ask, View:
		allowed = [][2]int{{1, 0}}
	case Replace:
		allowed = [][2]int{{1, 0}, {1, 1}}
	case Replay, Correct, Delete:
		allowed = [][2]int{{0, 1}}
	case Meta:
		allowed = [][2]int{{0, 0}}
	case Clarify:
		allowed = [][2]int{{2, 0}}
	default:
		panic("unknown type of action")
	}
	actual := [2]int{len(a.Args), len(a.IntArgs)}
	for _, possibility := range allowed {
		if actual == possibility {
			return true
		}
	}
	return false
}

func (a ActionC) Instantiate(ts []T) ActionT {
	args := make([]T, len(a.Args))
	for i, arg := range a.Args {
		args[i] = arg.Instantiate(ts)
	}
	return ActionT{a.Act, args, a.IntArgs}
}

func (a ActionC) Uninstantiate(names []string) ActionS {
	args := make([]S, len(a.Args))
	for i, arg := range a.Args {
		args[i] = arg.Uninstantiate(names)
	}
	return ActionS{a.Act, args, a.IntArgs}
}

func (a ActionS) Instantiate(names []string) ActionC {
	args := make([]C, len(a.Args))
	for i, arg := range a.Args {
		args[i] = arg.Instantiate(names)
	}
	return ActionC{a.Act, args, a.IntArgs}
}

func (a ActionCID) String() string {
	return a.ActionC().String()
}

func (a ActionC) String() string {
	b := new(bytes.Buffer)
	b.WriteString(a.Act.String())
	for _, arg := range a.IntArgs {
		b.WriteString(fmt.Sprintf(" %d", arg))
	}
	for _, arg := range a.Args {
		b.WriteString(fmt.Sprintf(" %v", arg))
	}
	return b.String()
}

func (a ActionS) String() string {
	b := new(bytes.Buffer)
	b.WriteString(a.Act.String())
	for _, arg := range a.IntArgs {
		b.WriteString(fmt.Sprintf(" %d", arg))
	}
	for _, arg := range a.Args {
		b.WriteString(fmt.Sprintf(" %v", arg))
	}
	return b.String()
}

func ReturnS(s S) ActionS {
	return ActionS{Act: Return, Args: []S{s}}
}
func AskS(s S) ActionS {
	return ActionS{Act: Ask, Args: []S{s}}
}
func ViewS(s S) ActionS {
	return ActionS{Act: View, Args: []S{s}}
}
func ReplaceS(s S, n int) ActionS {
	return ActionS{Replace, []S{s}, []int{n}}
}
func ClarifyS(s, t S) ActionS {
	return ActionS{Clarify, []S{s, t}, []int{}}
}
func CorrectS(n int) ActionS {
	return ActionS{Correct, []S{}, []int{n}}
}
func DeleteS(n int) ActionS {
	return ActionS{Delete, []S{}, []int{n}}
}
func ReturnC(c C) ActionC {
	return ActionC{Return, []C{c}, []int{}}
}
func AskC(c C) ActionC {
	return ActionC{Ask, []C{c}, []int{}}
}
func ViewC(c C) ActionC {
	return ActionC{View, []C{c}, []int{}}
}
func ReplaceC(c C, n int) ActionC {
	return ActionC{Replace, []C{c}, []int{n}}
}
func ClarifyC(c, d C) ActionC {
	return ActionC{Clarify, []C{c, d}, []int{}}
}
func CorrectC(n int) ActionC {
	return ActionC{Correct, []C{}, []int{n}}
}
func DeleteC(n int) ActionC {
	return ActionC{Delete, []C{}, []int{n}}
}
func ReturnT(t T) ActionT {
	return ActionT{Return, []T{t}, []int{}}
}
func AskT(t T) ActionT {
	return ActionT{Ask, []T{t}, []int{}}
}
func ViewT(t T) ActionT {
	return ActionT{View, []T{t}, []int{}}
}
func ReplaceT(t T, n int) ActionT {
	return ActionT{Replace, []T{t}, []int{n}}
}
func ClarifyT(t, u T) ActionT {
	return ActionT{Clarify, []T{t, u}, []int{}}
}
func CorrectT(n int) ActionT {
	return ActionT{Correct, []T{}, []int{n}}
}
func DeleteT(n int) ActionT {
	return ActionT{Delete, []T{}, []int{n}}
}

func MetaC() ActionC {
	return ActionC{Meta, []C{}, []int{}}
}
func ReplayC(n int) ActionC {
	return ActionC{Replay, []C{}, []int{n}}
}
