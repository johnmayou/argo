package hashset

import "iter"

type HashSet[T comparable] struct {
	m map[T]struct{}
}

func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{m: make(map[T]struct{})}
}

func (s *HashSet[T]) Add(val T) {
	s.m[val] = struct{}{}
}

func (s *HashSet[T]) Remove(val T) {
	delete(s.m, val)
}

func (s *HashSet[T]) Contains(val T) bool {
	_, ok := s.m[val]
	return ok
}

func (s *HashSet[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.m {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *HashSet[T]) IsEmpty() bool {
	return len(s.m) == 0
}
