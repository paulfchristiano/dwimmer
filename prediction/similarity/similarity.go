package similarity

import (
	"container/heap"

	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
	"github.com/xrash/smetrics"
)

var (
	RelatedSettings = term.Make("what settings have been encountered before that are most analogous " +
		"to the setting [], and what is their relationship to that setting?")
)

func Suggestions(d dynamics.Dwimmer, s *term.Setting, n int) ([]term.ActionC, []float32) {
	settings, settingPriorities := analogies(d, s, n)
	result := []term.ActionC{}
	priorities := []float32{}
buildingResult:
	for i, other := range settings {
		t, _ := d.Get(other)
		simple, isSimple := t.(dynamics.SimpleTransition)
		if isSimple {
			action := simple.Action
			for _, otherAction := range result {
				if action.Id() == otherAction.Id() {
					continue buildingResult
				}
			}
			result = append(result, action)
			priorities = append(priorities, settingPriorities[i])
		}
	}
	return result, priorities
}

func match(a, b term.SettingLine) (float32, bool) {
	if a.Slots() != b.Slots() {
		return 0.0, false
	}
	switch a := a.(type) {
	case term.ActionCId:
		switch b := b.(type) {
		case term.ActionCId:
			return 1 - distance(a.String(), b.String()), true
		default:
			return 0.0, false
		}
	case term.TemplateId:
		switch b := b.(type) {
		case term.TemplateId:
			return 1 - distance(a.String(), b.String()), true
		default:
			return 0.0, false
		}
	default:
		return 0.0, false
	}
}

//TODO the algorithms part could be more efficient, but I don't really care
func analogies(d dynamics.Dwimmer, s *term.Setting, n int) ([]*term.Setting, []float32) {
	if s.Size == 1 {
		return contenders(d, s.Last, n)
	}

	previousSetting := s.Previous
	lastLine := s.Last
	previousAnalogies, previousPriorities := analogies(d, previousSetting, n+1)

	possibilities, possiblePriorities := []*term.Setting{}, new(indexHeap)
	i := 0
	for j, priority := range previousPriorities {
		analogy := previousAnalogies[j]
		for _, setting := range d.Continuations(analogy) {
			fit, canMatch := match(setting.Last, lastLine)
			if canMatch {
				possiblePriorities.Push(prioritized{
					index:    i,
					priority: priority * fit,
				})
				i++
				possibilities = append(possibilities, setting)
			}
		}
	}
	heap.Init(possiblePriorities)
	result := make([]*term.Setting, 0)
	priorities := make([]float32, 0)
	for i := 0; i < n && possiblePriorities.Len() > 0; i++ {
		next := heap.Pop(possiblePriorities).(prioritized)
		priorities = append(priorities, next.priority)
		result = append(result, possibilities[next.index])
	}
	return result, priorities
}

func contenders(d dynamics.Dwimmer, l term.SettingLine, n int) ([]*term.Setting, []float32) {
	allSettings := d.Continuations(term.Init())
	result, priorities := Top(l, allSettings, n)
	return result, priorities
}

func Top(target term.SettingLine, options []*term.Setting, n int) ([]*term.Setting, []float32) {
	h := makeHeap()
	ptms := make([]prioritized, 0)
	for i, option := range options {
		fit, canMatch := match(target, option.Last)
		if canMatch {
			ptms = append(ptms, prioritized{fit, i})
		}
	}
	h.heap = ptms
	heap.Init(h)
	result := make([]*term.Setting, 0)
	priorities := make([]float32, 0)
	for i := 0; i < n && h.Len() > 0; i++ {
		next := heap.Pop(h).(prioritized)
		priorities = append(priorities, next.priority)
		result = append(result, options[next.index])
	}
	return result, priorities
}

type indexHeap struct {
	heap []prioritized
}

func makeHeap() *indexHeap {
	return &indexHeap{make([]prioritized, 0)}
}

type prioritized struct {
	priority float32
	index    int
}

func (t *indexHeap) Len() int {
	return len(t.heap)
}

func (t *indexHeap) Less(i, j int) bool {
	//NOTE reversed order
	return t.heap[i].priority > t.heap[j].priority
}

func (t *indexHeap) Swap(i, j int) {
	t.heap[i], t.heap[j] = t.heap[j], t.heap[i]
}

func (t *indexHeap) Push(x interface{}) {
	t.heap = append(t.heap, x.(prioritized))
}

func (t *indexHeap) Pop() (result interface{}) {
	t.heap, result = t.heap[:t.Len()-1], t.heap[t.Len()-1]
	return result
}

//TODO this needs to be improved a lot!
func distance(a, b string) float32 {
	return float32(smetrics.WagnerFischer(a, b, 1, 1, 2)) / float32(len(a)+len(b))
}
