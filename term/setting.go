package term

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/term/intern"
)

//NOTE Setting is immutable, SettingS and SettingT are mutable
//this seems good for performance, but may cause bugs...
//Also, it exploits the fact that a setting can't change after going into a channel :(

type Setting struct {
	Id       SettingId
	Previous *Setting
	Last     SettingLine
	Slots    int
	Size     int
}

func (s *Setting) Rollback(n int) *Setting {
	if n < 0 {
		n = s.Size + n
	}
	if s.Size <= n {
		return s
	}
	return s.Previous.Rollback(n)
}

func Init() *Setting {
	id := new(intern.Id)
	id.Empty()
	return &Setting{
		Id: SettingId(*id),
	}
}

func (s *Setting) Append(line SettingLine, slotss ...int) *Setting {
	var slots int
	if len(slotss) == 0 {
		slots = line.Slots()
	} else {
		slots = slotss[0]
	}
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

func (s *Setting) TotalSlots() int {
	if s.Size == 0 {
		return 0
	}
	return s.Slots + s.Previous.TotalSlots()
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
	Slots() int
}

func (a ActionCId) LineId() intern.Id {
	return intern.Id(a)
}

func (a ActionCId) Slots() int {
	return 0
}

func (t TemplateId) Slots() int {
	return t.Template().Slots()
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

func (s *SettingS) AppendTemplate(t TemplateId, names ...string) *SettingS {
	s.Setting = s.Setting.Append(t, len(names))
	for i := range names {
		for j := range s.Names {
			if names[i] == s.Names[j] {
				panic("duplicate name!")
			}
		}
	}
	s.Names = append(s.Names, names...)
	return s
}

func (s *SettingS) AppendAction(a ActionC) *SettingS {
	id := a.Id()
	s.Setting = s.Setting.Append(id, 0)
	return s
}

func (s *SettingS) Lines() []string {
	lines := s.Setting.Lines()
	result := make([]string, 0)
	index := 0
	for i, line := range lines {
		switch line := line.(type) {
		case ActionCId:
			a := line.ActionC()
			result = append(result, "", fmt.Sprintf("%d< %s", i, a.Uninstantiate(s.Names).String()))
		case TemplateId:
			t := line.Template()
			newindex := index + t.Slots()
			if len(s.Names) < newindex {
				panic(fmt.Sprintf(
					"s.Names = %v, template = %v, index = %d, newindex = %d",
					s.Names, t, index, newindex,
				))
			}
			result = append(result, fmt.Sprintf(
				"%d> %s", i,
				t.ShowWith(s.Names[index:newindex]...),
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

func InitT() *SettingT {
	return &SettingT{
		Setting: Init(),
		Args:    []T{},
	}
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

func (s *SettingT) AppendAction(a ActionC) *SettingT {
	id := a.Id()
	s.Setting = s.Setting.Append(id, 0)
	return s
}

func (s *SettingS) Rollback(n int) *SettingS {
	if n < 0 {
		n = s.Setting.Size + n
	}
	drop := 0
	for s.Setting.Size > n {
		drop += s.Setting.Slots
		s.Setting = s.Setting.Previous
	}
	s.Names = s.Names[:len(s.Names)-drop]
	return s
}

func (s *SettingT) Rollback(n int) *SettingT {
	if n < 0 {
		n = s.Setting.Size + n
	}
	drop := 0
	for s.Setting.Size > n {
		drop += s.Setting.Slots
		s.Setting = s.Setting.Previous
	}
	s.Args = s.Args[:len(s.Args)-drop]
	return s
}
