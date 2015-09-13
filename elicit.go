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

var ActionQ = term.Make("what action should be taken in the setting []?")

func init() {
	dynamics.AddNativeResponse(ActionQ, 1, dynamics.Args1(findAction))
}

func findAction(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to decide what to do in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	settingId := term.IdSetting(setting)
	transition, ok := d.Get(settingId)
	if !ok {
		action := ElicitAction(d, settingId, true)
		transition = dynamics.SimpleTransition{action}
		d.Save(settingId, transition)
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
	actions, _ := similarity.Suggestions(d, s.Setting(), n)
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

func ElicitAction(d dynamics.Dwimmer, id term.SettingId, hints bool) term.ActionC {
	setting := id.Setting()
	settingS := addNames(setting)
	ShowSettingS(d, settingS)
	hint_strings := []string{}
	if hints {
		hint_strings = GetHints(d, settingS, 6)
		for i, hint := range hint_strings {
			d.Println(fmt.Sprintf("%d. %s", i, hint))
		}
		d.Println("")
	}
	for {
		input := d.Readln(" < ", hint_strings)
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
					a.IntArgs[i] = len(setting.Outputs) - 1
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

var allNames = "xyzwijklmnstuvabcdefg"

func makeNames(n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = allNames[i : i+1]
	}
	return result
}

func addNames(s *term.Setting) *term.SettingS {
	names := makeNames(s.Slots())
	return &term.SettingS{term.IdSetting(s), names}
}
