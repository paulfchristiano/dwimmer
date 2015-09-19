package dwimmer

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"

	"github.com/paulfchristiano/dwimmer/data/core"
	_ "github.com/paulfchristiano/dwimmer/data/ints"
	"github.com/paulfchristiano/dwimmer/data/represent"
	_ "github.com/paulfchristiano/dwimmer/data/strings"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/storage"
	_ "github.com/paulfchristiano/dwimmer/storage/store"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
)

var (
	logger *log.Logger
)

func init() {
	f, err := os.Create("dwimmer-log")
	if err != nil {
		panic("failed to create log file")
	}
	logger = log.New(f, "", log.Lshortfile|log.Ltime)
}

type Dwimmer struct {
	dynamics.Transitions
	ui.UIImplementer
	storage.StorageImplementer
}

func (d *Dwimmer) Close() {
	d.CloseUI()
	d.CloseStorage()
	RecoverStackError(recover())
}

func (d *Dwimmer) Readln(s string, hintss ...[]string) string {
	defer func() {
		e := recover()
		if e != nil {
			switch e {
			case "interrupted":
				panic(StackError{[]*term.Template{}, "interrupted during input"})
			default:
				panic(e)
			}
		}
	}()
	return d.UIImplementer.Readln(s, hintss...)
}

func RecoverStackError(e interface{}) {
	if e != nil {
		switch e := e.(type) {
		case StackError:
			fmt.Printf("%s, printing top 20 questions in stack...\n", e.message)
			fmt.Printf("(stack size is %d)\n", len(e.stack))
			for _, t := range e.Top(20) {
				fmt.Println(t)
			}
		default:
			panic(e)
		}
	}
}

func NewDwimmer(impls ...ui.UIImplementer) *Dwimmer {
	var impl ui.UIImplementer
	if len(impls) == 1 {
		impl = impls[0]
	} else {
		impl = &ui.Term{}
	}
	result := &Dwimmer{
		Transitions:        dynamics.DefaultTransitions,
		UIImplementer:      impl,
		StorageImplementer: storage.NewStorage("state"),
	}
	result.InitUI()
	defer func() {
		e := recover()
		if e != nil {
			result.Close()
			panic(e)
		}
	}()
	RunInitializers(result)
	return result
}

var Initialization = term.Make("initializing a new dwimmer")

func RunInitializers(d dynamics.Dwimmer) {
	s := term.InitT()
	s.AppendTerm(Initialization.T())
	for _, t := range dynamics.DefaultInitializers {
		c := term.InitT()
		subRun(d, t, s, c)
	}
}

func DoC(d dynamics.Dwimmer, a term.ActionC, s *term.SettingT) term.T {
	return d.Do(a.Instantiate(s.Args), s)
}

func subAsk(d dynamics.Dwimmer, Q term.T, parent *term.SettingT) (term.T, *term.SettingT) {
	defer propagateStackError(Q)
	stackCheck()
	child := term.InitT()
	child.AppendTerm(dynamics.ParentChannel.T(term.MakeChannel(parent)))
	child.AppendTerm(Q)
	return d.Run(child), child
}

func subRun(d dynamics.Dwimmer, Q term.T, parent, child *term.SettingT) term.T {
	defer propagateStackError(Q)
	stackCheck()
	child.AppendTerm(dynamics.ParentChannel.T(term.MakeChannel(parent)))
	child.AppendTerm(Q)
	A := d.Run(child)
	parent.AppendTerm(dynamics.OpenChannel.T(term.MakeChannel(child)))
	parent.AppendTerm(A)
	return A
}

func (d *Dwimmer) Do(a term.ActionT, s *term.SettingT) term.T {
	switch a.Act {
	case term.Return:
		return a.Args[0]
	case term.Ask:
		Q := a.Args[0]
		child := term.InitT()
		subRun(d, Q, s, child)
		return nil
	case term.View:
		value := a.Args[0]
		if value != nil {
			s.AppendTerm(value)
		} else {
			s.AppendTerm(Closed.T())
		}
		return nil
	case term.Replace:
		//value := a.Args[0]
		//n := a.IntArgs[0]
		//s.Rollback(n).AppendTerm(value)
		return nil
	case term.Replay:
		n := a.IntArgs[0]
		s.Rollback(n)
		return nil
	case term.Clarify:
		Q := a.Args[1]
		//TODO handle null pointers much better...
		//(e.g. one part of an expression may refer to a deleted variable)
		if a.Args[0] == nil {
			s.AppendTerm(Closed.T())
			return nil
		}
		var target *term.SettingT
		channel, err := represent.ToChannel(d, a.Args[0])
		if err == nil {
			target = channel.(term.Channel).Instantiate()
		} else {
			var othererr term.T
			target, othererr = represent.ToSettingT(d, a.Args[0])
			if othererr != nil {
				s.AppendTerm(NotAChannel.T(err))
				return nil
			}
		}
		subRun(d, Q, s, target)
		return nil
	case term.Correct:
		n := a.IntArgs[0]
		old := s.Setting.Rollback(n)
		action := ElicitAction(d, old, true)
		d.Save(old, dynamics.SimpleTransition{action})
		s.AppendAction(action)
		s.AppendTerm(core.OK.T())
		return nil
	case term.Delete:
		n := a.IntArgs[0]
		s.Args[n] = nil
		s.AppendTerm(core.OK.T())
		return nil
	case term.Meta:
		s.AppendTerm(CurrentSetting.T(represent.SettingT(s)))
		return nil
	}
	panic("unknown kind of action")
}

var (
	DeleteNonVar   = term.Make("only variables can be deleted")
	CurrentSetting = term.Make("the current setting is []")
	Interrupted    = term.Make("execution was interrupted in setting []")
)

func (d *Dwimmer) Run(setting *term.SettingT) term.T {
	for {
		goMeta := func() {
			shell := term.InitT()
			shell.AppendTerm(Interrupted.T(represent.SettingT(setting)))
			StartShell(d, shell)
		}
		char, sent := d.CheckCh()
		if sent {
			if char == 'q' {
				panic("interrupted")
			}
			if char == 's' {
				goMeta()
			} else {
				d.Clear()
				d.Debug("(Type [s] to interrupt execution and drop into a shell)")
			}
		}
		transition, ok := d.Get(setting.Setting)
		if !ok {
			defer func() {
				e := recover()
				if e != nil {
					if e == "meta" {
						goMeta()
					} else {
						panic(e)
					}
				}
			}()
			Q := FallThrough.T(represent.SettingT(setting))
			result, _ := subAsk(d, Q, setting)
			var err term.T
			switch result.Head() {
			case TakeTransition:
				transition, err = represent.ToTransition(d, result.Values()[0])
				if err == nil {
					ok = true
				}
			case core.OK:
			default:
				result, err = d.Answer(TransitionGiven.T(result, Q))
				if err == nil {
					transition, err = represent.ToTransition(d, result)
					if err == nil {
						ok = true
					}
				}
			}
		}
		if ok {
			result := transition.Step(d, setting)
			if result != nil {
				return result
			}
		}
	}
}

func stackSize() uint64 {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	return mem.StackInuse
}

type StackError struct {
	stack   []*term.Template
	message string
}

func (e *StackError) Add(t *term.Template) {
	e.stack = append(e.stack, t)
}

func (e *StackError) Top(n int) []*term.Template {
	if n > len(e.stack) {
		n = len(e.stack)
	}
	return e.stack[len(e.stack)-n:]
}

func (e *StackError) StackSize() int {
	return len(e.stack)
}

func propagateStackError(Q term.T) {
	x := recover()
	if x != nil {
		switch y := x.(type) {
		case StackError:
			y.Add(Q.Head().Template())
			panic(y)
		}
		panic(x)
	}
}

func stackCheck() {
	if rand.Int()%100 == 0 {
		if stackSize() > 5e8 {
			panic(StackError{[]*term.Template{}, "stack is too large!"})
		}
	}
}

func (d *Dwimmer) Ask(Q term.T) (term.T, *term.SettingT) {
	setting := term.InitT().AppendTerm(Q)
	return d.Run(setting), setting
}

var (
	GetAnswerQ      = term.Make("what is the answer to a question whose representation satisfies property []?")
	WhileAttempting = term.Make("while trying to figure out what answer to produce, received the reply []")
	Closed          = term.Make("an error signal returned when a deleted argument is viewed")
	NotAChannel     = term.Make("received [] while trying to convert argument to a channel")
	BuiltinAnswerer = term.Make("the builtin Answer function is trying to answer []")
	IsAnswer        = term.Make("if [] is given as a reply to question [], does it provide an answer? " +
		"the reply should be either 'yes' or 'no'")
	WhatAnswer      = term.Make("if [] is given as a reply to question [], what answer does it imply?")
	IsAnswerClarify = term.Make("that response should be repeated as 'yes' or 'no'")
)

func (d *Dwimmer) Answer(q term.T) (term.T, term.T) {
	s := term.InitT()
	s.AppendTerm(BuiltinAnswerer.T(q))
	a, _ := subAsk(d, q, s)
	switch a.Head() {
	case core.Answer:
		return a.Values()[0], nil
	case core.Yes, core.No:
		return a, nil
	}
	follow := term.InitT()
	isAnswer := subRun(d, IsAnswer.T(a, q), s, follow)
	for {
		switch isAnswer.Head() {
		case core.Yes:
			result, err := d.Answer(WhatAnswer.T(a, q))
			if err != nil {
				return nil, a
			}
			return result, nil
		case core.No:
			return nil, a
		}
		isAnswer = subRun(d, IsAnswerClarify.T(), s, follow)
	}
}
