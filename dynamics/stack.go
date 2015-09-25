package dynamics

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/term"
)

type Stack interface {
	Push(t term.T)
	Pop() term.T
	ShowStack()
}

type BasicStack struct {
	stack []term.T
}

func (s *BasicStack) Push(t term.T) {
	s.stack = append(s.stack, t)
}

func (s *BasicStack) Pop() term.T {
	result := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return result
}

const (
	fromTop    = 10
	fromBottom = 10
)

func (s *BasicStack) ShowStack() {
	fmt.Println("Dwimmer stack:")
	n := len(s.stack)
	top, bottom := fromTop, fromBottom
	if n <= fromTop+fromBottom {
		top = n
		bottom = 0
	}
	for i := 0; i < top; i++ {
		fmt.Println(s.stack[n-1-i].Head())
	}
	if bottom > 0 {
		fmt.Printf("... [%d entries elided] ...\n", n-top-bottom)
		for i := 0; i < bottom; i++ {
			fmt.Println(s.stack[bottom-1-i].Head())
		}
	}
}
