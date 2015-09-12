package represent

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/lists"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	QuotedNil = term.Make("nothing")

	QuotedSetting  = term.Make("the setting with the list of lines []")
	QuotedSettingT = term.Make("the setting with the list of lines [], " +
		"list of arguments [], and list of subordinate settings []")
	QuotedTemplate = term.Make("the term template that has parts []")

	QuotedReturn  = term.Make("the action that returns the term constructed by instantiating []")
	QuotedView    = term.Make("the action that views the term constructed by instantiating []")
	QuotedAsk     = term.Make("the action that asks the question constructed by instantiating []")
	QuotedReplace = term.Make("the action that replaces output [] with the term constructed " +
		"by instantiating []")
	QuotedClarify = term.Make("the action that asks the subsetting that produced output [] " +
		"the clarifying question constructed by instantiating []")
	QuotedCorrect = term.Make("the action that prompts the user to correct action []")
	QuotedClose   = term.Make("the action that closes the channel in position []")
	QuotedDelete  = term.Make("the action that deletes the term referenced by the constructor []")

	QuotedCompoundC = term.Make("the constructor that has template [] and arguments formed by " +
		"instantiating each constructor in []")
	QuotedReferenceC = term.Make("the constructor that returns the []th argument of the environment " +
		"where it is instantiated")
	QuotedConstC = term.Make("the constructor that returns the term []")

	QuotedCompoundT = term.Make("the term with head [] and arguments []")
	QuotedIntT      = term.Make("a term wrapping the native integer []")
	QuotedStrT      = term.Make("a term wrapping the native string []")
	QuotedWrapperT  = term.Make("a term referring to the go object []")
	QuotedQuoteT    = term.Make("a term referring to the representation of []")

	QuotedSimpleTransition = term.Make("the transition that takes the action []")
	QuotedNativeTransition = term.Make("the transition that applies the Go function [] " +
		"to the Dwimmer object evolving the current setting, and a pointer to that setting")
)

func SettingT(s *term.SettingT) term.T {
	if s == nil {
		return QuotedNil.T()
	}
	args := make([]term.T, len(s.Args))
	children := make([]term.T, len(s.Children))
	for i, arg := range s.Args {
		args[i] = T(arg)
	}
	for i, child := range s.Children {
		children[i] = SettingT(child)
	}
	return QuotedSettingT.T(settingLines(s.Setting()), List(args), List(children))
}

func settingLines(s *term.Setting) term.T {
	l := make([]term.T, 0)
	for i, output := range s.Outputs {
		l = append(l, Template(output))
		if i < len(s.Inputs) {
			l = append(l, ActionC(s.Inputs[i].ActionC()))
		}
	}
	return List(l)
}

func Setting(s *term.Setting) term.T {
	return QuotedSetting.T(settingLines(s))
}

func Template(temp term.TemplateId) term.T {
	t := temp.Template()
	parts := make([]term.T, len(t.Parts))
	for i, part := range t.Parts {
		parts[i] = Str(part)
	}
	return QuotedTemplate.T(List(parts))
}

func ActionC(a term.ActionC) term.T {
	switch a.Act {
	case term.Return:
		return QuotedReturn.T(C(a.Args[0]))
	case term.View:
		return QuotedView.T(C(a.Args[0]))
	case term.Ask:
		return QuotedAsk.T(C(a.Args[0]))
	case term.Replace:
		return QuotedReplace.T(C(a.Args[0]), Int(a.IntArgs[0]))
	case term.Clarify:
		return QuotedClarify.T(C(a.Args[0]), Int(a.IntArgs[0]))
	case term.Correct:
		return QuotedCorrect.T(Int(a.IntArgs[0]))
	case term.Close:
		return QuotedClose.T(Int(a.IntArgs[0]))
	case term.Delete:
		return QuotedDelete.T(Int(a.IntArgs[0]))
	}
	panic("quoting an unknown type of action")
}

func Transition(t dynamics.Transition) term.T {
	switch t := t.(type) {
	case dynamics.SimpleTransition:
		return QuotedSimpleTransition.T(ActionC(t.Action))
	case dynamics.NativeTransition:
		return QuotedNativeTransition.T(term.Wrap(t))
	}
	panic("quoting an unknown type of transition")
}

func C(c term.C) term.T {
	switch c := c.(type) {
	case *term.CompoundC:
		args := make([]term.T, len(c.Values()))
		for i, arg := range c.Values() {
			args[i] = C(arg)
		}
		return QuotedCompoundC.T(Template(c.Head()), List(args))
	case term.ReferenceC:
		return QuotedReferenceC.T(term.Int(c.Index))
	case term.ConstC:
		return QuotedConstC.T(T(c.Val))
	}
	panic(fmt.Sprintf("trying to quote %v of unknown type", c))
}

func T(t term.T) term.T {
	return term.Quoted{t}
}

func List(ts []term.T) term.T {
	result := lists.Empty.T()
	for i := len(ts) - 1; i >= 0; i-- {
		result = lists.Cons.T(ts[i], result)
	}
	return result
}

func makeExplicit(d dynamics.Dwimmer, s *term.SettingT, quoted term.T) term.T {
	t := quoted.(term.Quoted).Value
	switch t := t.(type) {
	case *term.CompoundT:
		args := make([]term.T, len(t.Values()))
		for i, arg := range t.Values() {
			args[i] = T(arg)
		}
		return core.Answer.T(QuotedCompoundT.T(Template(t.Head()), List(args)))
	case term.Int:
		return core.Answer.T(QuotedIntT.T(t))
	case term.Str:
		return core.Answer.T(QuotedStrT.T(t))
	case term.Wrapper:
		return core.Answer.T(QuotedWrapperT.T(t))
	case term.Quoted:
		return core.Answer.T(QuotedQuoteT.T(t.Value))
	}
	panic("unknown type of term!")
}

var (
	Explicit = term.Make("what term is []? the result should be represented explicitly such " +
		"that its properties can be inspected")
)

func getQuotedHead(d dynamics.Dwimmer, s *term.SettingT, q term.T) term.T {
	t := q.(term.Quoted).Value.Head()
	return core.Answer.T(Template(t))
}

func init() {
	s := term.InitS()
	s.AppendTemplate(Explicit, "t")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("t")))
	s.AppendTemplate(term.Quoted{}.Head())
	dynamics.AddNative(s, dynamics.Args1(makeExplicit), "t")
}

func Ch(c rune) term.T {
	return term.Make(fmt.Sprintf("the character '%c'", c)).T()
}

func Str(s string) term.T {
	return term.Str(s)
}

func Int(n int) term.T {
	return term.Int(n)
}
