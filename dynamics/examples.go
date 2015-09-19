package dynamics

import "github.com/paulfchristiano/dwimmer/term"

func ExpectQuestion(s *term.SettingS, t term.TemplateId, names ...string) *term.SettingS {
	result := s.Copy()
	result.AppendTemplate(ParentChannel, names[0])
	result.AppendTemplate(t, names[1:]...)
	return result
}

func ExpectAnswer(s *term.SettingS, t term.TemplateId, names ...string) *term.SettingS {
	result := s.Copy()
	result.AppendTemplate(OpenChannel, names[0])
	result.AppendTemplate(t, names[1:]...)
	return result
}

var (
	OpenChannel   = term.Make("@[]")
	ParentChannel = term.Make("@[]*")
)
