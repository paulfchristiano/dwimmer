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

	Transitions
	ui.UIImplementer
	storage.StorageImplementer

	Close()
}
