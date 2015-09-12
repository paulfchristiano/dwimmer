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
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
)

type dwimmer struct {
	*dynamics.TransitionTable
}

func (d *dwimmer) Transitions() *dynamics.TransitionTable {
	return d.TransitionTable
}

func (d *dwimmer) Writeln(s string) {
	ui.Println(s)
}

func DisplayStackError() {
	e := recover()
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

func (d *dwimmer) Readln(s string, hintss ...[]string) string {
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
	var hints []string
	if len(hintss) == 0 {
		hints = []string{}
	} else {
		hints = hintss[0]
	}
	ui.Print(s)
	return ui.GetLine(hints)
}

func Dwimmer() *dwimmer {
	return &dwimmer{
		dynamics.DefaultTransitions,
	}
}

func (d *dwimmer) DoC(a term.ActionC, s *term.SettingT) term.T {
	return d.Do(a.Instantiate(s.Args), s)
}

func (d *dwimmer) Do(a term.ActionT, s *term.SettingT) term.T {
	switch a.Act {
	case term.Return:
		return a.Args[0]
	case term.Ask:
		question := a.Args[0]
		answer, setting := d.Ask(question)
		s.AppendTerm(answer)
		s.SetLastChild(setting)
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
		var numToRemove int
		if n == -1 {
			numToRemove = 1
		} else {
			numToRemove = len(s.Children) - a.IntArgs[0]
		}
		var saveChild *term.SettingT
		for i := 0; i < numToRemove; i++ {
			if i == numToRemove-1 {
				saveChild = s.LastChild()
			}
			s.RemoveAction()
			s.RemoveTerm()
		}
		s.AppendTerm(value)
		s.SetLastChild(saveChild)
		return nil
	case term.Clarify:
		question := a.Args[0]
		target := s.Children[a.IntArgs[0]]
		if target == nil {
			s.AppendTerm(term.Make("there is no one here").T())
			return nil
		}
		target.AppendTerm(question)
		result := d.Run(target)
		s.AppendTerm(result)
		return nil
	case term.Correct:
		n := a.IntArgs[0]
		oldSetting := s.Setting().RollBack(n)
		oldid := term.IdSetting(oldSetting)
		s.AppendTerm(ElicitCorrection.T())
		action := ElicitAction(d, oldid)
		d.Save(oldid, dynamics.SimpleTransition{action})
		s.AppendAction(action)
		s.AppendTerm(core.OK.T())
		return nil
	case term.Close:
		n := a.IntArgs[0]
		s.Children[n] = nil
		s.AppendTerm(core.OK.T())
		return nil
	case term.Delete:
		n := a.IntArgs[0]
		s.Args[n] = nil
		s.AppendTerm(core.OK.T())
		return nil
	}
	panic("unknown kind of action")
}

var (
	ElicitCorrection = term.Make("the next action you enter will be used as a correction")
)

func (d *dwimmer) Run(setting *term.SettingT) term.T {
	for {
		transition, ok := d.Get(setting.Head())
		if !ok {
			actQ := ActionQ.T(represent.Setting(setting.Setting()))
			actA, err := d.Answer(actQ)
			if err != nil {
				return WhileAttempting.T(err)
			}
			transition, err = represent.ToTransition(d, actA)
			if err != nil {
				return term.Make("received [] while trying to convert a term representing a transition " +
					"into a transition").T(err)
			}
			/*
				fmt.Println("have to go meta!")
				metaQ := meta.MakeMetaQ(setting)
				metaA, err := d.Answer(metaQ)
				if err != nil {
					return WhileAttempting.T(err)
				}
				result, err := represent.ToT(d, metaA)
				if err != nil {
					return term.Make("received [] while trying to convert the best reply to a native term").T(err)
				}
				return result
			*/
		}
		result := transition.Step(d, setting)
		if result != nil {
			return result
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

func (d *dwimmer) Ask(Q term.T) (term.T, *term.SettingT) {
	defer func() {
		x := recover()
		if x != nil {
			switch y := x.(type) {
			case StackError:
				y.Add(Q.Head().Template())
				panic(y)
			}
			panic(x)
		}
	}()
	if rand.Int()%100 == 0 {
		if stackSize() > 1e7 {
			panic(StackError{[]*term.Template{}, "stack is too large!"})
		}
	}
	setting := term.InitT().AppendTerm(Q)
	return d.Run(setting), setting
}

var (
	GetAnswerQ      = term.Make("what is the answer to a question whose representation satisfies property []?")
	WhileAttempting = term.Make("while trying to figure out what answer to produce, received the reply []")
	Closed          = term.Make("a standin term returned when a deleted argument is viewed")
)

func (d *dwimmer) Answer(q term.T) (term.T, term.T) {
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
