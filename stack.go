package main

import "errors"

type IntStack []int

func (s *IntStack) Add(v int) {
	*s = append(*s, v)
}

func (s *IntStack) Pop() int {
	ss := *s

	if len(ss) == 0 {
		panic(errors.New("stack: can't pop from empty stack"))
	}

	ret := ss[len(ss)-1]
	*s = ss[:len(ss)-1]
	return ret
}

func (s IntStack) Empty() bool {
	return len(s) == 0
}
