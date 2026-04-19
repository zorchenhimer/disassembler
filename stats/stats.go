package stats

import (
)

type Data interface {
	StatName() string
}

type Group map[string]int

type Set struct {
	Global Group
	Groups map[string]Group
}

func New() *Set {
	return &Set{
		Global: Group{},
		Groups: make(map[string]Group),
	}
}

func (s *Set) Incr(group, name string) {
	s.Global[name]++

	if _, ok := s.Groups[group]; !ok {
		s.Groups[group] = Group{}
	}
	s.Groups[group][name]++
}

func (s *Set) Add(group, name string) {
	if _, ok := s.Groups[group]; !ok {
		s.Groups[group] = Group{}
	}

	s.Global[name] = 0
	s.Groups[group][name] = 0
}
