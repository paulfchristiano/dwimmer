package term

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/paulfchristiano/dwimmer/term/intern"
)

type Template struct {
	Parts []string
}

func (t *Template) Slots() int {
	return len(t.Parts) - 1
}

func (t *Template) ShowWith(ss ...string) string {
	names := make([]string, len(ss))
	for i, s := range ss {
		names[i] = fmt.Sprintf("[%s]", s)
	}
	return interleave(t.Parts, names)
}

var sep = "[]"

func (t TemplateId) String() string {
	return t.Template().String()
}

func (t Template) String() string {
	return strings.Join(t.Parts, sep)
}

func (t TemplateId) Head() TemplateId {
	return t
}

func (t TemplateId) Parts() []string {
	return t.Template().Parts
}

type CompoundT struct {
	TemplateId
	args []T
}

type CompoundC struct {
	TemplateId
	args []C
}

type CompoundS struct {
	TemplateId
	args []S
}

type T interface {
	Head() TemplateId
	Values() []T
	Pack(intern.Packed) intern.Packed
	String() string
}

type C interface {
	String() string
	Pack(intern.Packed) intern.Packed
	Values() []C
	Instantiate([]T) T
	Uninstantiate([]string) S
}

type S interface {
	String() string
	Values() []S
	Instantiate([]string) C
}

func interleave(as, bs []string) string {
	buffer := new(bytes.Buffer)
	for i, a := range as {
		buffer.Write([]byte(a))
		if i < len(bs) {
			buffer.Write([]byte(bs[i]))
		}
	}
	return buffer.String()
}

func (t *CompoundT) String() string {
	args := make([]string, len(t.args))
	for i, arg := range t.args {
		args[i] = fmt.Sprintf("[%s]", arg.String())
	}
	return interleave(t.Parts(), args)
}

func (t *CompoundT) Values() []T {
	return t.args
}

func (t *CompoundC) Instantiate(ts []T) T {
	args := make([]T, len(t.args))
	for i, arg := range t.args {
		args[i] = arg.Instantiate(ts)
	}
	return &CompoundT{t.TemplateId, args}
}

func (t *CompoundC) Uninstantiate(names []string) S {
	args := make([]S, len(t.args))
	for i, arg := range t.args {
		args[i] = arg.Uninstantiate(names)
	}
	return t.S(args...)
}

func (t *CompoundS) Instantiate(names []string) C {
	args := make([]C, len(t.args))
	for i, arg := range t.args {
		args[i] = arg.Instantiate(names)
	}
	return &CompoundC{t.TemplateId, args}
}

func (t *CompoundC) String() string {
	args := make([]string, len(t.args))
	for i, arg := range t.args {
		args[i] = fmt.Sprintf("[%s]", arg.String())
	}
	return interleave(t.Parts(), args)
}
func (t *CompoundS) String() string {
	args := make([]string, len(t.args))
	for i, arg := range t.args {
		args[i] = fmt.Sprintf("[%s]", arg.String())
	}
	return interleave(t.Parts(), args)
}
func (t *CompoundC) Values() []C {
	return t.args
}
func (t *CompoundS) Values() []S {
	return t.args
}

func Make(ss ...string) TemplateId {
	parts := make([]string, 0)
	for _, s := range ss {
		parts = append(parts, strings.Split(s, sep)...)
	}
	return IdTemplate(&Template{parts})
}

func (t TemplateId) T(ts ...T) T {
	if t.Slots() != len(ts) {
		panic(fmt.Sprintf("instantiating %v with arguments %v", t.String(), ts))
	}
	return &CompoundT{t, ts}
}
func (t TemplateId) C(cs ...C) C {
	return &CompoundC{t, cs}
}
func (t TemplateId) S(ss ...S) S {
	if len(t.Template().Parts) != len(ss)+1 {
		panic(fmt.Sprintf("instatiating %v with arguments %v", t, ss))
	}
	return &CompoundS{t, ss}
}

type ReferenceC struct {
	Index int
}

func (r ReferenceC) Instantiate(ts []T) T {
	return ts[r.Index]
}

func (r ReferenceC) Uninstantiate(names []string) S {
	return ReferenceS{names[r.Index]}
}

func Sr(s string) S {
	return ReferenceS{s}
}

func Cr(n int) C {
	return ReferenceC{n}
}

func Sc(t T) S {
	return ConstS{t}
}

func Cc(t T) C {
	return ConstC{t}
}

func (r ReferenceC) Values() []C {
	return make([]C, 0)
}

func (r ReferenceC) String() string {
	return fmt.Sprintf("#%d", r.Index)
}

type ReferenceS struct {
	name string
}

func (r ReferenceS) Instantiate(names []string) C {
	for i, name := range names {
		if name == r.name {
			return ReferenceC{i}
		}
	}
	panic(fmt.Sprintf("tried to instantiate a ReferenceS with name %v in the enviornment %v", r.name, names))
}

func (r ReferenceS) String() string {
	return fmt.Sprintf("#%s", r.name)
}

func (r ReferenceS) Values() []S {
	return make([]S, 0)
}
