package dwimmer

import (
	"container/heap"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/parsing"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/paulfchristiano/dwimmer/ui"
	"github.com/xrash/smetrics"
)

var ActionQ = term.Make("what action should be taken in the setting []?")

func init() {
	dynamics.AddNativeResponse(ActionQ, 1, dynamics.Args1(findAction))
}

func findAction(d dynamics.Dwimmer, s *term.SettingT, quotedSetting term.T) term.T {
	setting, err := represent.ToSetting(d, quotedSetting)
	if err != nil {
		return term.Make("asked to decide what to do in setting [], "+
			"but while converting to a setting received []").T(quotedSetting, err)
	}
	settingId := term.IdSetting(setting)
	transition, ok := d.Get(settingId)
	if !ok {
		action := ElicitAction(d, settingId)
		transition = dynamics.SimpleTransition{action}
		d.Save(settingId, transition)
	}
	return core.Answer.T(represent.Transition(transition))
}

func ShowSettingS(d dynamics.Dwimmer, settingS *term.SettingS) {
	ui.Clear()
	for _, line := range settingS.Lines() {
		d.Writeln(line)
	}
}

//TODO only do this sometimes, or control better the length, or something...
func GetHints(d dynamics.Dwimmer, id term.SettingId, n int) []string {
	hs := hints(d, id, n)
	hint_strings := make([]string, len(hs))
	if len(hs) > 0 {
		for i, h := range hs {
			hint_strings[len(hs)-1-i] = "replace " + h.String()
			//d.Writeln(fmt.Sprintf("%d. %v", i, h))
		}
		//d.Writeln("")
	}
	return hint_strings
}

func ElicitAction(d dynamics.Dwimmer, id term.SettingId) term.ActionC {
	setting := id.Setting()
	settingS := addNames(setting)
	ShowSettingS(d, settingS)
	hint_strings := GetHints(d, id, 6)
	for {
		input := d.Readln(" < ", hint_strings)
		a := parsing.ParseAction(input, settingS.Names)
		if a == nil {
			c := parsing.ParseTerm(input, settingS.Names)
			if c != nil {
				switch c := c.(type) {
				case *term.CompoundC:
					if questionLike(c) {
						a = new(term.ActionC)
						*a = term.AskC(c)
					}
				case term.ReferenceC:
					a = new(term.ActionC)
					*a = term.ViewC(c)
				}
				d.Writeln("please input an action (ask, view, or return)")
			} else {
				d.Writeln("that response wasn't parsed correctly")
			}
		}
		if a != nil {
			for i, n := range a.IntArgs {
				if n == -1 {
					a.IntArgs[i] = len(setting.Outputs) - 1
				}
			}
			return *a
		}
	}
}

func questionLike(c term.C) bool {
	for _, char := range c.String() {
		if char == '?' {
			return true
		}
	}
	return false
}

type prioritized struct {
	priority float32
	index    int
}

type tmHeap struct {
	heap []prioritized
}

func makeHeap() *tmHeap {
	return &tmHeap{make([]prioritized, 0)}
}

func (t *tmHeap) Len() int {
	return len(t.heap)
}

func (t *tmHeap) Less(i, j int) bool {
	//NOTE reversed order
	return t.heap[i].priority > t.heap[j].priority
}

func (t *tmHeap) Swap(i, j int) {
	t.heap[i], t.heap[j] = t.heap[j], t.heap[i]
}

func (t *tmHeap) Push(x interface{}) {
	t.heap = append(t.heap, x.(prioritized))
}

func (t *tmHeap) Pop() (result interface{}) {
	t.heap, result = t.heap[:t.Len()-1], t.heap[t.Len()-1]
	return result
}

//TODO this needs to be improved a lot!
func distance(a, b string) float32 {
	return float32(smetrics.WagnerFischer(a, b, 1, 1, 2)) / float32(len(a)+len(b))
}

func hints(d dynamics.Dwimmer, settingId term.SettingId, num int) []*term.Template {
	prefixId := settingId.IdInit()
	s := settingId.IdLast().Template().String()
	h := makeHeap()
	alltms := make([]*term.Template, 0)
	ptms := make([]prioritized, 0)
	for _, tmid := range d.Continuations(prefixId) {
		tm := tmid.Template()
		ptms = append(ptms, prioritized{-distance(s, tm.String()), len(alltms)})
		alltms = append(alltms, tm)
	}
	h.heap = ptms
	heap.Init(h)
	result := make([]*term.Template, 0)
	for i := 0; i < num && h.Len() > 0; i++ {
		next := heap.Pop(h)
		result = append(result, alltms[next.(prioritized).index])
	}
	return result
}

var allNames = "xyzwijklmnstuvabcdefg"

func makeNames(n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = allNames[i : i+1]
	}
	return result
}

func addNames(s *term.Setting) *term.SettingS {
	names := makeNames(s.Slots())
	return &term.SettingS{term.IdSetting(s), names}
}
