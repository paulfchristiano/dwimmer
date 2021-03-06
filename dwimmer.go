package dwimmer

import (
	"math/rand"

	"github.com/paulfchristiano/dwimmer/data/core"
	_ "github.com/paulfchristiano/dwimmer/data/ints"
	"github.com/paulfchristiano/dwimmer/data/represent"
	_ "github.com/paulfchristiano/dwimmer/data/strings"
	"github.com/paulfchristiano/dwimmer/dynamics"
	_ "github.com/paulfchristiano/dwimmer/dynamics/meta"
	"github.com/paulfchristiano/dwimmer/storage"
	_ "github.com/paulfchristiano/dwimmer/storage/store"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
)

type Dwimmer struct {
	dynamics.Transitions
	dynamics.Stack
	ui.UIImplementer
	storage.StorageImplementer
	lastWatcher *term.SettingT
}

func (d *Dwimmer) Close() {
	d.CloseUI()
	d.ShowStack()
	d.CloseStorage()
}

func TestDwimmer() *Dwimmer {
	result := &Dwimmer{
		Transitions:        dynamics.DefaultTransitions,
		UIImplementer:      ui.Dummy(),
		StorageImplementer: storage.Dummy(),
		Stack:              &dynamics.BasicStack{},
	}
	return result
}

func NewDwimmer(stateName string, impls ...ui.UIImplementer) *Dwimmer {
	var impl ui.UIImplementer
	if len(impls) == 1 {
		impl = impls[0]
	} else {
		impl = &ui.Term{}
	}
	var store storage.StorageImplementer
	if stateName == "" {
		store = storage.Dummy()
	} else {
		store = storage.NewStorage(stateName)
	}
	result := &Dwimmer{
		Transitions:        dynamics.DefaultTransitions,
		UIImplementer:      impl,
		StorageImplementer: store,
		Stack:              &dynamics.BasicStack{},
	}
	result.InitUI()
	defer func() {
		e := recover()
		if e != nil {
			result.Close()
			panic(e)
		}
	}()
	RunDefaultInitializers(result)
	return result
}

var Initialization = term.Make("initializing a new dwimmer")

func RunDefaultInitializers(d dynamics.Dwimmer) {
	s := term.InitT()
	s.AppendTerm(Initialization.T())
	for _, t := range dynamics.DefaultInitializers {
		c := term.InitT()
		dynamics.SubRun(d, t, s, c)
	}
}

func DoC(d dynamics.Dwimmer, a term.ActionC, s *term.SettingT) term.T {
	return d.Do(a.Instantiate(s.Args), s)
}

func addToStack(d dynamics.Dwimmer, t term.T) {

}

func (d *Dwimmer) Do(a term.ActionT, s *term.SettingT) term.T {
	switch a.Act {
	case term.Return:
		return a.Args[0]
	case term.Ask:
		Q := a.Args[0]
		child := term.InitT()
		dynamics.SubRun(d, Q, s, child)
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
		dynamics.SubRun(d, Q, s, target)
		return nil
	case term.Correct:
		n := a.IntArgs[0]
		old := s.Setting.Rollback(n)
		transition := ElicitAction(d, term.InitT(), old)
		d.Save(old, transition)
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
	watchFrequency      = 1000
	watcherDepth   uint = 0
)

func (d *Dwimmer) watchdog(setting *term.SettingT) {
	if rand.Int()%(watchFrequency<<(3*watcherDepth)) == 0 {
		watcherDepth++
		defer func() { watcherDepth-- }()
		parent := setting.Copy().AppendAction(term.MetaC())
		oldWatcher := d.lastWatcher
		newWatcher := term.InitT()
		var Q term.T
		if d.lastWatcher == nil {
			Q = IsAllWell.T(term.MakeChannel(parent))
		} else {
			Q = IsAllWellPred.T(term.MakeChannel(parent), term.MakeChannel(oldWatcher))
		}
		d.lastWatcher = newWatcher
		dynamics.SubRun(d, Q, parent, newWatcher)
	}
}

var (
	IsAllWell     = term.Make("is everything OK in the state accessible over the channel []?")
	IsAllWellPred = term.Make("is everything OK in the state accessible over the channel []? " +
		"the last check was implemented by []")
	DeleteNonVar   = term.Make("only variables can be deleted")
	CurrentSetting = term.Make("the current setting is []")
	Interrupted    = term.Make("execution was interrupted in setting []")
)

func (d *Dwimmer) Run(setting *term.SettingT) term.T {
	goMeta := func() {
		StartShell(d, Interrupted.T(represent.SettingT(setting)))
	}
	for {
		//d.watchdog(setting)
		char, sent := d.CheckCh()
		if sent {
			if char == 'q' {
				panic("interrupted")
			} else if char == 's' {
				goMeta()
			} else {
				d.Clear()
				d.Debug("(Type [s] to interrupt execution and drop into a shell)")
				d.Debug("(Type [q] to interrupt execution and quit)")
			}
		}
		transition, ok := d.Get(setting.Setting)
		if !ok {
			transition = ElicitAction(d, setting.Copy(), setting.Setting)
		}
		result := transition.Step(d, setting)
		if result != nil {
			return result
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

func (d *Dwimmer) Answer(q term.T, optionalSetting ...*term.SettingT) (term.T, term.T) {
	var s *term.SettingT
	if len(optionalSetting) == 1 {
		s = optionalSetting[0]
	} else {
		s = term.InitT()
		s.AppendTerm(BuiltinAnswerer.T(q))
	}
	a := dynamics.SubRun(d, q, s)
	switch a.Head() {
	case core.Answer:
		return a.Values()[0], nil
	case core.Yes, core.No:
		return a, nil
	}
	follow := term.InitT()
	isAnswer := dynamics.SubRun(d, IsAnswer.T(a, q), s, follow)
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
		isAnswer = dynamics.SubRun(d, IsAnswerClarify.T(), s, follow)
	}
}
