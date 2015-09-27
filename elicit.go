package dwimmer

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/maps"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/dynamics/meta"
	"github.com/paulfchristiano/dwimmer/parsing"
	"github.com/paulfchristiano/dwimmer/prediction/similarity"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	Names       = term.Make("the setting was displayed with the sequence of names []")
	ActionQ     = term.Make("what action should be taken in the setting []?")
	FallThrough = term.Make("there is no compiled transition in the setting []; " +
		"compile one if desired and specify what should be done")
	FallThroughAlt = term.Make("there is no compiled transition in the setting []; " +
		"compile one if desired and specify what should be done. the process should eschew built in functions")
	SetTransition = term.Make("from now on, the transition [] should be made in setting []")
	GetTransition = term.Make("what transition will be reflexively taken in setting []?")
	AllSettings   = term.Make("what is the list of all settings " +
		"about which data is available?")
	GetContinuations = term.Make("what settings for which data is available " +
		"can be formed by adding one line to []?")
	NoCompiledAction = term.Make("there is no compiled transition in that setting")
	TakeTransition   = term.Make("transition [] should be taken")
	TransitionGiven  = term.Make("if [] was given as a reply to the question [], " +
		"what transition should be taken?")
	MakeTransition             = term.Make("what is a transition that takes action []?")
	GetTemplate                = term.Make("what setting template underlies []?")
	ElicitActionQ              = term.Make("elicit an action from the user to take in setting []")
	ElicitActionQTemplate      = term.Make("elicit an action from the user to take in setting template []")
	ElicitActionQTemplateHints = term.Make("elicit an action from the user to take in setting template [] " +
		"with list of hints [] and list of suggested term templates []")
	DisplayTemplate = term.Make("display the setting template []")
	GetActionInput  = term.Make("prompt the user for a string that will be parsed as an action; " +
		"all necessary context is implied by the information currently on the screen; " +
		"the list of hint strings is [] and the binding of keys to suggested term templates is []")
	ParseInputAsAction = term.Make("if the user typed the string [], " +
		"in a setting with list of names [], " +
		"what indexed action are they most likely referring to?")
	ShouldRecord = term.Make("if the user took action [] in a setting, " +
		"should the same action be taken when that setting arises in the future?")
)

func init() {
	//dynamics.AddNativeResponse(FallThrough, 1, dynamics.Args1(fallThrough))
	dynamics.AddNativeResponse(DisplayTemplate, 1, dynamics.Args1(displayTemplate))
	dynamics.AddNativeResponse(GetActionInput, 2, dynamics.Args2(getInput))
	dynamics.AddNativeResponse(ParseInputAsAction, 2, dynamics.Args2(parseInput))

	var s *term.SettingS
	s = dynamics.ExpectQuestion(term.InitS(), FallThrough, "Q", "setting")
	s = dynamics.AddSimple(s, term.AskS(GetTemplate.S(term.Sr("setting"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A", "template")
	s = dynamics.AddSimple(s, term.AskS(ElicitActionQ.S(term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A2", "action")
	s = dynamics.AddSimple(s, term.AskS(MakeTransition.S(term.Sr("action"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A3", "transition")
	s = dynamics.AddSimple(s, term.AskS(ShouldRecord.S(term.Sr("action"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "AA", "should")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("should")))
	s.AppendTemplate(core.Yes)
	s = dynamics.AddSimple(s, term.AskS(SetTransition.S(term.Sr("transition"), term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, core.OK, "A4")
	s = dynamics.AddSimple(s, term.ReturnS(TakeTransition.S(term.Sr("transition"))))

	s = dynamics.ExpectQuestion(term.InitS(), ShouldRecord, "Q", "action")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("action")))
	s.AppendTemplate(represent.QuotedActionC, "act", "args", "indices")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("act")))
	for act, template := range represent.ActionLookup {
		t := s.Copy().AppendTemplate(template)
		switch act {
		case term.Replace, term.Correct, term.Replay:
			dynamics.AddSimple(t, term.ReturnS(core.Answer.S(core.No.S())))
		default:
			dynamics.AddSimple(t, term.ReturnS(core.Answer.S(core.Yes.S())))
		}
	}

	s = dynamics.ExpectQuestion(term.InitS(), MakeTransition, "Q", "action")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(represent.QuotedSimpleTransition.S(term.Sr("action")))))

	s = dynamics.ExpectQuestion(term.InitS(), GetTemplate, "Q", "setting")
	t := dynamics.AddSimple(s, term.ViewS(term.Sr("setting")))
	s = t.Copy().AppendTemplate(represent.QuotedTemplate, "parts")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(term.Sr("setting"))))
	s = t.Copy().AppendTemplate(represent.QuotedSettingT, "template", "args")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(term.Sr("template"))))

	s = dynamics.ExpectQuestion(term.InitS(), ElicitActionQ, "Q", "setting")
	s = dynamics.AddSimple(s, term.AskS(GetTemplate.S(term.Sr("setting"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A", "template")
	s = dynamics.AddSimple(s, term.AskS(ElicitActionQTemplate.S(term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A2", "action")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(term.Sr("action"))))

	s = dynamics.ExpectQuestion(term.InitS(), ElicitActionQTemplate, "Q", "template")
	s = dynamics.AddSimple(s, term.AskS(similarity.GetHints.S(term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A", "hints")
	s = dynamics.AddSimple(s, term.AskS(similarity.GetTools.S(term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A2", "tools")
	s = dynamics.AddSimple(s, term.AskS(ElicitActionQTemplateHints.S(
		term.Sr("template"), term.Sr("hints"), term.Sr("tools"),
	)))
	s = dynamics.ExpectAnswer(s, core.Answer, "A3", "action")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(term.Sr("action"))))

	s = dynamics.ExpectQuestion(term.InitS(), ElicitActionQTemplateHints, "Q", "template", "hints", "tools")
	s = dynamics.AddSimple(s, term.AskS(meta.Clear.S()))
	s = dynamics.ExpectAnswer(s, core.OK, "A")
	s = dynamics.AddSimple(s, term.AskS(DisplayTemplate.S(term.Sr("template"))))
	s = dynamics.ExpectAnswer(s, Names, "A2", "names")
	s = dynamics.AddSimple(s, term.AskS(similarity.DisplayTools.S(term.Sr("tools"))))
	s = dynamics.ExpectAnswer(s, similarity.ToolMapping, "A3", "toolmap")
	s = dynamics.AddSimple(s, term.AskS(similarity.DisplayHints.S(term.Sr("hints"), term.Sr("names"))))
	s = dynamics.ExpectAnswer(s, similarity.HintStrings, "A4", "hintstrings")
	s = dynamics.AddSimple(s, term.AskS(GetActionInput.S(term.Sr("hintstrings"), term.Sr("toolmap"))))
	s = dynamics.ExpectAnswer(s, core.Answer, "A5", "input")
	s = dynamics.AddSimple(s, term.AskS(ParseInputAsAction.S(term.Sr("input"), term.Sr("names"))))

	t = dynamics.ExpectAnswer(s.Copy(), core.Answer, "A6", "action")
	t = dynamics.AddSimple(t, term.ReturnS(core.Answer.S(term.Sr("action"))))

	t = dynamics.ExpectAnswer(s.Copy(), GoMeta, "A6")
	t = dynamics.AddSimple(t, term.MetaS())

	dynamics.AddNativeResponse(ActionQ, 1, dynamics.Args1(findAction))
	dynamics.AddNativeResponse(SetTransition, 2, dynamics.Args2(setTransition))
	dynamics.AddNativeResponse(GetTransition, 1, dynamics.Args1(getTransition))
	dynamics.AddNativeResponse(GetContinuations, 1, dynamics.Args1(getContinuations))
	dynamics.AddNativeResponse(AllSettings, 0, dynamics.Args0(allSettings))
}

func getInput(d dynamics.Dwimmer, context *term.SettingT, hintstrings, toolmap term.T) term.T {
	quotedHints, err := represent.ToList(d, hintstrings)
	if err != nil {
		return represent.ConversionError.T(hintstrings, err)
	}
	hints := make([]string, len(quotedHints))
	for i, quoted := range quotedHints {
		hints[i], err = represent.ToStr(d, quoted)
		if err != nil {
			return represent.ConversionError.T(quoted, err)
		}
	}
	tools := make(map[rune]string)
	for toolmap.Head() != maps.Empty {
		switch toolmap.Head() {
		case maps.Cons:
			vs := toolmap.Values()
			c, err := represent.ToRune(d, vs[0])
			if err != nil {
				return represent.ConversionError.T(vs[0], err)
			}
			s, err := represent.ToStr(d, vs[1])
			if err != nil {
				return represent.ConversionError.T(vs[1], err)
			}
			tools[c] = s
			toolmap = vs[2]
		default:
			context.AppendTerm(UnrecognizedDictionary.T())
			return nil
		}
	}
	input := d.Readln(" < ", hints, tools)
	return core.Answer.T(represent.Str(input))
}

var (
	UnrecognizedDictionary = term.Make("the representation of the given tool bindings cannot be handled " +
		"by the builtin prompt")
	GoMeta     = term.Make("the user's input indicates that they want to jump up a meta level")
	ParseError = term.Make("the user's input could not be parsed into an action or term")
	IsTerm     = term.Make("the user's input could be parsed into a term, but not an action")
)

func parseInput(d dynamics.Dwimmer, context *term.SettingT, quotedInput, quotedNames term.T) term.T {
	input, err := represent.ToStr(d, quotedInput)
	if err != nil {
		return represent.ConversionError.T(quotedInput, err)
	}
	quotedList, err := represent.ToList(d, quotedNames)
	if err != nil {
		return represent.ConversionError.T(quotedNames, err)
	}
	names := make([]string, len(quotedList))
	for i, quoted := range quotedList {
		names[i], err = represent.ToStr(d, quoted)
		if err != nil {
			return represent.ConversionError.T(quoted, err)
		}
	}
	if input == "jump up" {
		return GoMeta.T()
	}
	a := parsing.ParseAction(input, names)
	if a == nil {
		c := parsing.ParseTerm(input, names)
		if c != nil {
			switch c := c.(type) {
			case *term.CompoundC:
				if questionLike(c) {
					a = new(term.ActionC)
					*a = term.AskC(c)
				}
			case term.ReferenceC:
				a = new(term.ActionC)
				*a = term.ViewC(c)
			}
			context.AppendTerm(IsTerm.T())
			return nil
		} else {
			context.AppendTerm(ParseError.T())
			return nil
		}
	} else {
		for _, n := range a.IntArgs {
			if n < 0 {
				context.AppendTerm(ParseError.T())
				return nil
			}
		}
		return core.Answer.T(represent.ActionC(*a))
	}
}

func getContinuations(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to return the continuations from a setting, " +
			"but while converting to a setting received []").T(err)
	}
	continuations := d.Continuations(setting)
	result := make([]term.T, len(continuations))
	for i, c := range continuations {
		result[i] = represent.Setting(c)
	}
	return core.Answer.T(represent.List(result))
}

func allSettings(d dynamics.Dwimmer, s *term.SettingT) term.T {
	queue := []*term.Setting{term.Init()}
	for k := 0; k < len(queue); k++ {
		top := queue[k]
		queue = append(queue, d.Continuations(top)...)
	}
	result := make([]term.T, len(queue))
	for i, setting := range queue {
		result[i] = represent.Setting(setting)
	}
	return core.Answer.T(represent.List(result))
}

func setTransition(d dynamics.Dwimmer, s *term.SettingT, quotedTransition, quotedSetting term.T) term.T {
	transition, err := represent.ToTransition(d, quotedTransition)
	if err != nil {
		return term.Make("asked to set a setting to transition [], "+
			"but while converting to a transition received []").T(quotedTransition, err)
	}
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to set a transition in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	d.Save(setting, transition)
	return core.OK.T()
}

func getTransition(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to get a transition in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	result, ok := d.Get(setting)
	if !ok {
		return NoCompiledAction.T()
	}
	return core.Answer.T(represent.Transition(result))
}

func fallThrough(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	settingT, err := represent.ToSettingT(d, quotedSetting)
	if err != nil {
		return term.Make("asked to decide what to do in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	transition := ElicitAction(d, s, settingT.Setting)
	if shouldSave(transition) {
		d.Save(settingT.Setting, transition)
	}
	return TakeTransition.T(represent.Transition(transition))
}

func shouldSave(t dynamics.Transition) bool {
	switch t := t.(type) {
	case dynamics.NativeTransition:
		return true
	case dynamics.SimpleTransition:
		switch t.Action.Act {
		case term.Replay, term.Replace, term.Correct:
			return false
		default:
			return true
		}
	}
	panic("unknown type of transition")
}

func findAction(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to decide what to do in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	transition, ok := d.Get(setting)
	if !ok {
		transition := ElicitAction(d, s, setting)
		d.Save(setting, transition)
	}
	return core.Answer.T(represent.Transition(transition))
}

func ShowSettingS(d dynamics.Dwimmer, settingS *term.SettingS) {
	d.Clear()
	for _, line := range settingS.Lines() {
		d.Println(line)
	}
}

//TODO only do this sometimes, or control better the length, or something...
func GetHints(d dynamics.Dwimmer, context *term.SettingT, s *term.SettingS, n int) []string {
	suggestions, err := d.Answer(similarity.SuggestedActions.T(
		represent.Setting(s.Setting),
		represent.Int(n),
	), context)
	if err != nil {
		return []string{}
	}
	suggestionList, err := represent.ToList(d, suggestions)
	if err != nil {
		return []string{}
	}
	result := []string{}
	for _, suggestion := range suggestionList {
		actionC, err := represent.ToActionC(d, suggestion)
		if err != nil {
			continue
		}
		actionS := actionC.Uninstantiate(s.Names)
		result = append(result, actionS.String())
	}
	reverseHints(result)
	return result
}

func reverseHints(l []string) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func displayTemplate(d dynamics.Dwimmer, context *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return represent.ConversionError.T(quotedSetting, err)
	}
	settingS := addNames(setting)
	ShowSettingS(d, settingS)
	d.Println("")
	d.Println("")
	quotedNames := make([]term.T, len(settingS.Names))
	for i, name := range settingS.Names {
		quotedNames[i] = represent.Str(name)
	}
	return Names.T(represent.List(quotedNames))
}

func ElicitAction(d dynamics.Dwimmer, context *term.SettingT, setting *term.Setting) dynamics.Transition {
	for {
		Q := FallThrough.T(represent.Setting(setting))
		result := dynamics.SubRun(d, Q, context)
		switch result.Head() {
		case TakeTransition:
			transition, err := represent.ToTransition(d, result.Values()[0])
			if err == nil {
				return transition
			}
		case core.OK:
		default:
			result, err := d.Answer(TransitionGiven.T(result, Q))
			if err == nil {
				transition, err := represent.ToTransition(d, result)
				if err == nil {
					return transition
				}
			}
		}
	}
}

func OldElicitAction(d dynamics.Dwimmer, context *term.SettingT, subject *term.Setting, hints bool) term.ActionC {
	settingS := addNames(subject)
	ShowSettingS(d, settingS)
	hint_strings := []string{}
	tool_map := make(map[rune]string)
	if hints {
		hint_strings = GetHints(d, context, settingS, 6)
		tools := similarity.SuggestedTemplates(d, subject, 6)
		if len(hint_strings) > 0 || len(tools) > 0 {
			d.Println("")
		}
		if len(tools) > 0 {
			d.Println("")
		}
		tips := []rune{'a', 's', 'd', 'w', 'e', 'j'}
		for i, tool := range tools {
			tool_map[tips[i]] = tool.String()
			d.Println(fmt.Sprintf("%c: %v", tips[i], tool))
		}
		if len(hint_strings) > 0 {
			d.Println("")
		}
		for i, hint := range hint_strings {
			d.Println(fmt.Sprintf("%d. %s", i, hint))
		}
	}
	d.Println("")
	for {
		input := d.Readln(" < ", hint_strings, tool_map)
		if input == "jump up" {
			StartShell(d, Interrupted.T(represent.SettingT(context)))
		}
		a := parsing.ParseAction(input, settingS.Names)
		if a == nil {
			c := parsing.ParseTerm(input, settingS.Names)
			if c != nil {
				switch c := c.(type) {
				case *term.CompoundC:
					if questionLike(c) {
						a = new(term.ActionC)
						*a = term.AskC(c)
					}
				case term.ReferenceC:
					a = new(term.ActionC)
					*a = term.ViewC(c)
				}
				d.Println("please input an action (ask, view, return, close, delete, correct, or tell)")
			} else {
				d.Println("that response wasn't parsed correctly")
			}
		}
		if a != nil {
			for i, n := range a.IntArgs {
				if n == -1 {
					a.IntArgs[i] = subject.Size - 1
				}
			}
			return *a
		}
	}
}

func questionLike(c term.C) bool {
	for _, char := range c.String() {
		if char == '?' {
			return true
		}
	}
	return false
}

var allNames = "wxyzXYZWbcdABCDefghstuvijklmn"

func makeNames(n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = allNames[i : i+1]
	}
	return result
}

func addNames(s *term.Setting) *term.SettingS {
	names := makeNames(s.TotalSlots())
	return &term.SettingS{s, names}
}
