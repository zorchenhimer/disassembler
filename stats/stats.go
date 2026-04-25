package stats

import (
)

type Data interface {
	StatName() string
}

type Group struct {
	Size int
	Items map[string]int
}

func NewGroup() *Group {
	return &Group{
		Items: make(map[string]int),
	}
}

type Set struct {
	Global *Group
	Groups map[string]*Group
}

func New() *Set {
	return &Set{
		Global: NewGroup(),
		Groups: make(map[string]*Group),
	}
}

func (s *Set) Incr(group, name string) {
	s.Global.Items[name]++

	if _, ok := s.Groups[group]; !ok {
		s.Groups[group] = NewGroup()
	}
	s.Groups[group].Items[name]++
}

func (s *Set) Add(group, name string) {
	if _, ok := s.Groups[group]; !ok {
		s.Groups[group] = NewGroup()
	}

	s.Global.Items[name] = 0
	s.Groups[group].Items[name] = 0
}

func (s *Set) SetSize(group string, size int) {
	if _, ok := s.Groups[group]; !ok {
		s.Groups[group] = NewGroup()
	}

	s.Global.Size += size
	s.Groups[group].Size = size
}
