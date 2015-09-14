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
	for i, other := range settings {
		t, _ := d.Get(other.Id())
		simple, isSimple := t.(dynamics.SimpleTransition)
		if isSimple {
			action := simple.Action
			result = append(result, action)
			priorities = append(priorities, settingPriorities[i])
		}
	}
	return result, priorities
}

func match(a, b *term.Template) (float32, bool) {
	if a.Slots() != b.Slots() {
		return 0.0, false
	}
	return 1 - distance(a.String(), b.String()), true
}

func analogies(d dynamics.Dwimmer, s *term.Setting, n int) ([]*term.Setting, []float32) {
	if len(s.Inputs) == 0 {
		return contenders(d, s.Outputs[0], n)
	}

	previousSetting := &term.Setting{Inputs: s.Inputs[:len(s.Inputs)-1], Outputs: s.Outputs[:len(s.Outputs)-1]}
	lastAction := s.Inputs[len(s.Inputs)-1].ActionC()
	lastTemplate := s.Outputs[len(s.Outputs)-1].Template()
	previousAnalogies, previousPriorities := analogies(d, previousSetting, n+1)

	possibilities, possiblePriorities := []*term.Setting{}, new(indexHeap)
	i := 0
	for j, priority := range previousPriorities {
		analogy := previousAnalogies[j]
		t, _ := d.Get(analogy.Id())
		simple, isSimple := t.(dynamics.SimpleTransition)
		if isSimple {
			action := simple.Action
			priority = priority * (1 - distance(action.String(), lastAction.String()))
			analogy.AppendAction(action.Id())
			for _, template := range d.Continuations(analogy.Id()) {
				fit, canMatch := match(template.Template(), lastTemplate)
				if canMatch {
					possiblePriorities.Push(prioritized{
						index:    i,
						priority: priority * fit,
					})
					i++
					possibilities = append(possibilities, analogy.Copy().AppendTemplate(template))
				}
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

func contenders(d dynamics.Dwimmer, t term.TemplateId, n int) ([]*term.Setting, []float32) {
	allQs := d.Continuations(term.InitId)
	templates, priorities := Top(t, allQs, n)
	result := make([]*term.Setting, len(templates))
	for i, template := range templates {
		result[i] = &term.Setting{Inputs: []term.ActionCId{}, Outputs: []term.TemplateId{template}}
	}
	return result, priorities
}

func Top(target term.TemplateId, options []term.TemplateId, n int) ([]term.TemplateId, []float32) {
	h := makeHeap()
	targetTemplate := target.Template()
	ptms := make([]prioritized, 0)
	for i, tmid := range options {
		fit, canMatch := match(targetTemplate, tmid.Template())
		if canMatch {
			ptms = append(ptms, prioritized{fit, i})
		}
	}
	h.heap = ptms
	heap.Init(h)
	result := make([]term.TemplateId, 0)
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
