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
	begincode
	code
	endcode
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

func (self *scanner) remainder() []byte {
	return self.content[self.index:]
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
	j, ttype, inchar, incode, outcode := self.index, other, false, false, false
	outcode = incode
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
			case r == '{':
				incode = true
				ttype = code
			default:
				ttype = term
			}
		} else if incode && r == '}' {
			incode = false
			outcode = true
		}
		if !incode && !inchar && unicode.IsSpace(r) {break}
		self.index += l
		if outcode {break}
	}
	token := string(self.content[j:self.index])
	if ttype == newline {token = ""}
	word = tok{token, ttype}
	return
}
