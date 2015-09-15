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
	table         map[term.SettingId]Transition
	continuations map[term.SettingId]([]*term.Setting)
}

type Transitions interface {
	Save(term.SettingId, Transition)
	Set(term.SettingId, Transition)
	Get(term.SettingId) (Transition, bool)
	Continuations(term.SettingId) []*term.Setting
}

var DefaultTransitions = NewTransitionTable(database.Collection("newtransitions"))

func NewTransitionTable(C database.C) *TransitionTable {
	result := &TransitionTable{
		collection:    C,
		table:         make(map[term.SettingId]Transition),
		continuations: make(map[term.SettingId]([]*term.Setting)),
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
		setting := term.LoadSetting(settingRecord)
		action := term.LoadActionC(actionRecord)
		result.SetSimpleC(term.IdSetting(setting), action)
	}
	return result
}

func (table *TransitionTable) AddContinuation(s *term.Setting) {
	if s.Size == 0 {
		return
	}
	continuations, ok := table.continuations[s.Previous.Id]
	if !ok {
		continuations = make([]*term.Setting, 0)
		table.continuations[s.Previous.Id] = continuations
	}
	for _, c := range continuations {
		if c.Last.LineId() == s.Last.LineId() {
			return
		}
	}
	continuations = append(continuations, s)
	if s.Size > 0 {
		table.AddContinuation(s.Previous)
	}
}

func (table *TransitionTable) Set(s term.SettingId, t Transition) {
	table.AddContinuation(s.Setting().Rollback(-1))
	table.table[s] = t
}

func (table *TransitionTable) Save(s term.SettingId, t Transition) {
	table.Set(s, t)
	switch t := t.(type) {
	case SimpleTransition:
		table.collection.Set(term.SaveSetting(s.Setting()), term.SaveActionC(t.Action))
	}
}

func (t *TransitionTable) SetSimpleC(s term.SettingId, a term.ActionC) {
	t.Set(s, SimpleTransition{a})
}

func (t *TransitionTable) SaveSimpleC(s term.SettingId, a term.ActionC) {
	t.Save(s, SimpleTransition{a})
}

func (t *TransitionTable) SaveSimpleS(s *term.SettingS, a term.ActionS) *term.SettingS {
	aC := a.Instantiate(s.Names)
	t.SaveSimpleC(s.Setting.Id, aC)
	return s.Copy().AppendAction(aC)
}

func (t *TransitionTable) Get(s term.SettingId) (Transition, bool) {
	result, ok := t.table[s]
	return result, ok
}

var allNames = []string{"x", "y", "z", "w", "i", "j", "k",
	"b", "c", "d", "e", "f", "g", "h",
	"l", "m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v"}

func AddSimple(s *term.SettingS, a term.ActionS) *term.SettingS {
	return DefaultTransitions.SaveSimpleS(s, a)
}

func AddNativeResponse(t term.TemplateId, n int, f func(Dwimmer, *term.SettingT, ...term.T) term.T) {
	names := allNames[:n]
	AddNative(term.InitS().AppendTemplate(t, names...), f, names...)
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
	table.Save(s.Setting.Id, NativeTransition(g))
}

func (table *TransitionTable) Continuations(id term.SettingId) []*term.Setting {
	return table.continuations[id]
}
