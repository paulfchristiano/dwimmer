package dynamics

import "github.com/paulfchristiano/dwimmer/term"

type Dwimmer interface {
	Ask(term.T) (term.T, *term.SettingT)
	Answer(term.T) (term.T, term.T)
	Run(*term.SettingT) term.T
	Do(term.ActionT, *term.SettingT) term.T
	Continuations(term.SettingId) []term.TemplateId
	Save(term.SettingId, Transition)
	Set(term.SettingId, Transition)
	Get(term.SettingId) (Transition, bool)
	Transitions() *TransitionTable
	Writeln(string)
	Readln(string, ...[]string) string
}
