package stack

import (
	"fmt"
)

type Stack struct {
	data []interface{}
	p    int
}

func New() *Stack {
	return &Stack{
		data: []interface{}{},
		p:    -1,
	}
}

func (s *Stack) Pop() (interface{}, bool) {
	r, ok := s.Top()
	if ok {
		s.p--
	}

	return r, ok
}

func (s *Stack) Push(i interface{}) {
	s.p++
	if s.p < len(s.data) {
		s.data[s.p] = i
		return
	}

	s.data = append(s.data, i)
}

func (s *Stack) Top() (interface{}, bool) {
	if s.p < 0 {
		return nil, false
	}

	return s.data[s.p], true
}

func (s *Stack) String() string {
	return fmt.Sprintf("%v", s.data[:s.p+1])
}
