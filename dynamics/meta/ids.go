package meta

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/term/intern"
)

var (
	GetID         = term.Make("what is the internal identifier of the object []?")
	SettingFromID = term.Make("what is the setting with ID []?")
	NoSetting     = term.Make("there is no setting with that ID")
)

func getID(d dynamics.Dwimmer, s *term.SettingT, quoted term.T) term.T {
	var result int
	switch quoted.Head() {
	case represent.QuotedSetting:
		setting, err := represent.ToSetting(d, quoted)
		if err != nil {
			return term.Make("was asked to find the ID of a setting, " +
				"but while converting to a setting received []").T(err)
		}
		result = int(setting.ID)
	case term.Int(0).Head():
		result = int(term.IDer.PackInt(int(quoted.(term.Int))).(intern.ID))
	case term.Str("").Head():
		result = int(term.IDer.PackString(string(quoted.(term.Str))).(intern.ID))
	}
	return core.Answer.T(represent.Int(result))
}

func settingFromID(d dynamics.Dwimmer, s *term.SettingT, quotedID term.T) term.T {
	id, err := represent.ToInt(d, quotedID)
	if err != nil {
		return represent.ConversionError.T(quotedID, err)
	}
	setting, ok := term.ToSetting(term.SettingID(id))
	if ok {
		return core.Answer.T(represent.Setting(setting))
	} else {
		return NoSetting.T()
	}
}

func init() {
	dynamics.AddNativeResponse(GetID, 1, dynamics.Args1(getID))
	dynamics.AddNativeResponse(SettingFromID, 1, dynamics.Args1(settingFromID))
}
