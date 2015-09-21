package functions

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/lists"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	Apply   = term.Make("what is the result of applying function [] to argument []?")
	Compose = term.Make("the function that applies [], then applies [] to the result")
	Map     = term.Make("the function that applies [] to each element of its argument, which should be a list")
	Fold    = term.Make("what is the result of using [] to combine all of the arguments " +
		"in the list [] into one, starting from [] and assuming that the function is associative?")

	Apply2 = term.Make("what is the result of applying function [] to the arguments [] and []?")
)

var applyState, applyState2 *term.SettingS

func ApplicationSetting() *term.SettingS {
	return applyState.Copy()
}

func ApplicationSetting2() *term.SettingS {
	return applyState2.Copy()
}

func init() {
	applyState = term.InitS()
	applyState = dynamics.ExpectQuestion(applyState, Apply, "Q1", "f", "x")
	applyState = dynamics.AddSimple(applyState, term.ViewS(term.Sr("f")))
	applyState2 = term.InitS()
	applyState2 = dynamics.ExpectQuestion(applyState2, Apply2, "f", "x", "y")
	applyState2 = dynamics.AddSimple(applyState2, term.ViewS(term.Sr("f")))
}

func init() {
	s := ApplicationSetting()
	s.AppendTemplate(Compose, "f1", "f2")
	s = dynamics.AddSimple(s, term.AskS(Apply.S(term.Sr("f1"), term.Sr("x"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A1", "y")
	s = dynamics.AddSimple(s, term.AskS(Apply.S(term.Sr("f2"), term.Sr("y"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A2", "z")
	s = dynamics.AddSimple(s, term.ReturnS(term.Sr("z")))
}

func init() {
	var s, t *term.SettingS
	s = ApplicationSetting()
	s.AppendTemplate(Map, "g")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("x")))

	t = s.Copy().AppendTemplate(lists.Empty)
	t = dynamics.AddSimple(t, term.ReturnS(core.Answer.S(lists.Empty.S())))

	t = s.Copy().AppendTemplate(lists.Singleton, "y")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("g"), term.Sr("y"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "newy")
	t = dynamics.AddSimple(t, term.ReturnS(lists.Singleton.S(term.Sr("y"))))

	t = s.Copy().AppendTemplate(lists.Cons, "head", "tail")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("g"), term.Sr("head"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "newhead")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("f"), term.Sr("tail"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "newtail")
	t = dynamics.AddSimple(t, term.ReturnS(lists.Cons.S(term.Sr("newhead"), term.Sr("newtail"))))

	t = s.Copy().AppendTemplate(lists.Snoc, "init", "last")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("g"), term.Sr("last"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "newlast")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("f"), term.Sr("init"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "newinit")
	t = dynamics.AddSimple(t, term.ReturnS(lists.Snoc.S(term.Sr("newinit"), term.Sr("newlast"))))

	t = s.Copy().AppendTemplate(lists.Concat, "a", "b")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("f"), term.Sr("a"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "newa")
	t = dynamics.AddSimple(t, term.AskS(Apply.S(term.Sr("f"), term.Sr("b"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "newb")
	t = dynamics.AddSimple(t, term.ReturnS(lists.Concat.S(term.Sr("newa"), term.Sr("newb"))))
}

func init() {
	var s, t *term.SettingS
	s = term.InitS()
	s.AppendTemplate(Fold, "f", "l", "x")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("l")))

	t = s.Copy().AppendTemplate(lists.Empty)
	t = dynamics.AddSimple(t, term.ReturnS(term.Sr("x")))

	t = s.Copy().AppendTemplate(lists.Cons, "head", "tail")
	t = dynamics.AddSimple(t, term.AskS(Apply2.S(term.Sr("f"), term.Sr("x"), term.Sr("head"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "y")
	t = dynamics.AddSimple(t, term.AskS(Fold.S(term.Sr("f"), term.Sr("tail"), term.Sr("y"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "A")
	t = dynamics.AddSimple(t, term.ReturnS(term.Sr("A")))

	t = s.Copy().AppendTemplate(lists.Snoc, "init", "last")
	t = dynamics.AddSimple(t, term.AskS(Fold.S(term.Sr("f"), term.Sr("init"), term.Sr("x"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "y")
	t = dynamics.AddSimple(t, term.AskS(Apply2.S(term.Sr("f"), term.Sr("y"), term.Sr("last"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "A")
	t = dynamics.AddSimple(t, term.ReturnS(term.Sr("A")))

	t = s.Copy().AppendTemplate(lists.Concat, "a", "b")
	t = dynamics.AddSimple(t, term.AskS(Fold.S(term.Sr("f"), term.Sr("a"), term.Sr("x"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "y")
	t = dynamics.AddSimple(t, term.AskS(Fold.S(term.Sr("f"), term.Sr("b"), term.Sr("y"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "A")
	t = dynamics.AddSimple(t, term.ReturnS(term.Sr("A")))

	t = s.Copy().AppendTemplate(lists.Singleton, "y")
	t = dynamics.AddSimple(t, term.AskS(Apply2.S(term.Sr("f"), term.Sr("x"), term.Sr("y"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A1", "A")
	t = dynamics.AddSimple(t, term.ReturnS(term.Sr("A")))
}
