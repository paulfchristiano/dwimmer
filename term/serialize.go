package term

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/paulfchristiano/dwimmer/storage/database"
	"github.com/paulfchristiano/dwimmer/term/intern"
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

const (
	templateCache = iota
	tCache
	cCache
	settingTCache
)

func packTemplate(p intern.Packer, t *Template) (result intern.Packed) {
	var raw intern.Packed
	if cached, ok := p.GetCachedPack(templateCache, t); ok {
		raw = cached
	} else {
		parts := make([]intern.Packed, len(t.Parts))
		for i, part := range t.Parts {
			parts[i] = p.PackString(part)
		}
		raw = p.PackList(parts)
		raw = p.CachePack(templateCache, t, raw)
	}
	return p.PackPair(p.PackInt(int(v1Template)), raw)
}

func unpackTemplate(p intern.Packer, x intern.Packed) (result *Template) {
	v, x := p.UnpackPair(x)
	switch versionAndType(p.UnpackInt(v)) {
	case v1Template:
		cached, ok := p.GetCachedUnpack(templateCache, x)
		if ok {
			result, ok = cached.(*Template)
			if !ok {
				panic(errors.New("bad cached value"))
			}
			return result
		}
		defer func() {
			if result != nil {
				p.CacheUnpack(templateCache, x, result)
			}
		}()
		parts := p.UnpackList(x)
		t := Template{make([]string, len(parts))}
		for i, part := range parts {
			t.Parts[i] = p.UnpackString(part)
		}
		return &t
	default:
		panic(errors.New("Unknown kind of template!"))
	}
}

var (
	IDer     = intern.NewIDer()
	Recorder = intern.NewRecorder(database.Collection("records"))
)

func ToTemplate(id TemplateID) *Template {
	return unpackTemplate(IDer, intern.ID(id))
}

func IDTemplate(t *Template) TemplateID {
	return TemplateID(packTemplate(IDer, t).(intern.ID))
}

func SaveTemplate(t *Template) interface{} {
	return intern.FromRecord(packTemplate(Recorder, t).(intern.Record))
}

func decodingRecover(recovered interface{}) error {
	if recovered != nil {
		if _, ok := recovered.(runtime.Error); ok {
			panic(recovered)
		}
		return recovered.(error)
	}
	return nil
}

func LoadTemplate(b interface{}) (result *Template, err error) {
	defer func() { err = decodingRecover(recover()) }()
	return unpackTemplate(Recorder, intern.MakeRecord(b)), nil
}

func (id TemplateID) Template() *Template {
	return ToTemplate(id)
}

func (t *Template) ID() TemplateID {
	return IDTemplate(t)
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

func IDT(t T) TID {
	return TID(packT(IDer, t).(intern.ID))
}

func packT(packer intern.Packer, t T) intern.Packed {
	var raw intern.Packed
	var kind kindT
	switch t := t.(type) {
	case *CompoundT:
		kind = compoundT
		if cached, ok := packer.GetCachedPack(tCache, t); ok {
			raw = cached
		} else {
			args := make([]intern.Packed, len(t.args))
			for i, arg := range t.args {
				args[i] = packT(packer, arg)
			}
			raw = packer.PackPair(
				packTemplate(packer, ToTemplate(t.Head())),
				packer.PackList(args),
			)
			raw = packer.CachePack(tCache, t, raw)
		}
	case Channel:
		raw = packSettingT(packer, t.Setting)
		kind = chanT
	case Str:
		raw = packer.PackString(string(t))
		kind = strT
	case Int:
		raw = packer.PackInt(int(t))
		kind = intT
	case Quoted:
		raw = packT(packer, t.Value)
		kind = quotedT
	case Wrapper:
		raw = packer.PackInt(-1)
		kind = wrapperT
	default:
		panic(errors.New("unknown kind of T"))
	}
	typed := packer.PackPair(packer.PackInt(int(kind)), raw)
	versioned := packer.PackPair(packer.PackInt(int(v1T)), typed)
	return versioned
}

var Unwrapped = Make("an unknown value that was destroyed by serialization")

func unpackT(packer intern.Packer, x intern.Packed) (result T) {
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1T:
		kind, val := packer.UnpackPair(x)
		switch kindT(packer.UnpackInt(kind)) {
		case compoundT:
			cached, ok := packer.GetCachedUnpack(tCache, val)
			if ok {
				result, ok = cached.(*CompoundT)
				if !ok {
					panic(fmt.Errorf("bad cached value: %v is not a CompoundT", cached))
				}
				return result
			}
			defer func() { packer.CacheUnpack(tCache, val, result) }()
			head, packedArgs := packer.UnpackPair(val)
			args := packer.UnpackList(packedArgs)
			resultArgs := make([]T, len(args))
			result = &CompoundT{
				IDTemplate(unpackTemplate(packer, head)),
				resultArgs,
			}
			for i, arg := range args {
				resultArgs[i] = unpackT(packer, arg)
			}
			return result
		case chanT:
			return MakeChannel(unpackSettingT(packer, val))
		case intT:
			return Int(packer.UnpackInt(val))
		case strT:
			return Str(packer.UnpackString(val))
		case quotedT:
			return Quoted{unpackT(packer, val)}
		case wrapperT:
			panic(errors.New("unwrapping the ununwrappable"))
			return Unwrapped.T()
		default:
			panic(errors.New("unknown kind of T"))
		}
	default:
		panic(fmt.Errorf("unknown version or wrong type of data %v", v))
	}
}

func ToT(id TID) T {
	return unpackT(IDer, intern.ID(id))
}

func SaveT(t T) interface{} {
	return intern.FromRecord(packT(Recorder, t).(intern.Record))
}

func LoadT(b interface{}) (t T, err error) {
	defer func() { err = decodingRecover(recover()) }()
	return unpackT(Recorder, intern.MakeRecord(b)), nil
}

func (id TID) T() T {
	return ToT(id)
}

func (t *CompoundT) ID() TID {
	return IDT(t)
}
func (t Int) ID() TID {
	return IDT(t)
}
func (t Str) ID() TID {
	return IDT(t)
}
func (t Channel) ID() TID {
	return IDT(t)
}
func (t Wrapper) ID() TID {
	return IDT(t)
}
func (t Quoted) ID() TID {
	return IDT(t)
}

//C

type CID intern.ID

type kindC int

const (
	compoundC = kindC(iota)
	referenceC
	constantC
)

func packC(packer intern.Packer, c C) intern.Packed {
	var kind kindC
	var raw intern.Packed
	switch c := c.(type) {
	case *CompoundC:
		kind = compoundC
		if cached, ok := packer.GetCachedPack(cCache, c); ok {
			raw = cached
		} else {
			args := make([]intern.Packed, len(c.args))
			for i, arg := range c.args {
				args[i] = packC(packer, arg)
			}
			raw = packer.PackPair(packTemplate(packer, ToTemplate(c.Head())), packer.PackList(args))
			raw = packer.CachePack(cCache, c, raw)
		}
	case ReferenceC:
		raw = packer.PackInt(int(c.Index))
		kind = referenceC
	case ConstC:
		raw = packT(packer, c.Val)
		kind = constantC
	default:
		panic(errors.New("packing unknown type of C"))
	}
	typed := packer.PackPair(packer.PackInt(int(kind)), raw)
	versioned := packer.PackPair(packer.PackInt(int(v1C)), typed)
	return versioned
}

func IDC(c C) CID {
	return CID(packC(IDer, c).(intern.ID))
}

func ToC(id CID) C {
	return unpackC(IDer, intern.ID(id))
}

func unpackC(packer intern.Packer, x intern.Packed) (result C) {
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1C:
		kind, val := packer.UnpackPair(x)
		switch kindC(packer.UnpackInt(kind)) {
		case compoundC:
			cached, ok := packer.GetCachedUnpack(cCache, val)
			if ok {
				result, ok = cached.(*CompoundC)
				if !ok {
					panic(fmt.Errorf("bad cached value: %v is not a CompoundC", cached))
				}
				return result
			}
			defer func() { packer.CacheUnpack(cCache, val, result) }()
			head, packedArgs := packer.UnpackPair(val)
			args := packer.UnpackList(packedArgs)
			resultArgs := make([]C, len(args))
			result = &CompoundC{
				IDTemplate(unpackTemplate(packer, head)),
				resultArgs,
			}
			for i, arg := range args {
				resultArgs[i] = unpackC(packer, arg)
			}
			return result
		case referenceC:
			return ReferenceC{packer.UnpackInt(val)}
		case constantC:
			return ConstC{unpackT(packer, val)}
		default:
			panic(errors.New("unknown kind of C"))
		}
	default:
		panic(errors.New("unknown version or wrong type of data"))
	}
}

func SaveC(c C) interface{} {
	return intern.FromRecord(packC(Recorder, c).(intern.Record))
}

func LoadC(b interface{}) (c C, err error) {
	defer func() { err = decodingRecover(recover()) }()
	return unpackC(Recorder, intern.MakeRecord(b)), nil
}

func (id CID) C() C {
	return ToC(id)
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

//Actions

type ActionCID intern.ID

func packActionC(packer intern.Packer, a ActionC) intern.Packed {
	intArgs := make([]intern.Packed, len(a.IntArgs))
	for i, arg := range a.IntArgs {
		intArgs[i] = packer.PackInt(arg)
	}
	packedIntArgs := packer.PackList(intArgs)
	args := make([]intern.Packed, len(a.Args))
	for i, arg := range a.Args {
		args[i] = packC(packer, arg)
	}
	packedArgs := packer.PackList(args)
	allArgs := packer.PackPair(packedIntArgs, packedArgs)
	act := packer.PackInt(int(a.Act))
	raw := packer.PackPair(act, allArgs)
	versioned := packer.PackPair(packer.PackInt(int(v1ActionC)), raw)
	return versioned
}

func unpackActionC(packer intern.Packer, x intern.Packed) ActionC {
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
			result.Args[i] = unpackC(packer, arg)
		}
		for i, arg := range intArgs {
			result.IntArgs[i] = packer.UnpackInt(arg)
		}
		return result
	default:
		panic(errors.New("unknown version or wrong type"))
	}
}

func IDActionC(a ActionC) ActionCID {
	return ActionCID(packActionC(IDer, a).(intern.ID))
}

func ToActionC(id ActionCID) ActionC {
	return unpackActionC(IDer, intern.ID(id))
}

func SaveActionC(a ActionC) interface{} {
	return intern.FromRecord(packActionC(Recorder, a).(intern.Record))
}

func LoadActionC(b interface{}) (act ActionC, err error) {
	defer func() { err = decodingRecover(recover()) }()
	return unpackActionC(Recorder, intern.MakeRecord(b)), nil
}

func (id ActionCID) ActionC() ActionC {
	return ToActionC(id)
}

func (action ActionC) ID() ActionCID {
	return IDActionC(action)
}

//Setting Line

func unpackSettingLine(packer intern.Packer, x intern.Packed) (SettingLine, int) {
	v, _ := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1Template:
		template := unpackTemplate(packer, x)
		return IDTemplate(template), template.Slots()
	case v1ActionC:
		return IDActionC(unpackActionC(packer, x)), 0
	default:
		panic(errors.New("unknown version or bad data type"))
	}
}

func packSettingLine(packer intern.Packer, l SettingLine) intern.Packed {
	switch l := l.(type) {
	case ActionCID:
		return packActionC(packer, ToActionC(l))
	case TemplateID:
		return packTemplate(packer, ToTemplate(l))
	default:
		panic(errors.New("unknown kind of setting line"))
	}
}

//SettingT

func unpackSettingT(packer intern.Packer, x intern.Packed) (result *SettingT) {
	if cached, ok := packer.GetCachedUnpack(settingTCache, x); ok {
		return cached.(*SettingT)
	}
	defer func() { packer.CacheUnpack(settingTCache, x, result) }()
	v, x := packer.UnpackPair(x)
	switch versionAndType(packer.UnpackInt(v)) {
	case v1SettingT:
		setting, args := packer.UnpackPair(x)
		result := &SettingT{}
		result.Setting = unpackSetting(packer, setting)
		for _, arg := range packer.UnpackList(args) {
			result.Args = append(result.Args, unpackT(packer, arg))
		}
		return result
	default:
		panic(errors.New("unknown version or bad data type"))
	}
}

func packSettingT(packer intern.Packer, s *SettingT) (result intern.Packed) {
	if cached, ok := packer.GetCachedPack(settingTCache, s); ok {
		return cached
	}
	defer func() { result = packer.CachePack(settingTCache, s, result) }()
	args := make([]intern.Packed, len(s.Args))
	for i, arg := range s.Args {
		args[i] = packT(packer, arg)
	}
	raw := packer.PackPair(packSetting(packer, s.Setting), packer.PackList(args))
	versioned := packer.PackPair(packer.PackInt(int(v1SettingT)), raw)
	return versioned
}

//Setting

type SettingID intern.ID

func unpackSetting(packer intern.Packer, x intern.Packed) (result *Setting) {
	/*
		if cached, ok := packer.GetCachedUnpack(settingCache, x); ok {
			return cached.(*Setting)
		}
		defer func() { packer.CacheUnpack(settingCache, x, result) }()
	*/
	packedLines := packer.UnpackList(x)
	result = Init()
	for _, packedLine := range packedLines {
		line, slots := unpackSettingLine(packer, packedLine)
		result = result.Append(line, slots)
	}
	return result
}

func packSetting(packer intern.Packer, s *Setting) (result intern.Packed) {
	/*
		if cached, ok := packer.GetCachedPack(settingCache, s); ok {
			return cached
		}
		defer func() {
			if result != nil {
				result = packer.CachePack(settingCache, s, result)
			}
		}()
	*/
	if s.Size == 0 {
		return packer.PackList([]intern.Packed{})
	}
	previous := packSetting(packer, s.Previous)
	last := packSettingLine(packer, s.Last)
	result = packer.AppendToPacked(previous, last)
	return result
}

func ToSetting(id SettingID) *Setting {
	return unpackSetting(IDer, intern.ID(id))
}

func IDSetting(s *Setting) SettingID {
	return SettingID(packSetting(IDer, s).(intern.ID))
}

func SaveSetting(s *Setting) interface{} {
	return intern.FromRecord(packSetting(Recorder, s).(intern.Record))
}

func LoadSetting(b interface{}) (s *Setting, err error) {
	defer func() { err = decodingRecover(recover()) }()
	return unpackSetting(Recorder, intern.MakeRecord(b)), nil
}

func (id SettingID) Setting() *Setting {
	return ToSetting(id)
}
