package stack

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func TestStack(t *testing.T) {
	s := New()
	n := 10
	arr := make([]int, n+1)

	for i := 0; i <= n; i++ {
		arr[i] = i
		s.Push(i)
		r, ok := s.Top()
		requireEquals(t, true, ok)
		requireEquals(t, i, r)
	}
	requireEquals(t, fmt.Sprintf("%v", arr), s.String())

	for i := 0; i <= n; i++ {
		r, ok := s.Top()
		requireEquals(t, true, ok)
		requireEquals(t, n, r)
	}

	for i := n; i >= 0; i-- {
		r, ok := s.Pop()
		requireEquals(t, true, ok)
		requireEquals(t, i, r)
	}
	requireEquals(t, "[]", s.String())

	for i := n; i >= 0; i-- {
		r, ok := s.Top()
		requireEquals(t, false, ok)
		requireEquals(t, nil, r)
	}
}

func requireEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		debug.PrintStack()
		t.Fatalf("not equals, expected: %v, actual: %v\n", expected, actual)
	}
}
