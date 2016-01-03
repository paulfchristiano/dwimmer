package meta

import (
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	Run = term.Make("what is the result of running the setting [] until it produces a result?")

	Result   = term.Make("the result is [] and the final state of the setting is []")
	NoOutput = term.Make("there is no result, and the final state of the setting is []")
)

func init() {
	dynamics.AddNativeResponse(Run, 1, dynamics.Args1(run))
}

func run(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSettingT(d, quotedSetting)
	if err != nil {
		return represent.ConversionError.T(quotedSetting, err)
	}
	result := d.Run(setting)
	if result == nil {
		return NoOutput.T(represent.SettingT(setting))
	}
	return Result.T(represent.T(result), represent.SettingT(setting))
}
