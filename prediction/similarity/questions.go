package similarity

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/maps"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	DisplayHints = term.Make("display the list of suggested actions [], sorted from most to least promising, " +
		"in a context with list of names []; " +
		"return the list of strings that they display as, in the order they are displayed; " +
		"the most probable suggestions should be last")
	HintStrings = term.Make("the sequence of strings corresponding to suggested actions is given by the list []")
	GetHints    = term.Make("what actions should be recommended to the user in setting template []? " +
		"the most promising options should be at the top of the list")
	ToolMapping = term.Make("each letter is bound to its image in []")
	GetTools    = term.Make("what term templates should be recommended to the user in setting template []? " +
		"the most promising options should be at the top of the list")
	DisplayTools = term.Make("display the list of suggested term templates [], and bind them to letters; " +
		"the list of letters should be easy to type after ctrl+F; " +
		"return the resulting mapping from characters to strings")
)

func init() {
	dynamics.AddNativeResponse(DisplayTools, 1, dynamics.Args1(displayTools))
	dynamics.AddNativeResponse(GetTools, 1, dynamics.Args1(getTools))

	s := dynamics.ExpectQuestion(term.InitS(), GetHints, "Q", "template")
	s = dynamics.AddSimple(s, term.AskS(SuggestedActions.S(
		term.Sr("template"),
		term.Sc(represent.Int(6)),
	)))
	s = dynamics.ExpectAnswer(s, core.Answer, "A", "result")
	s = dynamics.AddSimple(s, term.ReturnS(core.Answer.S(term.Sr("result"))))

	dynamics.AddNativeResponse(DisplayHints, 2, dynamics.Args2(displayHints))

}

func displayHints(d dynamics.Dwimmer, context *term.SettingT, quotedSuggestion, quotedNames term.T) term.T {
	quotedActions, err := represent.ToList(d, quotedSuggestion)
	if err != nil {
		return represent.ConversionError.T(quotedSuggestion, err)
	}
	actions := make([]term.ActionC, len(quotedActions))
	for i, quoted := range quotedActions {
		actions[i], err = represent.ToActionC(d, quoted)
		if err != nil {
			return represent.ConversionError.T(quoted, err)
		}
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
	reverseActions(actions)
	strings := make([]string, len(actions))
	for i, action := range actions {
		strings[i] = action.Uninstantiate(names).String()
		d.Println(fmt.Sprintf("%d. %s", i, strings[i]))
	}
	if len(actions) > 0 {
		d.Println("")
	}
	quotedStrings := make([]term.T, len(strings))
	for i, s := range strings {
		quotedStrings[i] = represent.Str(s)
	}
	return HintStrings.T(represent.List(quotedStrings))
}

func reverseActions(list []term.ActionC) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func getTools(d dynamics.Dwimmer, context *term.SettingT, quotedTemplate term.T) term.T {
	template, err := represent.ToSetting(d, quotedTemplate)
	if err != nil {
		return represent.ConversionError.T(quotedTemplate, err)
	}
	templates := SuggestedTemplates(d, template, 6)
	quotedTemplates := make([]term.T, len(templates))
	for i, template := range templates {
		quotedTemplates[i] = represent.Template(template.ID())
	}
	return core.Answer.T(represent.List(quotedTemplates))
}

func displayTools(d dynamics.Dwimmer, context *term.SettingT, quotedSuggestion term.T) term.T {
	quotedTemplates, err := represent.ToList(d, quotedSuggestion)
	if err != nil {
		return represent.ConversionError.T(quotedSuggestion, err)
	}
	templates := make([]*term.Template, len(quotedTemplates))
	for i, quoted := range quotedTemplates {
		id, err := represent.ToTemplate(d, quoted)
		if err != nil {
			return represent.ConversionError.T(quoted, err)
		}
		templates[i] = id.Template()
	}
	letters := []rune{'a', 's', 'd', 'w', 'e', 'j', 'r', 'z', 'x', 'u', 'i', 'l', 'k'}
	mapping := maps.Empty.T()
	for i, template := range templates {
		if i >= len(letters) {
			break
		}
		s := template.String()
		d.Println(fmt.Sprintf("%c. %s", letters[i], s))
		mapping = maps.Cons.T(represent.Rune(letters[i]), represent.Str(s), mapping)
	}
	if len(templates) > 0 {
		d.Println("")
	}
	return ToolMapping.T(mapping)
}
