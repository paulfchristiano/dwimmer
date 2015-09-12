package represent

import (
	"bytes"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/ints"
	"github.com/paulfchristiano/dwimmer/data/lists"
	"github.com/paulfchristiano/dwimmer/data/strings"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	UnquoteRune = term.Make("what character is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteActionC = term.Make("what action is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteC = term.Make("what term constructor is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteT = term.Make("what term is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteInt = term.Make("what integer is []? the answer should be represented " +
		"as a wrapped native integer")
	UnquoteStr = term.Make("what string is []? the answer should be represented " +
		"as a wrapped native string")
	UnquoteTemplate = term.Make("what template is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteList = term.Make("what list is []? the answer should be represented " +
		"with one of the heads referenced in the file quote.go")
	UnquoteSetting = term.Make("what setting is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
	UnquoteTransition = term.Make("what transition is []? the answer should be represented " +
		"with one of the heads defined in the file quote.go")
)

func ToTransition(d dynamics.Dwimmer, t term.T) (dynamics.Transition, term.T) {
	switch t.Head() {
	case QuotedSimpleTransition.Head():
		action, err := ToActionC(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert a simple transition, " +
				"but while converting the action received []").T(err)
		}
		return dynamics.SimpleTransition{action}, nil
	case QuotedNativeTransition.Head():
		wrapped := t.Values()[0].(term.Wrapper).Value
		return wrapped.Interface().(dynamics.NativeTransition), nil
	}
	reduced, err := d.Answer(UnquoteTransition.T(t))
	if err != nil {
		return nil, err
	}
	return ToTransition(d, reduced)
}

func ToActionC(d dynamics.Dwimmer, t term.T) (term.ActionC, term.T) {
	var nullaction term.ActionC
	switch t.Head() {
	case QuotedReturn.Head():
		arg, err := ToC(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled return action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		return term.ReturnC(arg), nil
	case QuotedAsk.Head():
		arg, err := ToC(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled ask action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		return term.AskC(arg), nil
	case QuotedView.Head():
		arg, err := ToC(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled view action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		return term.ViewC(arg), nil
	case QuotedReplace.Head():
		arg, err := ToC(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled replacement action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		n, err := ToInt(d, t.Values()[1])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled replacement action, "+
				"but received [] while converting its index[]").T(err, t.Values()[1])
		}
		return term.ReplaceC(arg, n), nil
	case QuotedClarify.Head():
		arg, err := ToC(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled clarify action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		n, err := ToInt(d, t.Values()[1])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled clarify action, "+
				"but received [] while converting its index[]").T(err, t.Values()[1])
		}
		return term.ClarifyC(arg, n), nil
	case QuotedCorrect.Head():
		n, err := ToInt(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled correct action, "+
				"but received [] while converting its index[]").T(err, t.Values()[0])
		}
		return term.CorrectC(n), nil
	case QuotedClose.Head():
		n, err := ToInt(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled close action, "+
				"but received [] while converting its index[]").T(err, t.Values()[0])
		}
		return term.CloseC(n), nil
	case QuotedDelete.Head():
		arg, err := ToInt(d, t.Values()[0])
		if err != nil {
			return nullaction, term.Make("asked to convert a compiled delete action, "+
				"but received [] while converting its argument []").T(err, t.Values()[0])
		}
		return term.DeleteC(arg), nil
	}
	reduced, err := d.Answer(UnquoteActionC.T(t))
	if err != nil {
		return nullaction, err
	}
	return ToActionC(d, reduced)
}

func ToC(d dynamics.Dwimmer, t term.T) (term.C, term.T) {
	switch t.Head() {
	case QuotedCompoundC.Head():
		quotedArgs, err := ToList(d, t.Values()[1])
		if err != nil {
			return nil, term.Make("asked to convert compound C, "+
				"but received [] when converting argument list []").T(err, t.Values()[1])
		}
		args := make([]term.C, len(quotedArgs))
		for i, quotedArg := range quotedArgs {
			args[i], err = ToC(d, quotedArg)
			if err != nil {
				return nil, term.Make("asked to convert compound C, "+
					"but received [] while converting one of its arguments, []").T(err, quotedArg)
			}
		}
		template, err := ToTemplate(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert compound C, "+
				"but received [] while converting its template []").T(err, t.Values()[0])
		}
		return template.C(args...), nil
	case QuotedReferenceC.Head():
		index, err := ToInt(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert reference C, "+
				"but received [] while converting its index []").T(err, t.Values()[0])
		}
		return term.ReferenceC{index}, nil
	case QuotedConstC.Head():
		val, err := ToT(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert constanc C, "+
				"but received [] while converting its value []").T(err, t.Values()[0])
		}
		return term.ConstC{val}, nil
	}
	reduced, err := d.Answer(UnquoteC.T(t))
	if err != nil {
		return nil, err
	}
	return ToC(d, reduced)
}

func ToInt(d dynamics.Dwimmer, t term.T) (int, term.T) {
	switch t := t.(type) {
	case term.Int:
		return int(t), nil
	case *term.CompoundT:
		switch t.Head() {
		case ints.Zero:
			return 0, nil
		case ints.One:
			return 1, nil
		case ints.Negative:
			k, err := ToInt(d, t.Values()[0])
			if err != nil {
				return 0, term.Make("asked to convert integer, but received [] while converting "+
					"subexpression []").T(err, t.Values()[0])
			}
			return -k, nil
		case ints.Double:
			k, err := ToInt(d, t.Values()[0])
			if err != nil {
				return 0, term.Make("asked to convert integer, but received [] while converting "+
					"subexpression []").T(err, t.Values()[0])
			}
			return 2 * k, nil
		case ints.DoublePlusOne:
			k, err := ToInt(d, t.Values()[0])
			if err != nil {
				return 0, term.Make("asked to convert integer, but received [] while converting "+
					"subexpression []").T(err, t.Values()[0])
			}
			return 2*k + 1, nil
		}
	}
	reduced, err := d.Answer(UnquoteInt.T(t))
	if err != nil {
		return 0, err
	}
	return ToInt(d, reduced)
}

func ToStr(d dynamics.Dwimmer, t term.T) (string, term.T) {
	switch t := t.(type) {
	case term.Str:
		return string(t), nil
	case *term.CompoundT:
		switch t.Head() {
		case strings.ByRunes:
			l, err := ToList(d, t.Values()[0])
			if err != nil {
				return "", term.Make("asked to convert string, "+
					"but received [] while converting its list of characters []").T(err, t.Values()[0])
			}
			var b bytes.Buffer
			for _, quotedRune := range l {
				r, err := ToRune(d, quotedRune)
				if err != nil {
					return "", term.Make("asked to convert string, but received [] while converting "+
						"character []").T(err, quotedRune)
				}
				b.WriteRune(r)
			}
			return b.String(), nil
		}
	}
	reduced, err := d.Answer(UnquoteStr.T(t))
	if err != nil {
		return "", err
	}
	return ToStr(d, reduced)
}

func ToRune(d dynamics.Dwimmer, t term.T) (rune, term.T) {
	switch t.Head() {
	case strings.Rune:
		unicode, err := ToInt(d, t.Values()[0])
		if err != nil {
			return 0, term.Make("asked to convert character, "+
				"but received [] while converting its unicode encoding []").T(err, t.Values()[0])
		}
		return rune(unicode), nil
	}
	reduced, err := d.Answer(UnquoteRune.T(t))
	if err != nil {
		return 0, reduced
	}
	return ToRune(d, reduced)
}

func ToTemplate(d dynamics.Dwimmer, t term.T) (term.TemplateId, term.T) {
	switch t.Head() {
	case QuotedTemplate.Head():
		quotedParts, err := ToList(d, t.Values()[0])
		if err != nil {
			return 0, term.Make("asked to convert template, "+
				"but received [] while converting list of arguments []").T(err, t.Values()[0])
		}
		parts := make([]string, len(quotedParts))
		for i, part := range quotedParts {
			parts[i], err = ToStr(d, part)
			if err != nil {
				return 0, term.Make("asked to convert template, "+
					"but received [] while converting one of its parts []").T(err, part)
			}
		}
		return term.IdTemplate(&term.Template{Parts: parts}), nil
	}
	reduced, err := d.Answer(UnquoteTemplate.T(t))
	if err != nil {
		return 0, reduced
	}
	return ToTemplate(d, reduced)
}

func ToT(d dynamics.Dwimmer, t term.T) (term.T, term.T) {
	switch t.Head() {
	case term.Quoted{}.Head():
		return t.(term.Quoted).Value, nil
	case QuotedCompoundT.Head():
		quotedArgs, err := ToList(d, t.Values()[1])
		if err != nil {
			return nil, term.Make("asked to convert compound term, "+
				"but received [] while converting its list of arguments []").T(err, t.Values()[1])
		}
		args := make([]term.T, len(quotedArgs))
		for i, quotedArg := range quotedArgs {
			args[i], err = ToT(d, quotedArg)
			if err != nil {
				return nil, term.Make("asked to convert compound term, "+
					"but received [] while converting one of its arguments, []").T(err, quotedArg)
			}
		}
		template, err := ToTemplate(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert compound term, "+
				"but received [] while converting its template []").T(err, t.Values()[0])
		}
		return template.T(args...), nil
	case QuotedIntT.Head():
		return t.Values()[0], nil
	case QuotedStrT.Head():
		return t.Values()[0], nil
	case QuotedWrapperT.Head():
		return t.Values()[0], nil
	case QuotedQuoteT.Head():
		return term.Quoted{t.Values()[0]}, nil
	}
	reduced, err := d.Answer(UnquoteT.T(t))
	if err != nil {
		return nil, err
	}
	return ToT(d, reduced)
}

func reverse(l []term.T) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func ToList(d dynamics.Dwimmer, t term.T) ([]term.T, term.T) {
	result, err := ToReverseList(d, t)
	if err != nil {
		return nil, term.Make("asked to convert list, " +
			"but received [] while converting the reversed list").T(err)
	}
	reverse(result)
	return result, nil
}

func ToReverseList(d dynamics.Dwimmer, t term.T) ([]term.T, term.T) {
	switch t.Head() {
	case lists.Empty.Head():
		return []term.T{}, nil
	case lists.Cons.Head():
		reversetail, err := ToReverseList(d, t.Values()[1])
		if err != nil {
			return nil, term.Make("asked to convert list, "+
				"but received [] while converting its tail []").T(err, t.Values()[1])
		}
		return append(reversetail, t.Values()[0]), nil
	}
	reduced, err := d.Answer(UnquoteList.T(t))
	if err != nil {
		return nil, err
	}
	return ToReverseList(d, reduced)
}

func toSettingLines(d dynamics.Dwimmer, t term.T) ([]term.TemplateId, []term.ActionCId, term.T) {
	outputs := make([]term.TemplateId, 0)
	inputs := make([]term.ActionCId, 0)
	var err term.T
	l, err := ToList(d, t)
	if err != nil {
		return nil, nil, term.Make("asked to convert a list of lines in a setting, " +
			"but received [] while converting to a list at all").T(err)
	}
	for i := 0; 2*i < len(l); i++ {
		output, err := ToTemplate(d, l[2*i])
		if err != nil {
			return nil, nil, term.Make("asked to convert a list of lines in a setting, "+
				"but received [] while converting even numbered argument [] to a template").T(err, l[2*i])
		}
		outputs = append(outputs, output)
		if 2*i+1 < len(l) {
			input, err := ToActionC(d, l[2*i+1])
			if err != nil {
				return nil, nil, term.Make("asked to convert a list of lines in a setting, "+
					"but received [] while converting odd numbered argument [] to an action").T(err, l[2*i+1])
			}
			inputs = append(inputs, term.IdActionC(input))
		}
	}
	return outputs, inputs, nil
}
func ToSetting(d dynamics.Dwimmer, t term.T) (*term.Setting, term.T) {
	switch t.Head() {
	case QuotedSetting.Head():
		outputs, inputs, err := toSettingLines(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert setting, "+
				"but received [] while trying to convert its lines []").T(err, t.Values()[0])
		}
		return &term.Setting{outputs, inputs}, nil
	}
	reduced, err := d.Answer(UnquoteSetting.T(t))
	if err != nil {
		return nil, err
	}
	return ToSetting(d, reduced)
}

func ToSettingT(d dynamics.Dwimmer, t term.T) (*term.SettingT, term.T) {
	switch t.Head() {
	case QuotedSettingT.Head():
		outputs, inputs, err := toSettingLines(d, t.Values()[0])
		if err != nil {
			return nil, term.Make("asked to convert concrete setting, "+
				"but received [] while trying to convert its lines []").T(err, t.Values()[0])
		}
		id := term.IdSetting(&term.Setting{outputs, inputs})
		args, err := ToList(d, t.Values()[1])
		if err != nil {
			return nil, term.Make("asked to convert concrete setting, "+
				"but received [] while trying to convert its arguments []").T(err, t.Values()[1])
		}
		children, err := ToList(d, t.Values()[2])
		result := &term.SettingT{
			SettingId: id,
			Args:      make([]term.T, len(args)),
			Children:  make([]*term.SettingT, len(children)),
		}
		for i, arg := range args {
			result.Args[i], err = ToT(d, arg)
			if err != nil {
				return nil, term.Make("asked to convert concrete setting, "+
					"but received [] while trying to convert its argument []").T(err, arg)
			}
		}
		for i, child := range children {
			result.Children[i], err = ToSettingT(d, child)
			if err != nil {
				return nil, term.Make("asked to convert concrete setting, "+
					"but received [] while trying to convert its subordinate setwting []").T(err, child)
			}
		}
		return result, nil
	case QuotedNil.Head():
		return nil, nil
	}
	reduced, err := d.Answer(UnquoteSetting.T(t))
	if err != nil {
		return nil, err
	}
	return ToSettingT(d, reduced)
}

func init() {
	s := term.InitS().AppendTemplate(UnquoteInt, "x")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("x")))
	dynamics.AddSimple(s.Copy().AppendTemplate(ints.Zero), term.ReturnS(core.Answer.S(term.Sc(term.Int(0)))))
	dynamics.AddSimple(s.Copy().AppendTemplate(ints.One), term.ReturnS(core.Answer.S(term.Sc(term.Int(1)))))
}
