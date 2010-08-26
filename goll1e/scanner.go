package main

import (
	"os"
	"unicode"
	"utf8"
)

type tokType int

const (
	empty tokType = iota
	nonterm
	term
	newline
	enddef
	alternate
	begindef
	pcent
	other
)

type tok struct {
	text string
	ttype tokType
}

type scanner struct {
	index int
	content []uint8
}

func (self *scanner) nextWord() (word tok, err os.Error) {
	if self.index >= len(self.content) {
		err = os.NewError("EOF")
		return
	}

	for ; self.index < len(self.content); {
		r, l := utf8.DecodeRune(self.content[self.index:])
		if !unicode.IsSpace(r) || r == '\n' {break}
		self.index += l;
	}
	j, ttype, inchar := self.index, other, false
	for ; self.index < len(self.content); {
		r, l := utf8.DecodeRune(self.content[self.index:])
		if r == '\'' {inchar = !inchar}
		if self.index == j {
			switch {
			case unicode.IsUpper(r):
				ttype = nonterm
			case r == '\n':
				self.index++
				ttype = newline
				break
			case r == ':':
				ttype = begindef
			case r == ';':
				ttype = enddef
			case r == '|':
				ttype = alternate
			default:
				ttype = term
			}
		}
		if !inchar && unicode.IsSpace(r) {break}
		self.index += l
	}
	token := string(self.content[j:self.index])
	if ttype == newline {token = ""}
	word = tok{token, ttype}
	return
}
