package dwimmer

import (
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	RunShell     = term.Make("run a shell with prompt [], and reply with the user's return value")
	shellSetting = term.InitT().AppendTerm(term.Make("[Home]").T())
	oldShells    []*term.SettingT
)

func startWithPrompt(d dynamics.Dwimmer, setting *term.SettingT, question term.T) term.T {
	return StartShell(d, dynamics.Parent(setting), question)
}

func pushShell() {
	oldShells = append(oldShells, shellSetting.Copy())
}

func popShell() {
	shellSetting = oldShells[len(oldShells)-1]
	oldShells = oldShells[:len(oldShells)-1]
}

func Show(t term.T) {
	shellSetting.AppendTerm(t)
}

func StartShell(d dynamics.Dwimmer, ts ...term.T) term.T {
	pushShell()
	defer popShell()

	for _, t := range ts {
		shellSetting.AppendTerm(t)
	}

	for {
		actionC := ElicitAction(d, term.InitT(), shellSetting.Setting, false)
		shellSetting.AppendAction(actionC)
		result := DoC(d, actionC, shellSetting)
		if result != nil {
			return result
		}
	}
}

func init() {
	dynamics.AddNativeResponse(RunShell, 1, dynamics.Args1(startWithPrompt))
}
