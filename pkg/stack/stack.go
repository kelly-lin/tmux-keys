package stack

import "errors"

type Stack[T any] struct {
	items []T
}

func NewStack[T any](items []T) Stack[T] {
	return Stack[T]{items: items}
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var result T
		return result, errors.New("stack is empty")
	}

	result := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return result, nil
}

func (s *Stack[T]) Push(items ...T) {
	s.items = append(s.items, items...)
}

func (s *Stack[T]) HasItems() bool {
  return len(s.items) > 0
}
