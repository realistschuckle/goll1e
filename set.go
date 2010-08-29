package main

import (
	"container/vector"
)

type set vector.Vector

func (self *set) Push(i int) bool {
	if self.IndexOf(i) != -1 {return false}
	(*vector.Vector)(self).Push(i)
	return true
}

func (self *set) IndexOf(t int) int {
	for i, w := range *self {
		word := w.(int)
		if word == t {return i}
	}
	return -1
}

func (self *set) Union(s *set) bool {
	changed := false
	for _, w := range *s {
		word := w.(int)
		changed = self.Push(word) || changed
	}
	return changed
}                   

func (self *set) NoE() (output *set) {
	output = new(set)
	for _, w := range *self {
		word := w.(int)
		if word != 0 {
			output.Push(word)
		}
	}
	return
}

func (self *set) HasE() bool {
	for _, w := range *self {
		word := w.(int)
		if word == 0 {return true}
	}
	return false
}

