package dynamics

import (
	"github.com/paulfchristiano/dwimmer/storage"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
)

type Dwimmer interface {
	Ask(term.T) (term.T, *term.SettingT)
	Answer(term.T) (term.T, term.T)
	Run(*term.SettingT) term.T
	Do(term.ActionT, *term.SettingT) term.T

	Stack
	Transitions
	ui.UIImplementer
	storage.StorageImplementer

	Close()
}

var DefaultInitializers []term.T

func AddInitializer(t term.T) {
	DefaultInitializers = append(DefaultInitializers, t)
}

var SetupState = term.Make("initialize the interpreter's state")

func init() {
	AddInitializer(SetupState.T())
}
