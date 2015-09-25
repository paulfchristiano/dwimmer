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
	QuotedSettingT = term.Make("the concrete setting with template [] " +
		"and list of arguments []")
	QuotedTemplate = term.Make("the term template that has parts []")

	ActionLookup = map[term.Action]term.TemplateID{
		term.Return: term.Make("the parametrized action that returns the instantiation of its first argument"),
		term.View:   term.Make("the parametrized action that views the instantiation of its first argument"),
		term.Ask:    term.Make("the parametrized action that asks the instantiation of its first argument"),
		term.Replace: term.Make("the parametrized action that replaces the line given by its first index with " +
			"with the instantiation of its first argument"),
		term.Replay: term.Make("the parametrized action that replays the line given by its first index"),
		term.Clarify: term.Make("the parametrized action that sends the instantiation of its first argument " +
			"to the instantiation of its second argument"),
		term.Correct: term.Make("the parametrized action that prompts the user to correct its first index"),
		term.Delete:  term.Make("the action that deletes the variable indexed by its first index"),
		term.Meta:    term.Make("the action that returns a reference to the current setting"),
	}

	QuotedCompoundC = term.Make("the constructor that has template [] and arguments formed by " +
		"instantiating each constructor in []")
	QuotedReferenceC = term.Make("the constructor that returns the argument that comes after [] " +
		"others in the environment where it is instantiated")
	QuotedConstC = term.Make("the constructor that returns the term []")

	QuotedCompoundT = term.Make("the term with head [] and arguments []")
	QuotedIntT      = term.Make("a term wrapping the native integer []")
	QuotedStrT      = term.Make("a term wrapping the native string []")
	QuotedWrapperT  = term.Make("a term referring to the go object []")
	QuotedQuoteT    = term.Make("a term referring to the representation of []")

	QuotedChannel = term.Make("a channel for sending messages to the setting []")

	QuotedSimpleTransition = term.Make("the transition that takes the action []")
	QuotedNativeTransition = term.Make("the transition that applies the Go function [] " +
		"to the Dwimmer object evolving the current setting, and a pointer to that setting")

	QuotedActionC = term.Make("the action that performs [] with arguments [] and indices []")
)

func SettingT(s *term.SettingT) term.T {
	if s == nil {
		return QuotedNil.T()
	}
	args := make([]term.T, len(s.Args))
	for i, arg := range s.Args {
		args[i] = T(arg)
	}
	return QuotedSettingT.T(Setting(s.Setting), List(args))
}

func SettingLine(l term.SettingLine) term.T {
	switch l := l.(type) {
	case term.TemplateID:
		return Template(l)
	case term.ActionCID:
		return ActionC(l.ActionC())
	default:
		panic("quoting unknown type of setting line!")
	}
}

func Setting(s *term.Setting) term.T {
	lines := s.Lines()
	quotedLines := make([]term.T, len(lines))
	for i, line := range lines {
		quotedLines[i] = SettingLine(line)
	}
	return QuotedSetting.T(List(quotedLines))
}

func Template(temp term.TemplateID) term.T {
	t := temp.Template()
	parts := make([]term.T, len(t.Parts))
	for i, part := range t.Parts {
		parts[i] = Str(part)
	}
	return QuotedTemplate.T(List(parts))
}

func ActionC(a term.ActionC) term.T {
	args := make([]term.T, len(a.Args))
	for i, arg := range a.Args {
		args[i] = C(arg)
	}
	intargs := make([]term.T, len(a.IntArgs))
	for i, intarg := range a.IntArgs {
		intargs[i] = Int(intarg)
	}
	return QuotedActionC.T(Action(a.Act), List(args), List(intargs))
}

func Action(a term.Action) term.T {
	result, ok := ActionLookup[a]
	if !ok {
		panic("quoting unknown action")
	}
	return result.T()
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
		return QuotedReferenceC.T(Int(c.Index))
	case term.ConstC:
		return QuotedConstC.T(T(c.Val))
	}
	panic(fmt.Sprintf("trying to quote %v of unknown type %T", c, c))
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
	s = dynamics.ExpectQuestion(s, Explicit, "Q", "t")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("t")))
	s.AppendTemplate(term.Quoted{}.Head())
	dynamics.AddNative(s, dynamics.Args1(makeExplicit), "t")
}

func Ch(c rune) term.T {
	return term.Make(fmt.Sprintf("the character '%c'", c)).T()
}

func Str(s string) term.T {
	return term.Str(s)
	/*
		runes := make([]term.T, 0)
		for _, c := range s {
			runes = append(runes, strings.Rune.T(Int(int(c))))
		}
		return strings.ByRunes.T(List(runes))
	*/
}

func Int(n int) term.T {
	return term.Int(n)
	/*
		switch {
		case n == 0:
			return ints.Zero.T()
		case n == 1:
			return ints.One.T()
		case n < 0:
			return ints.Negative.T(Int(-n))
		case n%2 == 0:
			return ints.Double.T(Int(n / 2))
		case n%2 == 1:
			return ints.DoublePlusOne.T(Int(n / 2))
		default:
			panic("unreachable")
		}
	*/
}
