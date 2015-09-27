package maps

import "github.com/paulfchristiano/dwimmer/term"

var (
	Empty    = term.Make("the map that doesn't send anything to anything")
	Cons     = term.Make("the map sends [] to [] and maps all other inputs according to []")
	Lookup   = term.Make("what is the image of [] in the map []?")
	Insert   = term.Make("what map maps [] to [] and all other inputs according to []?")
	NotFound = term.Make("the requested key was not found")
)
