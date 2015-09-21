package dynamics

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/storage/database"
	"github.com/paulfchristiano/dwimmer/term"
)

type Transition interface {
	Step(Dwimmer, *term.SettingT) term.T
}

type NativeTransition func(Dwimmer, *term.SettingT) term.T

func (t NativeTransition) Step(d Dwimmer, s *term.SettingT) term.T {
	return t(d, s)
}

type SimpleTransition struct {
	Action term.ActionC
}

func (t SimpleTransition) Step(d Dwimmer, s *term.SettingT) term.T {
	actC := t.Action
	s.AppendAction(actC)
	actT := actC.Instantiate(s.Args)
	return d.Do(actT, s)
}

type TransitionTable struct {
	collection    database.C
	table         map[term.SettingID]Transition
	continuations map[term.SettingID]([]*term.Setting)
}

type Transitions interface {
	Save(*term.Setting, Transition)
	Set(*term.Setting, Transition)
	Get(*term.Setting) (Transition, bool)
	Continuations(*term.Setting) []*term.Setting
}

var DefaultTransitions = NewTransitionTable(database.Collection("newtransitions"))

func NewTransitionTable(C database.C) *TransitionTable {
	result := &TransitionTable{
		collection:    C,
		table:         make(map[term.SettingID]Transition),
		continuations: make(map[term.SettingID]([]*term.Setting)),
	}
	for _, transition := range C.All() {
		settingRecord, ok := transition["key"]
		if !ok {
			settingRecord, ok = transition["setting"]
		}
		if !ok {
			continue
		}
		actionRecord, ok := transition["value"]
		if !ok {
			actionRecord, ok = transition["action"]
		}
		if !ok {
			continue
		}
		setting, err := term.LoadSetting(settingRecord)
		if err != nil {
			continue
		}
		action, err := term.LoadActionC(actionRecord)
		if err != nil {
			continue
		}
		result.SetSimpleC(setting, action)
	}
	return result
}

func (table *TransitionTable) AddContinuation(s *term.Setting) {
	if s.Size == 0 {
		return
	}
	continuations := table.continuations[s.Previous.ID]
	for _, c := range continuations {
		if c.Last.LineID() == s.Last.LineID() {
			return
		}
	}
	table.continuations[s.Previous.ID] = append(continuations, s)
	if s.Size > 0 {
		table.AddContinuation(s.Previous)
	}
}

func (table *TransitionTable) Set(s *term.Setting, t Transition) {
	table.AddContinuation(s)
	table.table[s.ID] = t
}

func (table *TransitionTable) Save(s *term.Setting, t Transition) {
	table.Set(s, t)
	switch t := t.(type) {
	case SimpleTransition:
		table.collection.Set(term.SaveSetting(s), term.SaveActionC(t.Action))
	}
}

func (t *TransitionTable) SetSimpleC(s *term.Setting, a term.ActionC) {
	t.Set(s, SimpleTransition{a})
}

func (t *TransitionTable) SaveSimpleC(s *term.Setting, a term.ActionC) {
	t.Save(s, SimpleTransition{a})
}

func (t *TransitionTable) SaveSimpleS(s *term.SettingS, a term.ActionS) *term.SettingS {
	aC := a.Instantiate(s.Names)
	t.SaveSimpleC(s.Setting, aC)
	return s.Copy().AppendAction(aC)
}

func (t *TransitionTable) Get(s *term.Setting) (Transition, bool) {
	result, ok := t.table[s.ID]
	return result, ok
}

var allNames = []string{"x", "y", "z", "w", "i", "j", "k",
	"b", "c", "d", "e", "f", "g", "h",
	"l", "m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v"}

func AddSimple(s *term.SettingS, a term.ActionS) *term.SettingS {
	return DefaultTransitions.SaveSimpleS(s, a)
}

func AddNativeResponse(t term.TemplateID, n int, f func(Dwimmer, *term.SettingT, ...term.T) term.T) {
	names := append([]string{"_Q"}, allNames[:n]...)
	s := ExpectQuestion(term.InitS(), t, names...)
	AddNative(s, f, allNames[:n]...)
}

func AddNative(s *term.SettingS, f func(Dwimmer, *term.SettingT, ...term.T) term.T, names ...string) {
	AddNativeTo(DefaultTransitions, s, f, names...)
}

var (
	NativeQ = term.Make("what does the appropriate native function return in this setting?")
)

func AddNativeTo(table *TransitionTable, s *term.SettingS,
	f func(Dwimmer, *term.SettingT, ...term.T) term.T, names ...string) {

	indices := make([]int, len(names))
	for i, name := range names {
		for j, key := range s.Names {
			if name == key {
				indices[i] = j
			}
		}
	}
	g := func(d Dwimmer, s *term.SettingT) term.T {
		args := make([]term.T, len(indices))
		for i, index := range indices {
			args[i] = s.Args[index]
		}
		result := f(d, s, args...)
		if result != nil {
			s.AppendAction(term.AskC(NativeQ.C()))
			s.AppendTerm(core.Answer.T(result))
			s.AppendAction(term.ReturnC(term.Cr(len(s.Args) - 1)))
		}
		return result
	}
	table.Save(s.Setting, NativeTransition(g))
}

func (table *TransitionTable) Continuations(s *term.Setting) []*term.Setting {
	return table.continuations[s.ID]
}
