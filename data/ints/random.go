package ints

import (
	"math/rand"
	"time"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	source = rand.New(rand.NewSource(time.Now().UnixNano()))
	Random = term.Make("generate a random integer")
)

func init() {
	dynamics.AddNativeResponse(Random, 0, dynamics.Args0(random))
}

func random(d dynamics.Dwimmer, context *term.SettingT) term.T {
	return core.Answer.T(represent.Int(source.Int()))
}
