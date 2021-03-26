// Package stack implements the logic of data structure stack, element type is unlimited.

package stack

import (
	"fmt"
)

// Stack the data structure stack.
type Stack struct {
	data []interface{}
	p    int
}

// New returns a new Stack instance.
func New() *Stack {
	return &Stack{
		data: []interface{}{},
		p:    -1,
	}
}

// Pop return and remove the latest added element.
func (s *Stack) Pop() (interface{}, bool) {
	r, ok := s.Top()
	if ok {
		s.p--
	}

	return r, ok
}

// Push add element to stack.
func (s *Stack) Push(i interface{}) {
	s.p++
	if s.p < len(s.data) {
		s.data[s.p] = i
		return
	}

	s.data = append(s.data, i)
}

// Top returns the latest added element.
func (s *Stack) Top() (interface{}, bool) {
	if s.p < 0 {
		return nil, false
	}

	return s.data[s.p], true
}

// Clear flushes stack.
func (s *Stack) Clear() {
	s.p = -1
}

// String returns the data in string format.
func (s *Stack) String() string {
	return fmt.Sprintf("%v", s.data[:s.p+1])
}
