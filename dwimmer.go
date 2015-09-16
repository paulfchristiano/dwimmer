package dwimmer

import (
	"fmt"
	"math/rand"
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
	return result
}

func (d *Dwimmer) DoC(a term.ActionC, s *term.SettingT) term.T {
	return d.Do(a.Instantiate(s.Args), s)
}

func subRun(d dynamics.Dwimmer, Q term.T, parent, child *term.SettingT) {
	defer propagateStackError(Q)
	stackCheck()
	child.AppendTerm(ParentChannel.T(term.MakeChannel(parent)))
	child.AppendTerm(Q)
	A := d.Run(child)
	parent.AppendTerm(OpenChannel.T(term.MakeChannel(child)))
	parent.AppendTerm(A)
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
		value := a.Args[0]
		n := a.IntArgs[0]
		s.Rollback(n).AppendTerm(value)
		return nil
	case term.Clarify:
		Q := a.Args[1]
		//TODO handle null pointers much better...
		//(e.g. one part of an expression may refer to a deleted variable)
		if a.Args[0] == nil {
			s.AppendTerm(Closed.T())
			return nil
		}
		channel, ok := a.Args[0].(term.Channel)
		if !ok {
			s.AppendTerm(NotAChannel.T())
			return nil
		}
		target := channel.Instantiate()
		subRun(d, Q, s, target)
		return nil
	case term.Correct:
		n := a.IntArgs[0]
		oldSetting := s.Setting.Rollback(n)
		oldid := term.IdSetting(oldSetting)
		action := ElicitAction(d, oldid, true)
		d.Save(oldid, dynamics.SimpleTransition{action})
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
)

func (d *Dwimmer) Run(setting *term.SettingT) term.T {
	for {
		transition, ok := d.Get(setting.Setting.Id)
		if !ok {
			//TODO have this method look at advice and determine what to do
			//(probably need to convert to one of the desired forms)
			d.Ask(FallThrough.T(represent.SettingT(setting)))
		} else {
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
		if stackSize() > 1e7 {
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
	OpenChannel     = term.Make("@[]")
	ParentChannel   = term.Make("@[]*")
	NotAChannel     = term.Make("an error signal returned when a message is sent to an invalid target")
)

func (d *Dwimmer) Answer(q term.T) (term.T, term.T) {
	a, _ := d.Ask(q)
	switch a.Head() {
	case core.Answer.Head():
		return a.Values()[0], nil
	case core.NoAnswer.Head():
		return nil, a
	}
	answer, err := d.Answer(GetAnswerQ.T(a))
	if err != nil {
		return nil, err
	}
	return answer, nil
}
