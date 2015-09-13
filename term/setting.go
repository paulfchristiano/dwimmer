package term

import "fmt"

type Setting struct {
	Outputs []TemplateId
	Inputs  []ActionCId
}

func (s *Setting) RollBack(n int) *Setting {
	return &Setting{
		s.Outputs[:n],
		s.Inputs[:n-1],
	}
}

func (s *Setting) Slots() int {
	result := 0
	for _, output := range s.Outputs {
		result += output.Template().Slots()
	}
	return result
}

func (s *Setting) Question() *Template {
	return s.Outputs[0].Template()
}

func (id SettingId) Head() SettingId {
	return id
}

//TODO I could make the initialization faster if the singleton setting
//had the same ID as the starting term, which is already in hand when initializing...
var InitId = IdSetting(Init())

func InitT() *SettingT {
	return &SettingT{InitId, make([]*SettingT, 0), make([]T, 0)}
}

func InitS() *SettingS {
	return &SettingS{InitId, make([]string, 0)}
}

func Init() *Setting {
	return &Setting{make([]TemplateId, 0), make([]ActionCId, 0)}
}

type SettingT struct {
	SettingId
	Children []*SettingT
	Args     []T
}
type SettingS struct {
	SettingId
	Names []string
}

func (s *Setting) Id() SettingId {
	return IdSetting(s)
}

func (s *SettingS) Copy() *SettingS {
	result := &SettingS{s.SettingId, make([]string, len(s.Names))}
	copy(result.Names, s.Names)
	return result
}

func (s *Setting) Copy() *Setting {
	result := &Setting{
		Inputs:  make([]ActionCId, len(s.Inputs)),
		Outputs: make([]TemplateId, len(s.Outputs)),
	}
	copy(result.Inputs, s.Inputs)
	copy(result.Outputs, s.Outputs)
	return result
}

func (s *Setting) AppendAction(a ActionCId) *Setting {
	s.Inputs = append(s.Inputs, a)
	return s
}

func (s *Setting) AppendTemplate(t TemplateId) *Setting {
	s.Outputs = append(s.Outputs, t)
	return s
}

func (s *SettingS) AppendAction(a ActionC) *SettingS {
	if len(s.Setting().Outputs) <= len(s.Setting().Inputs) {
		panic("appending action to something that already ends in an action!")
	}
	s.SettingId = s.SettingId.ExtendByAction(IdActionC(a))
	return s
}

func (s *SettingS) AppendTemplate(a TemplateId, names ...string) *SettingS {
	if len(s.Setting().Outputs) > len(s.Setting().Inputs) {
		panic("appending template to something that already ends in a template!")
	}
	if a.Template().Slots() != len(names) {
		panic(fmt.Sprintf("appending template %v with names %v", a.Template(), names))
	}
	s.SettingId = s.SettingId.ExtendByTemplate(a)
	s.Names = append(s.Names, names...)
	return s
}

func (s *SettingT) SetLastChild(child *SettingT) {
	s.Children[len(s.Children)-1] = child
}

func (s *SettingT) LastChild() *SettingT {
	return s.Children[len(s.Children)-1]
}

func (s *SettingT) AppendTerm(t T) *SettingT {
	s.SettingId = s.SettingId.ExtendByTemplate(t.Head())
	s.Children = append(s.Children, nil)
	switch t := t.(type) {
	case *CompoundT:
		s.Args = append(s.Args, t.args...)
	}
	return s
}

func (s *SettingT) AppendAction(a ActionC) *SettingT {
	s.SettingId = s.SettingId.ExtendByAction(IdActionC(a))
	return s
}

func (s *SettingT) RemoveTerm() *SettingT {
	outputs := s.Setting().Outputs
	last := outputs[len(outputs)-1].Template()
	s.Args = s.Args[:len(s.Args)-last.Slots()]
	s.Children = s.Children[:len(s.Children)-1]
	s.SettingId = s.IdInit()
	return s
}

func (s *SettingT) RemoveAction() *SettingT {
	s.SettingId = s.IdInit()
	return s
}

func (s *SettingS) Lines() []string {
	result := make([]string, 0)
	inputs := s.Setting().Inputs
	outputs := s.Setting().Outputs
	used := 0
	for i, outputId := range outputs {
		output := outputId.Template()
		slots := output.Slots()
		result = append(result, fmt.Sprintf("%d> %s", i, output.ShowWith(s.Names[used:used+slots]...)))
		result = append(result, "")
		used = used + slots
		if i < len(inputs) {
			input := inputs[i].ActionC()
			result = append(result, fmt.Sprintf("%d< %v", i+1, input.Uninstantiate(s.Names)))
		}
	}
	return result
}
