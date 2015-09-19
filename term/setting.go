package term

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/term/intern"
)

//NOTE Setting is immutable, SettingS and SettingT are mutable
//this seems good for performance, but may cause bugs...
//Also, it exploits the fact that a setting can't change after going into a channel :(

type Setting struct {
	ID       SettingID
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

func Init(ider intern.Packer) *Setting {
	return &Setting{
		ID: SettingID(ider.PackList([]intern.Packed{}).(intern.ID)),
	}
}

func (s *Setting) Append(ider intern.Packer, line SettingLine, slotss ...int) *Setting {
	var slots int
	if len(slotss) == 0 {
		slots = line.Slots(ider)
	} else {
		slots = slotss[0]
	}
	return &Setting{
		Previous: s,
		Last:     line,
		Slots:    slots,
		Size:     s.Size + 1,
		ID: SettingID(ider.AppendToPacked(
			intern.ID(s.ID),
			line.LineID(ider),
		).(intern.ID)),
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
	LineID(intern.Packer) intern.ID
	Slots(intern.Packer) int
}

func (a ActionCID) LineID(ider intern.Packer) intern.ID {
	return intern.ID(a)
}

func (a ActionCID) Slots(ider intern.Packer) int {
	return 0
}

func (t TemplateID) Slots(ider intern.Packer) int {
	return ToTemplate(ider, t).Slots()
}

func (t TemplateID) LineID(ider intern.Packer) intern.ID {
	return intern.ID(t)
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

func InitS(ider intern.Packer) *SettingS {
	return &SettingS{
		Setting: Init(ider),
		Names:   []string{},
	}
}

func (s *SettingS) AppendTemplate(ider intern.Packer, t TemplateID, names ...string) *SettingS {
	s.Setting = s.Setting.Append(ider, t, len(names))
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

func (s *SettingS) AppendAction(ider intern.Packer, a ActionC) *SettingS {
	id := IDActionC(ider, a)
	s.Setting = s.Setting.Append(ider, id, 0)
	return s
}

func (s *SettingS) Lines(ider intern.Packer) []string {
	lines := s.Setting.Lines()
	result := make([]string, 0)
	index := 0
	for i, line := range lines {
		switch line := line.(type) {
		case ActionCID:
			a := ToActionC(ider, line)
			result = append(result, "", fmt.Sprintf("%d< %s", i, a.Uninstantiate(s.Names).String(ider)))
		case TemplateID:
			t := ToTemplate(ider, line)
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

func InitT(ider intern.Packer) *SettingT {
	return &SettingT{
		Setting: Init(ider),
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

func (s *SettingT) AppendTerm(ider intern.Packer, t T) *SettingT {
	s.Setting = s.Setting.Append(ider, t.Head(), len(t.Values()))
	s.Args = append(s.Args, t.Values()...)
	return s
}

func (s *SettingT) AppendAction(ider intern.Packer, a ActionC) *SettingT {
	id := IDActionC(ider, a)
	s.Setting = s.Setting.Append(ider, id, 0)
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
