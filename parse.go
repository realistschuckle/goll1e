package main

import (
	"fmt"
	"os"
	"container/vector"
)

func parseHeader() {
	for word, err := nextWord(); err == nil && word.text != "%%"; word, err = nextWord() {
		if word.ttype == newline {
			continue
		}
		switch {
		case word.text == "%dev":
			dev = true
		case word.text == "%package":
			word, err = nextWord()
			packageName = word.text
		case word.text == "%union":
			word, err = nextWord() // {
			word, err = nextWord() // newline
			parseUnionEntries()
		case word.text == "%import":
			for word, err := nextWord(); err == nil && word.ttype != newline; word, err = nextWord() {
				imports.Push(word.text)
			}
		case len(word.text) > 6 && word.text[0:5] == "%type":
			parseTypedEntries(word.text[6 : len(word.text)-1])
		case len(word.text) > 7 && word.text[0:6] == "%token":
			parseTypedEntries(word.text[7 : len(word.text)-1])
		case word.text == "%defaultcode":
			memorizeTerms = true
			code, _ := nextWord()
			defaultcode = code.text
			memorizeTerms = false
		default:
			fmt.Println("Unrecognized header entry:", word.text)
		}
	}
	memorizeTerms = true
}

func parseTypedEntries(etype string) {
	for name, err := nextWord(); err == nil && name.ttype != newline; name, err = nextWord() {
		typedEntries[name.text] = etype
	}
}

func parseUnionEntries() {
	for name, err := nextWord(); err == nil && name.text != "}"; name, err = nextWord() {
		etype := ""
		for typespec, err := nextWord(); err == nil && typespec.ttype != newline; typespec, err = nextWord() {
			etype += typespec.text
		}
		unionEntries[name.text] = etype
	}
}

func parseProductions() {
	for word, err := nextWord(); err == nil && word.text != "%%"; word, err = nextWord() {
		if word.ttype == newline {
			continue
		}
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
	for _, p := range prods {
		prod := p.(*production)
		if len(prod.code) == 0 {
			prod.code = defaultcode
		}
	}
	memorizeTerms = false
	terms["%%"] = 0, false
}

func parseProduction(prod *production) {
	for word, err := nextWord(); err == nil && word.ttype != enddef; word, err = nextWord() {
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
			if word.ttype == newline {
				break
			}
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
		if prod.name == name {
			output.Union(firsts[i])
		}
	}
	return
}


func followsFor(name string) *set {
	return follows[nonterms[name]]
}
