package represent

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	GetTemplate             = term.Make("what is the template of term []?")
	GetArguments            = term.Make("what are the arguments of term []?")
	GetTemplateAndArguments = term.Make("what are the template and arguments of term []?")
	GetSetting              = term.Make("what is the abstract setting associated with the concrete setting []?")
	TemplateAndArguments    = term.Make("the template is [] and the arguments are []")
	NumArguments            = term.Make("how many arguments does the term [] have?")
)

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, GetSetting, "Q", "s")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("s")))
	s.AppendTemplate(QuotedSettingT, "lines", "arguments", "children")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(QuotedSetting.S(term.Sr("lines")))))
}

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, GetTemplateAndArguments, "Q", "t")
	s = dynamics.AddSimple(s, term.AskS(GetTemplate.S(term.Sr("t"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A", "template")
	s = dynamics.AddSimple(s, term.AskS(GetArguments.S(term.Sr("t"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A1", "arguments")
	s = dynamics.AddSimple(s, term.ReturnS(TemplateAndArguments.S(term.Sr("template"), term.Sr("arguments"))))
}

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, GetTemplate, "Q", "t")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("t")))

	t := s.Copy().AppendTemplate(term.Quoted{}.Head())
	t = dynamics.AddSimple(t, term.AskS(Explicit.S(term.Sr("t"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A", "explicit")
	t = dynamics.AddSimple(t, term.AskS(GetTemplate.S(term.Sr("explicit"))))
	t = dynamics.ExpectAnswer(t, core.Answer, "A2", "result")
	t = dynamics.AddSimple(t, term.ReturnS(core.Answer.S(term.Sr("result"))))

	t = s.Copy().AppendTemplate(QuotedQuoteT, "q")
	dynamics.AddSimple(t, term.ReturnS(core.Answer.S(
		term.Sc(Template(term.Quoted{}.Head())),
	)))

	t = s.Copy().AppendTemplate(QuotedIntT, "q")
	dynamics.AddSimple(t, term.ReturnS(core.Answer.S(
		term.Sc(Template(term.Int(0).Head())),
	)))

	t = s.Copy().AppendTemplate(QuotedStrT, "q")
	dynamics.AddSimple(t, term.ReturnS(core.Answer.S(
		term.Sc(Template(term.Str(0).Head())),
	)))

	t = s.Copy().AppendTemplate(QuotedWrapperT, "q")
	dynamics.AddSimple(t, term.ReturnS(core.Answer.S(
		term.Sc(Template(term.Wrap(nil).Head())),
	)))

	t = s.Copy().AppendTemplate(QuotedCompoundT, "template", "args")
	t = dynamics.AddSimple(t, term.ReturnS(core.Answer.S(term.Sr("template"))))
}
