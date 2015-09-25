package dynamics

import (
	"math/rand"
	"runtime"

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

func SubAsk(d Dwimmer, Q term.T, parent *term.SettingT) (term.T, *term.SettingT) {
	d.Push(Q)
	stackCheck()
	child := term.InitT()
	child.AppendTerm(ParentChannel.T(term.MakeChannel(parent)))
	child.AppendTerm(Q)
	A := d.Run(child)
	d.Pop()
	return A, child
}

func SubRun(d Dwimmer, Q term.T, parent, child *term.SettingT) term.T {
	d.Push(Q)
	stackCheck()
	child.AppendTerm(ParentChannel.T(term.MakeChannel(parent)))
	child.AppendTerm(Q)
	A := d.Run(child)
	parent.AppendTerm(OpenChannel.T(term.MakeChannel(child)))
	if A != nil {
		parent.AppendTerm(A)
	}
	d.Pop()
	return A
}

func stackCheck() {
	if rand.Int()%100 == 0 {
		if stackSize() > 5e8 {
			panic("stack is too large!")
		}
	}
}

func stackSize() uint64 {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	return mem.StackInuse
}
