package main

import (
	"github.com/paulfchristiano/dwimmer"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
)

var (
	Home = term.Make("Home")
)

func main() {
	defer dwimmer.DisplayStackError()
	ui.Init()
	defer ui.Close()
	d := dwimmer.Dwimmer()
	setting := term.InitT()
	setting.AppendTerm(Home.T())
	for {
		actionC := dwimmer.ElicitAction(d, setting.SettingId, false)
		setting.AppendAction(actionC)
		d.DoC(actionC, setting)
	}
}
