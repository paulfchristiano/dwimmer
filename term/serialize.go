package term

import (
	"github.com/paulfchristiano/dwimmer/term/intern"
	"gopkg.in/mgo.v2/bson"
)

type versionAndType int

const (
	v1Template = versionAndType(iota)
	v1Setting
	v1ActionC
	v1T
	v1C
	v1SettingT
)

//Template

type TemplateID intern.ID

func packTemplate(p intern.Packer, t *Template) intern.Packed {
	parts := make([]intern.Packed, len(t.Parts))
	for i, part := range t.Parts {
		parts[i] = p.PackString(part)
	}
	return p.PackPair(p.PackInt(int(v1Template)), p.PackList(parts))
}

func unpackTemplate(p intern.Packer, x intern.Packed) *Template {
	v, x := p.UnpackPair(x)
	switch versionAndType(p.UnpackInt(v)) {
	case v1Template:
		parts := p.UnpackList(x)
		t := Template{make([]string, len(parts))}
		for i, part := range parts {
			t.Parts[i] = p.UnpackString(part)
		}
		return &t
	default:
		panic("Unknown kind of template!")
	}
}

func ToTemplate(ider intern.Packer, id TemplateID) *Template {
	return unpackTemplate(ider, intern.ID(id))
}

func IDTemplate(ider intern.Packer, t *Template) TemplateID {
	return TemplateID(packTemplate(ider, t).(intern.ID))
}

func SaveTemplate(t *Template) interface{} {
	return packTemplate(intern.Recorder{}, t).(intern.Record).Value
}

func LoadTemplate(b bson.M) *Template {
	return unpackTemplate(intern.Recorder{}, intern.Record{b})
}

//T

type TID intern.ID

type kindT int

const (
	compoundT = kindT(iota)
	intT
	strT
	quotedT
	wrapperT
	chanT
)

func IDT(ider intern.Packer, t T) TID {
	return TID(packT(ider, ider, t).(intern.ID))
}

func packT(packer, ider intern.Packer, t T) intern.Packed {
	var raw intern.Packed
	var kind kindT
	switch t := t.(type) {
	case *CompoundT:
		args := make([]intern.Packed, len(t.args))
		for i, arg := range t.args {
			args[i] = packT(packer, ider, arg)
		}
		raw = packer.PackPair(
			packTemplate(packer, ToTemplate(ider, t.Head())),
			packer.PackList(args),
		)
		kind = compoundT
	case Channel:
		raw = packSettingT(packer, ider, t.Setting)
		kind = chanT
	case Str:
		raw = packer.PackString(string(t))
		kind = strT
	case Int:
		raw = packer.PackInt(int(t))
		kind = intT
	case Quoted:
		raw = packT(packer, ider, t.Value)
		kind = quotedT
	case Wrapper:
		raw = packer.PackInt(-1)
		kind = wrapperT
	default:
		panic("unknown kind of T")
	}
	typed := packer.PackPair(packer.PackInt(int(kind)), raw)
	versioned := packer.PackPair(packer.PackInt(int(v1T)), typed)
	return versioned
}

var Unwrapped = Make("an unknown value that was destroyed by serialization")

func unpackT(packer, ider intern.Packer, x intern.Packed) T {
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1T:
		kind, val := packer.UnpackPair(x)
		switch kindT(packer.UnpackInt(kind)) {
		case compoundT:
			head, packedArgs := packer.UnpackPair(val)
			args := packer.UnpackList(packedArgs)
			result := &CompoundT{
				IDTemplate(ider, unpackTemplate(packer, head)),
				make([]T, len(args)),
			}
			for i, arg := range args {
				result.args[i] = unpackT(packer, ider, arg)
			}
			return result
		case chanT:
			return MakeChannel(unpackSettingT(packer, ider, val))
		case intT:
			return Int(packer.UnpackInt(val))
		case strT:
			return Str(packer.UnpackString(val))
		case quotedT:
			return Quoted{unpackT(packer, ider, val)}
		case wrapperT:
			panic("unwrapping the ununwrappable")
			return Unwrapped.T()
		default:
			panic("unknown kind of T")
		}
	default:
		panic("unknown version or wrong type of data")
	}
}

func ToT(ider intern.Packer, id TID) T {
	return unpackT(ider, ider, intern.ID(id))
}

func SaveT(ider intern.Packer, t T) interface{} {
	return packT(intern.Recorder{}, ider, t).(intern.Record).Value
}

func LoadT(ider intern.Packer, b interface{}) T {
	return unpackT(intern.Recorder{}, ider, intern.Record{b})
}

//C

type CID intern.ID

type kindC int

const (
	compoundC = kindC(iota)
	referenceC
	constantC
)

func packC(packer, ider intern.Packer, c C) intern.Packed {
	var kind kindC
	var raw intern.Packed
	switch c := c.(type) {
	case *CompoundC:
		args := make([]intern.Packed, len(c.args))
		for i, arg := range c.args {
			args[i] = packC(packer, ider, arg)
		}
		raw = packer.PackPair(packTemplate(packer, ToTemplate(ider, c.Head())), packer.PackList(args))
		kind = compoundC
	case ReferenceC:
		raw = packer.PackInt(int(c.Index))
		kind = referenceC
	case ConstC:
		raw = packT(packer, ider, c.Val)
		kind = constantC
	default:
		panic("packing unknown type of C")
	}
	typed := packer.PackPair(packer.PackInt(int(kind)), raw)
	versioned := packer.PackPair(packer.PackInt(int(v1C)), typed)
	return versioned
}

func IDC(ider intern.Packer, c C) CID {
	return CID(packC(ider, ider, c).(intern.ID))
}

func ToC(ider intern.Packer, id CID) C {
	return unpackC(ider, ider, intern.ID(id))
}

func unpackC(packer, ider intern.Packer, x intern.Packed) C {
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1C:
		kind, val := packer.UnpackPair(x)
		switch kindC(packer.UnpackInt(kind)) {
		case compoundC:
			head, packedArgs := packer.UnpackPair(val)
			args := packer.UnpackList(packedArgs)
			result := &CompoundC{
				IDTemplate(ider, unpackTemplate(packer, head)),
				make([]C, len(args)),
			}
			for i, arg := range args {
				result.args[i] = unpackC(packer, ider, arg)
			}
			return result
		case referenceC:
			return ReferenceC{packer.UnpackInt(val)}
		case constantC:
			return ConstC{unpackT(packer, ider, val)}
		default:
			panic("unknown kind of C")
		}
	default:
		panic("unknown version or wrong type of data")
	}
}

func SaveC(ider intern.Packer, c C) interface{} {
	return packC(intern.Recorder{}, ider, c).(intern.Record).Value
}

func LoadC(ider intern.Packer, b interface{}) C {
	return unpackC(intern.Recorder{}, ider, intern.Record{b})
}

//Actions

type ActionCID intern.ID

func packActionC(packer, ider intern.Packer, a ActionC) intern.Packed {
	intArgs := make([]intern.Packed, len(a.IntArgs))
	for i, arg := range a.IntArgs {
		intArgs[i] = packer.PackInt(arg)
	}
	packedIntArgs := packer.PackList(intArgs)
	args := make([]intern.Packed, len(a.Args))
	for i, arg := range a.Args {
		args[i] = packC(packer, ider, arg)
	}
	packedArgs := packer.PackList(args)
	allArgs := packer.PackPair(packedIntArgs, packedArgs)
	act := packer.PackInt(int(a.Act))
	raw := packer.PackPair(act, allArgs)
	versioned := packer.PackPair(packer.PackInt(int(v1ActionC)), raw)
	return versioned
}

func unpackActionC(packer, ider intern.Packer, x intern.Packed) ActionC {
	version, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(version)) {
	case v1ActionC:
		act, allArgs := packer.UnpackPair(x)
		packedIntArgs, packedArgs := packer.UnpackPair(allArgs)
		args := packer.UnpackList(packedArgs)
		intArgs := packer.UnpackList(packedIntArgs)
		result := ActionC{
			Act:     Action(packer.UnpackInt(act)),
			Args:    make([]C, len(args)),
			IntArgs: make([]int, len(intArgs)),
		}
		for i, arg := range args {
			result.Args[i] = unpackC(packer, ider, arg)
		}
		for i, arg := range intArgs {
			result.IntArgs[i] = packer.UnpackInt(arg)
		}
		return result
	default:
		panic("unknown version or wrong type")
	}
}

func IDActionC(ider intern.Packer, a ActionC) ActionCID {
	return ActionCID(packActionC(ider, ider, a).(intern.ID))
}

func ToActionC(ider intern.Packer, id ActionCID) ActionC {
	return unpackActionC(ider, ider, intern.ID(id))
}

func SaveActionC(ider intern.Packer, a ActionC) interface{} {
	return packActionC(intern.Recorder{}, ider, a).(intern.Record).Value
}

func LoadActionC(ider intern.Packer, b interface{}) ActionC {
	return unpackActionC(intern.Recorder{}, ider, intern.Record{b})
}

//Setting Line

func unpackSettingLine(packer, ider intern.Packer, x intern.Packed) (SettingLine, int) {
	v, _ := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1Template:
		template := unpackTemplate(packer, x)
		return IDTemplate(ider, template), template.Slots()
	case v1ActionC:
		return IDActionC(ider, unpackActionC(packer, ider, x)), 0
	default:
		panic("unknown version or bad data type")
	}
}

func packSettingLine(packer, ider intern.Packer, l SettingLine) intern.Packed {
	switch l := l.(type) {
	case ActionCID:
		return packActionC(packer, ider, ToActionC(ider, l))
	case TemplateID:
		return packTemplate(packer, ToTemplate(ider, l))
	default:
		panic("unknown kind of setting line")
	}
}

//SettingT

func unpackSettingT(packer, ider intern.Packer, x intern.Packed) *SettingT {
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1SettingT:
		setting, args := packer.UnpackPair(x)
		result := &SettingT{}
		result.Setting = unpackSetting(packer, ider, setting)
		for _, arg := range packer.UnpackList(args) {
			result.Args = append(result.Args, unpackT(packer, ider, arg))
		}
		return result
	default:
		panic("unknown version or bad data type")
	}
}

func packSettingT(packer, ider intern.Packer, s *SettingT) intern.Packed {
	args := make([]intern.Packed, len(s.Args))
	for i, arg := range s.Args {
		args[i] = packT(packer, ider, arg)
	}
	raw := packer.PackPair(packSetting(packer, ider, s.Setting), packer.PackList(args))
	versioned := packer.PackPair(packer.PackInt(int(v1SettingT)), raw)
	return versioned
}

//Setting

type SettingID intern.ID

func unpackSetting(packer, ider intern.Packer, x intern.Packed) *Setting {
	packedLines := packer.UnpackList(x)
	result := Init(ider)
	for _, packedLine := range packedLines {
		line, slots := unpackSettingLine(packer, ider, packedLine)
		result = result.Append(ider, line, slots)
	}
	return result
}

func packSetting(packer, ider intern.Packer, s *Setting) intern.Packed {
	if s.Size == 0 {
		return packer.PackList([]intern.Packed{})
	}
	previous := packSetting(packer, ider, s.Previous)
	last := packSettingLine(packer, ider, s.Last)
	return packer.AppendToPacked(previous, last)
}

func ToSetting(ider intern.Packer, id SettingID) *Setting {
	return unpackSetting(ider, ider, intern.ID(id))
}

func IDSetting(ider intern.Packer, s *Setting) SettingID {
	return SettingID(packSetting(ider, ider, s).(intern.ID))
}

func SaveSetting(ider intern.Packer, s *Setting) interface{} {
	return packSetting(intern.Recorder{}, ider, s)
}

func LoadSetting(ider intern.Packer, b interface{}) *Setting {
	return unpackSetting(intern.Recorder{}, ider, intern.Record{b})
}
