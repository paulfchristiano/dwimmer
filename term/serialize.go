package term

import (
	"github.com/paulfchristiano/dwimmer/term/intern"
	"gopkg.in/mgo.v2/bson"
)

//TODO the old Id machinery seems awkward, probably some easier way to do it using Packers?

//Template

type TemplateId intern.Id

func (t *Template) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for _, s := range t.Parts {
		x.Append(x.New().StoreStr(s))
	}
	return x
}

func unpackTemplate(x intern.Packed) *Template {
	parts := x.List()
	t := Template{make([]string, len(parts))}
	for i, part := range parts {
		t.Parts[i] = part.Str()
	}
	return &t
}

func (id TemplateId) Template() *Template {
	return unpackTemplate((*intern.Id)(&id))
}

func IdTemplate(t *Template) TemplateId {
	x := new(intern.Id)
	t.Pack(x)
	return TemplateId(*x)
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
	return x.StorePair(x.New().StoreInt(int(compoundT)), x)
}

func (t Str) Pack(x intern.Packed) intern.Packed {
	return x.StorePair(x.New().StoreInt(int(strT)), x.New().StoreStr(string(t)))
}

func (t Int) Pack(x intern.Packed) intern.Packed {
	return x.StorePair(x.New().StoreInt(int(intT)), x.New().StoreInt(int(t)))
}

func (t Quoted) Pack(x intern.Packed) intern.Packed {
	return x.StorePair(x.New().StoreInt(int(quotedT)), t.Value.Pack(x.New()))
}

func (w Wrapper) Pack(x intern.Packed) intern.Packed {
	return x.StoreInt(int(wrapperT))
}

var Unwrapped = Make("an unknown value that was destroyed by serialization")

func unpackT(x intern.Packed) T {
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
	return x
}

func (r ReferenceC) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(referenceC)), x.New().StoreInt(r.Index))
	return x
}

func (c ConstC) Pack(x intern.Packed) intern.Packed {
	x.StorePair(x.New().StoreInt(int(constantC)), c.Val.Pack(x.New()))
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
		panic("unknown kind of T")
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

//Setting

type SettingId intern.Id

func unpackSetting(x intern.Packed) *Setting {
	result := Setting{make([]TemplateId, 0), make([]ActionCId, 0)}
	args := x.List()
	for i := 0; 2*i < len(args); i++ {
		result.Outputs = append(result.Outputs, IdTemplate(unpackTemplate(args[2*i])))
		if 2*i+1 < len(args) {
			result.Inputs = append(result.Inputs, IdActionC(unpackActionC(args[2*i+1])))
		}
	}
	return &result
}

func (s *Setting) Pack(x intern.Packed) intern.Packed {
	x.Empty()
	for i, output := range s.Outputs {
		x.Append(output.Template().Pack(x.New()))
		if i < len(s.Inputs) {
			x.Append(s.Inputs[i].ActionC().Pack(x.New()))
		}
	}
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

func (settingId SettingId) IdLast() TemplateId {
	resultId := (*intern.Id)(&settingId).Last().(*intern.Id)
	return TemplateId(*resultId)
}

func (settingId SettingId) IdInit() SettingId {
	resultId := (*intern.Id)(&settingId).Init().(*intern.Id)
	return SettingId(*resultId)
}

func (settingId SettingId) ExtendByAction(actionId ActionCId) SettingId {
	result := (*intern.Id)(&settingId)
	result.Append((*intern.Id)(&actionId))
	return SettingId(*result)
}

func (settingId SettingId) ExtendByTemplate(templateId TemplateId) SettingId {
	result := (*intern.Id)(&settingId)
	result.Append((*intern.Id)(&templateId))
	return SettingId(*result)
}

func SaveSetting(t *Setting) interface{} {
	x := new(intern.Record)
	t.Pack(x)
	return x.Value
}

func LoadSetting(b interface{}) *Setting {
	return unpackSetting(&intern.Record{b})
}

//Actions

type ActionCId intern.Id
type Version int

const (
	v1 = Version(-27)
)

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
	version := x.New().StoreInt(int(v1))
	return x.StorePair(version, x)
}

func unpackActionC(x intern.Packed) ActionC {
	version, x := x.Pair()
	switch Version(version.Int()) {
	case v1:
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
		//This is the old implementation prior to versioning...
		act, packedArgs := version, x
		args := packedArgs.List()
		result := ActionC{
			Act:     Action(act.Int()),
			Args:    make([]C, len(args)),
			IntArgs: make([]int, 0),
		}
		for i, arg := range args {
			result.Args[i] = unpackC(arg)
		}
		switch result.Act {
		case Clarify:
			result.Act = Return
		case Replace:
			result.IntArgs = []int{-1}
		}
		return result
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
