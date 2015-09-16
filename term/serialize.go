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

type TemplateId intern.Id

func (t *Template) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for _, s := range t.Parts {
		x.Append(x.New().StoreStr(s))
	}
	x.StorePair(x.New().StoreInt(int(v1Template)), x)
	return x
}

func unpackTemplate(x intern.Packed) *Template {
	v, x := x.Pair()
	switch versionAndType(v.Int()) {
	case v1Template:
		parts := x.List()
		t := Template{make([]string, len(parts))}
		for i, part := range parts {
			t.Parts[i] = part.Str()
		}
		return &t
	default:
		panic("Unknown kind of template!")
	}
}

func (id TemplateId) Template() *Template {
	return unpackTemplate((*intern.Id)(&id))
}

func IdTemplate(t *Template) TemplateId {
	x := new(intern.Id)
	t.Pack(x)
	return TemplateId(*x)
}

func (t *Template) Id() TemplateId {
	return IdTemplate(t)
}

func SaveTemplate(t *Template) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadTemplate(b bson.M) *Template {
	return unpackTemplate(&intern.Record{b})
}

//T

type TId intern.Id

type kindT int

const (
	compoundT = kindT(iota)
	intT
	strT
	quotedT
	wrapperT
	chanT
)

func IdT(t T) TId {
	x := new(intern.Id)
	t.Pack(x)
	return TId(*x)
}

func (t *CompoundT) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for _, arg := range t.args {
		x.Append(arg.Pack(x.New()))
	}
	x.StorePair(t.Template().Pack(x.New()), x)
	x.StorePair(x.New().StoreInt(int(compoundT)), x)
	x.StorePair(x.New().StoreInt(int(v1T)), x)
	return x
}

func (t Channel) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(chanT)), t.Setting.Pack(x.New()))
	x.StorePair(x.New().StoreInt(int(v1T)), x)
	return x
}

func (t Str) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(strT)), x.New().StoreStr(string(t)))
	x.StorePair(x.New().StoreInt(int(v1T)), x)
	return x
}

func (t Int) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(intT)), x.New().StoreInt(int(t)))
	x.StorePair(x.New().StoreInt(int(v1T)), x)
	return x
}

func (t Quoted) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(quotedT)), t.Value.Pack(x.New()))
	x.StorePair(x.New().StoreInt(int(v1T)), x)
	return x
}

func (w Wrapper) Pack(x intern.Packed) intern.Packed {
	return x.StoreInt(int(wrapperT))
}

var Unwrapped = Make("an unknown value that was destroyed by serialization")

func unpackT(x intern.Packed) T {
	v, x := x.Pair()
	switch versionAndType(v.Int()) {
	case v1T:
		kind, val := x.Pair()
		switch kindT(kind.Int()) {
		case compoundT:
			head, packedArgs := val.Pair()
			args := packedArgs.List()
			result := &CompoundT{IdTemplate(unpackTemplate(head)), make([]T, len(args))}
			for i, arg := range args {
				result.args[i] = unpackT(arg)
			}
			return result
		case chanT:
			return MakeChannel(unpackSettingT(val))
		case intT:
			return Int(val.Int())
		case strT:
			return Str(val.Str())
		case quotedT:
			return Quoted{unpackT(val)}
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

func (id TId) T() T {
	return unpackT((*intern.Id)(&id))
}

func SaveT(t T) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadT(b interface{}) T {
	return unpackT(&intern.Record{b})
}

//C

type CId intern.Id

type kindC int

const (
	compoundC = kindC(iota)
	referenceC
	constantC
)

func (t *CompoundC) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for _, arg := range t.args {
		x.Append(arg.Pack(x.New()))
	}
	x.StorePair(t.Template().Pack(x.New()), x)
	x.StorePair(x.New().StoreInt(int(compoundC)), x)
	x.StorePair(x.New().StoreInt(int(v1C)), x)
	return x
}

func (r ReferenceC) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(referenceC)), x.New().StoreInt(r.Index))
	x.StorePair(x.New().StoreInt(int(v1C)), x)
	return x
}

func (c ConstC) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(constantC)), c.Val.Pack(x.New()))
	x.StorePair(x.New().StoreInt(int(v1C)), x)
	return x
}

func IdC(c C) CId {
	x := new(intern.Id)
	c.Pack(x)
	return CId(*x)
}

func (id CId) C() C {
	return unpackC((*intern.Id)(&id))
}

func unpackC(x intern.Packed) C {
	v, x := x.Pair()
	switch versionAndType(v.Int()) {
	case v1C:
		kind, val := x.Pair()
		switch kindC(kind.Int()) {
		case compoundC:
			head, packedArgs := val.Pair()
			args := packedArgs.List()
			result := &CompoundC{IdTemplate(unpackTemplate(head)), make([]C, len(args))}
			for i, arg := range args {
				result.args[i] = unpackC(arg)
			}
			return result
		case referenceC:
			return ReferenceC{val.Int()}
		case constantC:
			return ConstC{unpackT(val)}
		default:
			panic("unknown kind of C")
		}
	default:
		panic("unknown version or wrong data type")
	}
}

func SaveC(t C) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadC(b interface{}) C {
	return unpackC(&intern.Record{b})
}

//Actions

type ActionCId intern.Id

func (a ActionC) Pack(x intern.Packed) intern.Packed {
	packedIntArgs := x.New().Empty()
	for _, arg := range a.IntArgs {
		packedIntArgs.Append(x.New().StoreInt(arg))
	}
	packedArgs := x.New().Empty()
	for _, arg := range a.Args {
		packedArgs.Append(arg.Pack(x.New()))
	}
	allArgs := x.New().StorePair(packedIntArgs, packedArgs)
	act := x.New().StoreInt(int(a.Act))
	x.StorePair(act, allArgs)
	version := x.New().StoreInt(int(v1ActionC))
	x.StorePair(version, x)
	return x
}

func unpackActionC(x intern.Packed) ActionC {
	version, x := x.Pair()
	switch versionAndType(version.Int()) {
	case v1ActionC:
		act, allArgs := x.Pair()
		packedIntArgs, packedArgs := allArgs.Pair()
		args := packedArgs.List()
		intArgs := packedIntArgs.List()
		result := ActionC{
			Act:     Action(act.Int()),
			Args:    make([]C, len(args)),
			IntArgs: make([]int, len(intArgs)),
		}
		for i, arg := range args {
			result.Args[i] = unpackC(arg)
		}
		for i, arg := range intArgs {
			result.IntArgs[i] = arg.Int()
		}
		return result
	default:
		panic("unknown version or wrong type")
	}
}

func IdActionC(a ActionC) ActionCId {
	x := new(intern.Id)
	a.Pack(x)
	return ActionCId(*x)
}

func (id ActionCId) ActionC() ActionC {
	return unpackActionC((*intern.Id)(&id))
}

func SaveActionC(t ActionC) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadActionC(b interface{}) ActionC {
	return unpackActionC(&intern.Record{b})
}

//Setting Line

func unpackSettingLine(x intern.Packed) (SettingLine, int) {
	v, _ := x.Pair()
	switch versionAndType(v.Int()) {
	case v1Template:
		template := unpackTemplate(x)
		return template.Id(), template.Slots()
	case v1ActionC:
		return unpackActionC(x).Id(), 0
	default:
		panic("unknown version or bad data type")
	}
}

//SettingT

type SettingTId intern.Id

func unpackSettingT(x intern.Packed) *SettingT {
	v, x := x.Pair()
	switch versionAndType(v.Int()) {
	case v1SettingT:
		setting, args := x.Pair()
		result := &SettingT{}
		result.Setting = unpackSetting(setting)
		for _, arg := range args.List() {
			result.Args = append(result.Args, unpackT(arg))
		}
		return result
	default:
		panic("unknown version or bad data type")
	}
}

func (s *SettingT) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for _, arg := range s.Args {
		x.Append(arg.Pack(x.New()))
	}
	x.StorePair(s.Setting.Pack(x.New()), x)
	x.StorePair(x.New().StoreInt(int(v1SettingT)), x)
	return x
}

//Setting

type SettingId intern.Id

func unpackSetting(x intern.Packed) *Setting {
	packedLines := x.List()
	result := Init()
	for _, packedLine := range packedLines {
		line, slots := unpackSettingLine(packedLine)
		result = result.Append(line, slots)
	}
	return result
}

func (s *Setting) Pack(x intern.Packed) intern.Packed {
	if s.Size == 0 {
		return x.Empty()
	}
	s.Previous.Pack(x)
	x.Append(s.Last.Pack(x.New()))
	return x
}

func (id SettingId) Setting() *Setting {
	return unpackSetting((*intern.Id)(&id))
}

func IdSetting(s *Setting) SettingId {
	x := new(intern.Id)
	s.Pack(x)
	return SettingId(*x)
}

func SaveSetting(t *Setting) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadSetting(b interface{}) *Setting {
	return unpackSetting(&intern.Record{b})
}
