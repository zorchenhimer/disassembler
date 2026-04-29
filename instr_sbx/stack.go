package instrsbx

import (
	"fmt"
)

type Stack[T any] struct {
	data []T
	bottom int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{data: []T{}, bottom: 0}
}

func (s *Stack[T]) Array() []T {
	return s.data[0:s.bottom]
}

func (s *Stack[T]) Push(v T) error {
	if s.bottom > len(s.data) {
		return fmt.Errorf("stack bottom larger than stack")
	}

	if len(s.data) == s.bottom {
		s.data = append(s.data, v)
		s.bottom++
		return nil
	}

	s.data[s.bottom] = v
	s.bottom++
	return nil
}

func (s *Stack[T]) Pop() (T, error) {
	var z T
	if s.bottom <= 0 {
		return z, fmt.Errorf("empty on stack.Pop()")
	}

	s.bottom--
	return s.data[s.bottom], nil
}

func (s *Stack[T]) Get(i int64) (T, error) {
	var z T
	if i < 0 {
		return z, fmt.Errorf("negative value in stack.Get()")
	}

	idx := int64(s.bottom) - i - 1
	if idx < 0 {
		return z, fmt.Errorf("too deep in stack.Get()")
	}

	return s.data[idx], nil
}
