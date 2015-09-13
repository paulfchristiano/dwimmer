package main

import (
	"github.com/paulfchristiano/dwimmer"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	Home = term.Make("Home")
)

func main() {
	d := dwimmer.NewDwimmer()
	defer d.Close()
	setting := term.InitT()
	setting.AppendTerm(Home.T())
	d.Debug("testing!")
	for {
		actionC := dwimmer.ElicitAction(d, setting.SettingId, false)
		setting.AppendAction(actionC)
		d.DoC(actionC, setting)
	}
}
