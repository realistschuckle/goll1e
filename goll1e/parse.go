package main

import (
	"fmt"
	"os"
	"container/vector"
)

func parseHeader() {
	for word, err := nextWord();
		err == nil && word.text != "%%";
		word, err = nextWord() {
		switch word.text {
		case "%package":
			word, err = nextWord()
			packageName = word.text
		case "%import":
			for word, err := nextWord();
				err == nil && word.ttype != newline;
				word, err = nextWord() {
				imports.Push(word.text)
			}
		}
	}
	memorizeTerms = true
}

func parseProductions() {
	for word, err := nextWord();
		err == nil && word.text != "%%";
		word, err = nextWord() {
		if word.ttype == newline {continue}
		if word.ttype != nonterm {
			fmt.Println("Expected non-terminal but got", word.text)
			return
		}
		prod := new(production)
		prod.name = word.text
		prods.Push(prod)
		word, err = nextWord()
		if word.ttype != begindef {
			fmt.Println("Expected ':' but got", word.text)
			return
		}
		parseProduction(prod)
	}
}

func parseProduction(prod *production) {
	for word, err := nextWord();
		err == nil && word.ttype != enddef;
		word, err = nextWord() {
		HandleProductionToken:
		switch word.ttype {
		case code:
			prod.code = word.text
		case newline:
			word, err = nextWord()
			goto HandleProductionToken
		case alternate:
			prod = &production{prod.name, vector.Vector{}, ""}
			prods.Push(prod)
		case enddef:
			return
		default:
			prod.seq.Push(word)
		}
	}
}

func nextWord() (word tok, err os.Error) {
	word, err = s.nextWord()
	if word.text == "#" {
		for true {
			word, err = s.nextWord()
			if word.ttype == newline {break}
		}
	}
	if memorizeTerms {
		switch word.ttype {
		case term:
		if _, ok := terms[word.text]; !ok {
			terms[word.text] = len(terms)
		}
		case nonterm:
			if _, ok := nonterms[word.text]; !ok {
				nonterms[word.text] = len(nonterms)
			}
		}
	}
	return
}

func firstsFor(name string) (output *set) {
	output = &set{}
	for i, p := range prods {
		prod := p.(*production)
		if prod.name == name {output.Union(firsts[i])}
	}
	return
}


func followsFor(name string) *set {
	return follows[nonterms[name]]
}
