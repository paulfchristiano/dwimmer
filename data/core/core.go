package core

import "github.com/paulfchristiano/dwimmer/term"

var (
	Answer    = term.Make("the answer to the posed question is []")
	OK        = term.Make("OK")
	Yes       = term.Make("yes")
	No        = term.Make("no")
	NoAnswer  = term.Make("no answer was given to the posed question")
	AnswerAnd = term.Make("the answer to the posed question is [], and the question also satisfies []")
)
