package represent_test

import (
	"testing"

	"github.com/paulfchristiano/dwimmer"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/term"
)

func TestRepresentations(t *testing.T) {
	d := dwimmer.Dwimmer()
	template := term.Make("term with argument [] and second half here")
	template2, err := represent.ToTemplate(d, represent.Template(template))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if template != template2 {
		t.Errorf("%v != %v", template != template2)
	}
	setting := &term.Setting{Outputs: []term.TemplateId{template}, Inputs: []term.ActionCId{}}
	setting2, err := represent.ToSetting(d, represent.Setting(setting))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if term.IdSetting(setting) != term.IdSetting(setting2) {
		t.Errorf("%v != %v", setting, setting2)
	}
	action := term.ReturnC(term.Cr(3))
	action2, err := represent.ToActionC(d, represent.ActionC(action))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if term.IdActionC(action) != term.IdActionC(action2) {
		t.Errorf("%v != %v", action, action2)
	}
	stub := term.Make("stub")
	tm := template.T(stub.T())
	tm2, err := represent.ToT(d, represent.T(tm))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if tm2.String() != tm.String() {
		t.Errorf("%v != %v", tm2, tm)
	}
	rep, err := d.Answer(represent.Explicit.T(represent.T(tm)))
	if err != nil {
		t.Errorf("failed to make representation explicit: %v", err)
	}
	tm3, err := represent.ToT(d, rep)
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if tm3.String() != tm.String() {
		t.Errorf("%v != %v", tm3, tm)
	}

	n := -127
	n2, err := represent.ToInt(d, represent.Int(n))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if n != n2 {
		t.Errorf("%v != %v", n, n2)
	}

	s := "hello ₳"
	s2, err := represent.ToStr(d, represent.Str(s))
	if err != nil {
		t.Errorf("received error %v", err)
	}
	if s != s2 {
		t.Errorf("%s != %s", s, s2)
	}
}
