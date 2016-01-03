package main

import "github.com/paulfchristiano/dwimmer/ui"

func main() {
	t := &ui.Term{}
	t.InitUI()
	defer t.CloseUI()
	t.Readln(" >> ", []string{"A"}, map[rune]string{'a': "BBB[testing]CCC"})

}
