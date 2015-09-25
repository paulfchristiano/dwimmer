package term

import (
	"errors"
	"fmt"

	"github.com/paulfchristiano/dwimmer/storage/database"
	"github.com/paulfchristiano/dwimmer/term/intern"
)

const (
	v1Template = iota
	v1Setting
	v1ActionC
	v1T
	v1C
	v1SettingT
	v1TemplateID
	v1ActionCID
)

var (
	IDer     = intern.NewIDer()
	Recorder = intern.NewRecorder(database.Collection("records"), IDer)
)

func pair(packer intern.Packer, a, b interface{}) intern.Packed {
	args := make([]intern.Packed, 2)
	for i, x := range []interface{}{a, b} {
		args[i] = pack(packer, x)
	}
	return packer.PackPair(args[0], args[1])
}

func pack(packer intern.Packer, x interface{}) intern.Packed {
	switch x := x.(type) {
	case intern.Packed:
		return x
	case int:
		return packer.PackInt(x)
	case string:
		return packer.PackString(x)
	case intern.Pickler:
		return packer.PackPickler(x)
	case []int:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []string:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []intern.Pickler:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []T:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []C:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []intern.Packed:
		args := make([]interface{}, len(x))
		for i, y := range x {
			args[i] = interface{}(y)
		}
		return p(packer, args...)
	case []interface{}:
		return p(packer, x...)
	}
	panic(fmt.Errorf("bad argument %v to p(...) (type %T)", x, x))

}

func p(packer intern.Packer, xs ...interface{}) intern.Packed {
	args := make([]intern.Packed, len(xs))
	for i, x := range xs {
		args[i] = pack(packer, x)
	}
	return packer.PackList(args)
}

//Template---

type TemplateID intern.ID

func (t *Template) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(*Template)
	return ok
}

func (t *Template) Key() interface{} {
	return t
}

func (t *Template) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1Template, t.Parts)
}

func (t *Template) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	result := &Template{}
	if v, parts, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1Template {
			if parts, ok := packer.UnpackList(parts); ok {
				result.Parts = make([]string, len(parts))
				for i, part := range parts {
					if part, ok := packer.UnpackString(part); ok {
						result.Parts[i] = part
					}
				}
				return result, true
			}
		}
	}
	return nil, false
}

func ToTemplate(id TemplateID) (*Template, bool) {
	unpacked, ok := IDer.UnpackPickler(intern.ID(id), &Template{})
	if !ok {
		return nil, false
	}
	return unpacked.(*Template), ok
}

func IDTemplate(t *Template) TemplateID {
	return TemplateID(IDer.PackPickler(t).(intern.ID))
}

func SaveTemplate(t *Template) interface{} {
	return Recorder.PackPickler(t).(intern.Record).Value
}

func LoadTemplate(b interface{}) (*Template, bool) {
	result, ok := Recorder.UnpackPickler(intern.Record{b}, &Template{})
	if !ok {
		return nil, false
	}
	return result.(*Template), ok
}

func (id TemplateID) Template() *Template {
	result, ok := ToTemplate(id)
	if !ok {
		panic(errors.New("failed to convert ID to template!"))
	}
	return result
}

func (t *Template) ID() TemplateID {
	return IDTemplate(t)
}

//TemplateID---

func (t TemplateID) Key() interface{} { return t }

func (t TemplateID) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(TemplateID)
	return ok
}

func (t TemplateID) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1TemplateID, t.Template())
}

func (t TemplateID) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, template, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1TemplateID {
			if template, ok := packer.UnpackPickler(template, &Template{}); ok {
				return template.(*Template).ID(), true
			}
		}
	}
	return nil, false
}

//ActionCID---

func (a ActionCID) Key() interface{} { return a }

func (a ActionCID) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(ActionCID)
	return ok
}

func (a ActionCID) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1ActionCID, a.ActionC())
}

func (a ActionCID) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, action, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1ActionCID {
			if action, ok := packer.UnpackPickler(action, ActionC{}); ok {
				return action.(ActionC).ID(), true
			}
		}
	}
	return nil, false
}

//T---

type TID intern.ID

const (
	compoundT = iota
	intT
	strT
	quotedT
	wrapperT
	chanT
)

func (t Channel) Key() interface{} { return t }
func (t Channel) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}
func (t Quoted) Key() interface{} { return t }
func (t Quoted) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}
func (t Wrapper) Key() interface{} { return t }
func (t Wrapper) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}
func (t Str) Key() interface{} { return t }
func (t Str) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}
func (t Int) Key() interface{} { return t }
func (t Int) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}
func (t *CompoundT) Key() interface{} { return t }
func (t *CompoundT) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(T)
	return ok
}

func (t Channel) Unpickle(packer intern.Packer, b intern.Packed) (intern.Pickler, bool) {
	return (&CompoundT{}).Unpickle(packer, b)
}
func (t Quoted) Unpickle(packer intern.Packer, b intern.Packed) (intern.Pickler, bool) {
	return (&CompoundT{}).Unpickle(packer, b)
}
func (t Wrapper) Unpickle(packer intern.Packer, b intern.Packed) (intern.Pickler, bool) {
	return (&CompoundT{}).Unpickle(packer, b)
}
func (t Str) Unpickle(packer intern.Packer, b intern.Packed) (intern.Pickler, bool) {
	return (&CompoundT{}).Unpickle(packer, b)
}
func (t Int) Unpickle(packer intern.Packer, b intern.Packed) (intern.Pickler, bool) {
	return (&CompoundT{}).Unpickle(packer, b)
}

func (t *CompoundT) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, val, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1T {
			if kind, val, ok := packer.UnpackPair(val); ok {
				if kind, ok := packer.UnpackInt(kind); ok {
					return unpickleT(packer, kind, val)
				}
			}
			//panic(fmt.Errorf("failed to unpack pair from %v", val))
		}
	}
	//panic(fmt.Errorf("faield to unpack from pickled %v", pickled))
	return nil, false
}

func unpickleT(packer intern.Packer, kind int, packed intern.Packed) (T, bool) {
	switch kind {
	case chanT:
		unpickled, ok := packer.UnpackPickler(packed, &SettingT{})
		if !ok {
			return nil, false
		}
		setting, ok := unpickled.(*SettingT)
		return MakeChannel(setting), ok
	case intT:
		result, ok := packer.UnpackInt(packed)
		if !ok {
			//panic(fmt.Errorf("failed to unpickle int from %v", packed))
		}
		return Int(result), ok
	case strT:
		result, ok := packer.UnpackString(packed)
		return Str(result), ok
	case quotedT:
		unpickled, ok := packer.UnpackPickler(packed, &CompoundT{})
		if !ok {
			return nil, false
		}
		return Quote(unpickled.(T)), true
	case wrapperT:
		return Make("a term that stands in for a wrapped go object that could not be pickled").T(), true
	case compoundT:
		result := &CompoundT{}
		if packed, args, ok := packer.UnpackPair(packed); ok {
			if args, ok := packer.UnpackList(args); ok {
				result.args = make([]T, len(args))
				for i, arg := range args {
					unpacked, ok := packer.UnpackPickler(arg, &CompoundT{})
					if !ok {
						return nil, false
					}
					result.args[i] = unpacked.(T)
				}
				id, ok := unpickleTemplateID(packer, packed)
				if !ok {
					return nil, false
				}
				result.TemplateID = id
				return result, true
			}
		}
	}
	return nil, false
}

func (t *CompoundT) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1T, pair(packer, compoundT, pair(packer, t.Head(), t.args)))
}
func (t Channel) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1T, pair(packer, chanT, t.Setting))
}
func (t Int) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1T, pair(packer, intT, int(t)))
}
func (t Str) Pickle(packer intern.Packer) intern.Packed {
	result := pair(packer, v1T, pair(packer, strT, string(t)))
	return result
}
func (t Wrapper) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1T, pair(packer, wrapperT, 0))
}
func (t Quoted) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1T, pair(packer, quotedT, t.Value))
}

func SaveT(t T) interface{} {
	return Recorder.PackPickler(t).(intern.Record).Value
}

func LoadT(b interface{}) (t T, ok bool) {
	result, ok := Recorder.UnpackPickler(intern.Record{b}, &CompoundT{})
	if !ok {
		return nil, false
	}
	return result.(T), ok
}

func (id TID) T() T {
	result, ok := ToT(id)
	if !ok {
		//panic(errors.New("failed to convert ID to term!"))
	}
	return result
}

func ToT(id TID) (T, bool) {
	unpacked, ok := IDer.UnpackPickler(intern.ID(id), &CompoundT{})
	if !ok {
		return nil, false
	}
	return unpacked.(T), ok
}

func IDT(t T) TID {
	return TID(IDer.PackPickler(t).(intern.ID))
}

func (t *CompoundT) ID() TID { return IDT(t) }
func (t Int) ID() TID        { return IDT(t) }
func (t Str) ID() TID        { return IDT(t) }
func (t Channel) ID() TID    { return IDT(t) }
func (t Wrapper) ID() TID    { return IDT(t) }
func (t Quoted) ID() TID     { return IDT(t) }

//C---

type CID intern.ID

const (
	compoundC = iota
	referenceC
	constantC
)

func (c *CompoundC) Key() interface{} { return c }
func (c *CompoundC) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(C)
	return ok
}
func (c ReferenceC) Key() interface{} { return c }
func (c ReferenceC) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(C)
	return ok
}
func (c ConstC) Key() interface{} { return c }
func (c ConstC) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(C)
	return ok
}

func (c *CompoundC) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1C, pair(packer, compoundC, pair(packer, c.TemplateID, c.args)))
}

func (c ReferenceC) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1C, pair(packer, referenceC, int(c.Index)))
}

func (c ConstC) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1C, pair(packer, constantC, c.Val))
}

func (c ReferenceC) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	return (&CompoundC{}).Unpickle(packer, pickled)
}

func (c ConstC) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	return (&CompoundC{}).Unpickle(packer, pickled)
}

func (c *CompoundC) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, val, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1C {
			if kind, val, ok := packer.UnpackPair(val); ok {
				if kind, ok := packer.UnpackInt(kind); ok {
					return unpickleC(packer, kind, val)
				}

			}
		}
	}
	return nil, false
}

func unpickleC(packer intern.Packer, kind int, packed intern.Packed) (C, bool) {
	switch kind {
	case compoundC:
		result := &CompoundC{}
		if packed, args, ok := packer.UnpackPair(packed); ok {
			if args, ok := packer.UnpackList(args); ok {
				result.args = make([]C, len(args))
				for i, arg := range args {
					unpacked, ok := packer.UnpackPickler(arg, &CompoundC{})
					if !ok {
						return nil, false
					}
					result.args[i] = unpacked.(C)
				}
				id, ok := unpickleTemplateID(packer, packed)
				if !ok {
					return nil, false
				}
				result.TemplateID = id
				return result, true
			}
		}
	case referenceC:
		n, ok := packer.UnpackInt(packed)
		return ReferenceC{n}, ok
	case constantC:
		unpacked, ok := packer.UnpackPickler(packed, &CompoundT{})
		if !ok {
			return nil, false
		}
		return ConstC{unpacked.(T)}, ok
	}
	return nil, false
}

func unpickleTemplateID(packer intern.Packer, packed intern.Packed) (TemplateID, bool) {
	unpickled, ok := packer.UnpackPickler(packed, TemplateID(0))
	if !ok {
		unpickled, ok = packer.UnpackPickler(packed, &Template{})
		if !ok {
			return TemplateID(0), false
		}
		unpickled = unpickled.(*Template).ID()
	}
	return unpickled.(TemplateID), true
}

func SaveC(c C) interface{} {
	return Recorder.PackPickler(c).(intern.Record).Value
}

func LoadC(b interface{}) (C, bool) {
	result, ok := Recorder.UnpackPickler(intern.Record{b}, &CompoundC{})
	return result.(C), ok
}

func IDC(c C) CID {
	return CID(IDer.PackPickler(c).(intern.ID))
}

func ToC(id CID) (C, bool) {
	unpacked, ok := IDer.UnpackPickler(intern.ID(id), &CompoundC{})
	return unpacked.(C), ok
}

func (id CID) C() C {
	result, ok := ToC(id)
	if !ok {
		panic(errors.New("unable to convert ID to C"))
	}
	return result
}

func (c *CompoundC) ID() CID {
	return IDC(c)
}
func (c ConstC) ID() CID {
	return IDC(c)
}
func (c ReferenceC) ID() CID {
	return IDC(c)
}

//Actions---

type ActionCID intern.ID

func (a ActionC) Key() interface{} { return nil }
func (a ActionC) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(ActionC)
	return ok
}

func (a ActionC) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1ActionC, pair(packer, int(a.Act), pair(packer, a.IntArgs, a.Args)))
}

func (a ActionC) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, val, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1ActionC {
			if act, allArgs, ok := packer.UnpackPair(val); ok {
				result := ActionC{}
				if act, ok := packer.UnpackInt(act); ok {
					result.Act = Action(act)
				} else {
					return nil, false
				}
				if intArgs, args, ok := packer.UnpackPair(allArgs); ok {
					if intArgs, ok := packer.UnpackList(intArgs); ok {
						result.IntArgs = make([]int, len(intArgs))
						for i, packedArg := range intArgs {
							arg, ok := packer.UnpackInt(packedArg)
							if !ok {
								panic(fmt.Errorf("failed to unpack int! %v", val))
								return nil, false
							}
							result.IntArgs[i] = arg
						}
					} else {
						panic(fmt.Errorf("failed to unpack list of int args! %v", val))
						return nil, false
					}
					if args, ok := packer.UnpackList(args); ok {
						result.Args = make([]C, len(args))
						for i, arg := range args {
							unpacked, ok := packer.UnpackPickler(arg, &CompoundC{})
							if !ok {
								return nil, false
							}
							result.Args[i] = unpacked.(C)
						}
					} else {
						return nil, false
					}
					return result, true
				}
			}
		}
	}
	return nil, false
}

func IDActionC(a ActionC) ActionCID {
	return ActionCID(IDer.PackPickler(a).(intern.ID))
}

func ToActionC(id ActionCID) (ActionC, bool) {
	result, ok := IDer.UnpackPickler(intern.ID(id), ActionC{})
	if !ok {
		return ActionC{}, false
	}
	return result.(ActionC), true
}

func SaveActionC(a ActionC) interface{} {
	return Recorder.PackPickler(a).(intern.Record).Value
}

func LoadActionC(b interface{}) (act ActionC, ok bool) {
	result, ok := Recorder.UnpackPickler(intern.Record{b}, ActionC{})
	if !ok {
		return ActionC{}, false
	}
	return result.(ActionC), ok
}

func (id ActionCID) ActionC() ActionC {
	result, ok := ToActionC(id)
	if !ok {
		panic(errors.New("can't convert ID to ActionC"))
	}
	return result
}

func (action ActionC) ID() ActionCID {
	return IDActionC(action)
}

//Setting Line---

func unpickleSettingLine(packer intern.Packer, x intern.Packed) (SettingLine, int, bool) {
	result, ok := packer.UnpackPickler(x, ActionCID(0))
	if ok {
		return result.(ActionCID), 0, true
	}
	id, ok := unpickleTemplateID(packer, x)
	if ok {
		return id, id.Template().Slots(), true
	}
	result, ok = packer.UnpackPickler(x, ActionC{})
	if ok {
		return result.(ActionC).ID(), 0, true
	}
	return nil, 0, false
}

//SettingT---

func (s *SettingT) Key() interface{} { return s }
func (s *SettingT) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(*SettingT)
	return ok
}

func (s *SettingT) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	if v, val, ok := packer.UnpackPair(pickled); ok {
		if v, ok := packer.UnpackInt(v); ok && v == v1SettingT {
			if setting, packedArgs, ok := packer.UnpackPair(val); ok {
				result := &SettingT{}
				unpacked, ok := packer.UnpackPickler(setting, &Setting{})
				if !ok {
					return nil, false
				}
				result.Setting = unpacked.(*Setting)
				args, ok := packer.UnpackList(packedArgs)
				if !ok {
					return nil, false
				}
				result.Args = make([]T, len(args))
				for i, arg := range args {
					unpacked, ok := packer.UnpackPickler(arg, &CompoundT{})
					if !ok {
						return nil, false
					}
					result.Args[i] = unpacked.(T)
				}
				return result, true
			}
		}
	}
	return nil, false
}

func (s *SettingT) Pickle(packer intern.Packer) intern.Packed {
	return pair(packer, v1SettingT, pair(packer, s.Setting, s.Args))
}

//Setting---

type SettingID intern.ID

func (s *Setting) Key() interface{} { return s }
func (s *Setting) Test(pickler intern.Pickler) bool {
	_, ok := pickler.(*Setting)
	return ok
}

func (s *Setting) Unpickle(packer intern.Packer, pickled intern.Packed) (intern.Pickler, bool) {
	lines, ok := packer.UnpackList(pickled)
	if !ok {
		return nil, false
	}
	result := Init()
	for _, packedLine := range lines {
		line, slots, ok := unpickleSettingLine(packer, packedLine)
		if !ok {
			return nil, false
		}
		result = result.Append(line, slots)
	}
	return result, true
}

func (s *Setting) Pickle(packer intern.Packer) intern.Packed {
	if s.Size == 0 {
		return packer.PackList([]intern.Packed{})
	}
	previous := s.Previous.Pickle(packer)
	return packer.AppendToPacked(previous, packer.PackPickler(s.Last))
}

func ToSetting(id SettingID) (*Setting, bool) {
	result, ok := IDer.UnpackPickler(intern.ID(id), &Setting{})
	if !ok {
		return nil, false
	}
	return result.(*Setting), true
}

func IDSetting(s *Setting) SettingID {
	return SettingID(IDer.PackPickler(s).(intern.ID))
}

func SaveSetting(s *Setting) interface{} {
	return Recorder.PackPickler(s).(intern.Record).Value
}

func LoadSetting(b interface{}) (s *Setting, ok bool) {
	result, ok := Recorder.UnpackPickler(intern.Record{b}, &Setting{})
	if !ok {
		return nil, false
	}
	return result.(*Setting), true
}

func (id SettingID) Setting() *Setting {
	result, ok := ToSetting(id)
	if !ok {
		panic(errors.New("failed to convert ID to setting"))
	}
	return result
}
