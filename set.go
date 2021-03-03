package main

type IntSet map[int]struct{}

func (s IntSet) Add(v int) {
	s[v] = struct{}{}
}

func (s IntSet) Contains(v int) bool {
	_, ok := s[v]
	return ok
}

func (s IntSet) Slice() (r []int) {
	for k := range s {
		r = append(r, k)
	}
	return
}

func (s *IntSet) Reset() {
	*s = map[int]struct{}{}
}
