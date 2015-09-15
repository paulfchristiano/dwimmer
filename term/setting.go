package term

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/term/intern"
)

type Setting struct {
	Id       SettingId
	Previous *Setting
	Last     SettingLine
	Slots    int
	Size     int
}

func Init() *Setting {
	id := new(intern.Id)
	id.Empty()
	return &Setting{
		Id: SettingId(*id),
	}
}

func (s *Setting) Append(line SettingLine, slots int) *Setting {
	id := new(intern.Id)
	*id = intern.Id(s.Id)
	lineId := line.LineId()
	id.Append(&lineId)
	return &Setting{
		Previous: s,
		Last:     line,
		Slots:    slots,
		Size:     s.Size + 1,
		Id:       SettingId(*id),
	}
}

func (s *Setting) Lines() []SettingLine {
	if s.Size == 0 {
		return []SettingLine{}
	}
	previous := s.Previous.Lines()
	return append(previous, s.Last)
}

type SettingLine interface {
	LineId() intern.Id
	Pack(intern.Packed) intern.Packed
}

func (a ActionCId) LineId() intern.Id {
	return intern.Id(a)
}

func (a ActionCId) Pack(x intern.Packed) intern.Packed {
	switch x := x.(type) {
	case *intern.Id:
		*x = intern.Id(a)
	}
	return a.ActionC().Pack(x)
}

func (a TemplateId) Pack(x intern.Packed) intern.Packed {
	switch x := x.(type) {
	case *intern.Id:
		*x = intern.Id(a)
	}
	return a.Template().Pack(x)
}

func (t TemplateId) LineId() intern.Id {
	return intern.Id(t)
}

type SettingS struct {
	Setting *Setting
	Names   []string
}

func (s *SettingS) Copy() *SettingS {
	result := &SettingS{
		Setting: s.Setting,
		Names:   make([]string, len(s.Names)),
	}
	copy(result.Names, s.Names)
	return result
}

func InitS() *SettingS {
	return &SettingS{
		Setting: Init(),
		Names:   []string{},
	}
}

func (s *SettingS) AppendTemplate(t TemplateId, names []string) *SettingS {
	s.Setting = s.Setting.Append(t, len(names))
	s.Names = append(s.Names, names...)
	return s
}

func (s *SettingS) AppendActionC(a ActionCId) *SettingS {
	s.Setting = s.Setting.Append(a, 0)
	return s
}

func (s *SettingS) Lines(names []string) []string {
	lines := s.Setting.Lines()
	result := make([]string, 0)
	index := 0
	for i, line := range lines {
		switch line := line.(type) {
		case ActionCId:
			a := line.ActionC()
			result = append(result, "")
			result = append(result, fmt.Sprintf("%d>%s", i, a.Uninstantiate(names).String()))
		case TemplateId:
			t := line.Template()
			newindex := index + t.Slots()
			result = append(result, fmt.Sprintf(
				"%d<%s", i,
				t.ShowWith(names[index:newindex]...),
			))
			index = newindex
		default:
			panic("Unknown kind of line!")
		}
	}
	return result
}

type SettingT struct {
	Setting *Setting
	Args    []T
}

func (s *SettingT) Copy() *SettingT {
	result := &SettingT{
		Setting: s.Setting,
		Args:    make([]T, len(s.Args)),
	}
	copy(result.Args, s.Args)
	return result
}

func (s *SettingT) AppendTerm(t T) *SettingT {
	s.Setting = s.Setting.Append(t.Head(), len(t.Values()))
	s.Args = append(s.Args, t.Values()...)
	return s
}

func (s *SettingT) AppendActionC(a ActionCId) *SettingT {
	s.Setting = s.Setting.Append(a, 0)
	return s
}

func (s *SettingS) Rollback(n int) {
	drop := 0
	for s.Setting.Size > n {
		drop += s.Setting.Slots
		s.Setting = s.Setting.Previous
	}
	s.Names = s.Names[:len(s.Names)-drop]
}

func (s *SettingT) Rollback(n int) {
	drop := 0
	for s.Setting.Size > n {
		drop += s.Setting.Slots
		s.Setting = s.Setting.Previous
	}
	s.Args = s.Args[:len(s.Args)-drop]
}
