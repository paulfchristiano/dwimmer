package similarity

import (
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

func SuggestedTemplates(d dynamics.Dwimmer, s *term.Setting, n int, templatess ...[]*term.Template) (templates []*term.Template) {
	for _, x := range templatess {
		templates = append(templates, x...)
	}
	if s.Size == 0 {
		return
	}
	defer func() {
		if len(templates) > n {
			templates = templates[:n]
		}
	}()
	ShouldAddTemplate := func(template *term.Template) bool {
		for _, other := range templates {
			if other.String() == template.String() {
				return false
			}
		}
		excluded := []term.T{
			term.Int(0), term.Str(""), term.Quoted{}, term.Wrapper{},
		}
		for _, exclude := range excluded {
			if template == exclude.Head().Template() {
				return false
			}
		}
		excludedIDs := []term.TemplateID{
			dynamics.ParentChannel, dynamics.OpenChannel,
		}
		for _, exclude := range excludedIDs {
			if template.ID() == exclude {
				return false
			}
		}
		return true
	}
	AddTemplate := func(template *term.Template) {
		if ShouldAddTemplate(template) {
			templates = append(templates, template)
		}
	}
	AddTemplates := func(action term.ActionC) {
		actionTemplates := action.AllTemplates()
		for _, template := range actionTemplates {
			AddTemplate(template)
		}
	}

	actions, _ := Suggestions(d, s, n)
	for _, action := range actions {
		AddTemplates(action)
	}
	lastLine := s.Last
	switch line := lastLine.(type) {
	case term.ActionCID:
		AddTemplates(line.ActionC())
	case term.TemplateID:
		AddTemplate(line.Template())
	}
	if len(templates) < n {
		templates = SuggestedTemplates(d, s.Previous, n, templates)
	}
	return
}
