package main

import (
	"container/vector"
)

type set vector.Vector

func (self *set) Push(t tok) bool {
	for _, s := range *self {
		st := s.(tok)
		if st.text == t.text {return false}
	}
	(*vector.Vector)(self).Push(t)
	return true
}

func (self *set) Union(s *set) bool {
	changed := false
	for _, w := range *s {
		word := w.(tok)
		changed = changed || self.Push(word)
	}
	return changed
}                   

func (self *set) NoE() (output *set) {
	output = new(set)
	for _, w := range *self {
		word := w.(tok)
		if word.text != e.text {
			output.Push(word)
		}
	}
	return
}

func (self *set) HasE() bool {
	for _, w := range *self {
		word := w.(tok)
		if word.ttype == e.ttype {return true}
	}
	return false
}

