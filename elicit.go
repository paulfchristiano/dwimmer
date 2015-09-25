package dwimmer

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/parsing"
	"github.com/paulfchristiano/dwimmer/prediction/similarity"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	ActionQ     = term.Make("what action should be taken in the setting []?")
	FallThrough = term.Make("there is no compiled transition in the setting []; " +
		"compile one if desired and specify what should be done")
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
)

func init() {
	dynamics.AddNativeResponse(ActionQ, 1, dynamics.Args1(findAction))
	dynamics.AddNativeResponse(FallThrough, 1, dynamics.Args1(fallThrough))
	dynamics.AddNativeResponse(SetTransition, 2, dynamics.Args2(setTransition))
	dynamics.AddNativeResponse(GetTransition, 1, dynamics.Args1(getTransition))
	dynamics.AddNativeResponse(GetContinuations, 1, dynamics.Args1(getContinuations))
	dynamics.AddNativeResponse(AllSettings, 0, dynamics.Args0(allSettings))
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
	action := ElicitAction(d, settingT.Setting, true)
	transition := dynamics.SimpleTransition{action}
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
		action := ElicitAction(d, setting, true)
		transition = dynamics.SimpleTransition{action}
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
func GetHints(d dynamics.Dwimmer, s *term.SettingS, n int) []string {
	result := []string{}
	actions, _ := similarity.Suggestions(d, s.Setting, n)
	for _, actionC := range actions {
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

func ElicitAction(d dynamics.Dwimmer, s *term.Setting, hints bool) term.ActionC {
	settingS := addNames(s)
	ShowSettingS(d, settingS)
	hint_strings := []string{}
	tool_map := make(map[rune]string)
	if hints {
		hint_strings = GetHints(d, settingS, 6)
		tools := similarity.SuggestedTemplates(d, s, 6)
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
					a.IntArgs[i] = s.Size - 1
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
