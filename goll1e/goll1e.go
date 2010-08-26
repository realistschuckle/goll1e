package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"container/vector"
)

var e tok
var s scanner
var terms set
var nonterms set
var prods vector.Vector
var firsts map[string]*set
var follows map[string]*set

func init() {
	firsts = make(map[string]*set)
	follows = make(map[string]*set)
}

type production struct {
	num int
	name string	
	seq vector.Vector
}

func main() {
	filename := "stdin"
	in, out := os.Stdin, os.Stdout
	var err os.Error
	defer func() {
		if in != os.Stdin {in.Close()}
		if out != os.Stdout {out.Close()}
	}()

	if len(os.Args) > 1 {
		filename = os.Args[1]
		in, err = os.Open(filename, os.O_RDONLY, 0)
		if nil != err {
			fmt.Println("Cannot", err)
			os.Exit(-1)
		}
	}
	if len(os.Args) > 2 {
		flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		out, err = os.Open(os.Args[2], flags, 0666)
		if nil != err {
			fmt.Println("Cannot create output file", os.Args[2], ":", err)
			os.Exit(-1)
		}
	}
	
	content, err := ioutil.ReadAll(in)
	if nil != err {
		fmt.Println("Cannot", err)
		os.Exit(-2)
	}

	s.content = content
	parseGrammar()
	computeFirsts()
	computeFollows()
}

func computeFollows() {
	for _, w := range nonterms {
		word := w.(tok)
		follows[word.text] = &set{}
	}
	for true {
		changed := false
		for _, p := range prods {
			prod := p.(*production)
			for i, w := range prod.seq {
				word := w.(tok)
				if word.ttype != nonterm {continue}
				for j, x := range prod.seq[i + 1:] {
					xord := x.(tok)
					switch xord.ttype {
					case term:
						follows[word.text].Push(xord)
						goto NextProd
					case nonterm:
						changed = follows[word.text].Union(firsts[xord.text]) || changed
						if !firsts[xord.text].HasE() {goto NextProd}
						if j == len(prod.seq[i + 1:]) - 1 {
							changed = follows[word.text].Union(follows[prod.name]) || changed
						}
					}
				}
				NextProd:
			}
		}
		if !changed {break}
	}
}

func computeFirsts() {
	firsts[""] = &set{tok{"", empty}}
	for _, w := range nonterms {
		word := w.(tok)
		firsts[word.text] = &set{}
	}
	for _, w := range terms {
		word := w.(tok)
		s := &set{}
		s.Push(word)
		firsts[word.text] = s
	}
	for true {
		changed := false
		for _, p := range prods {
			prod := p.(*production)
			if len(prod.seq) == 0 {
				firsts[prod.name].Push(e)
				continue
			}
			first := prod.seq[0].(tok)
			switch first.ttype {
			case term:
				// fmt.Println("Found term", first)
				changed = firsts[prod.name].Push(first) || changed
			case nonterm:
				// fmt.Println("Found nonterm", first)
				if firsts[first.text].HasE() {                                     
					noe := firsts[first.text].NoE()
					changed = firsts[prod.name].Union(noe) || changed
					 for i := 1; i < len(prod.seq); i++ {
						t := prod.seq[i].(tok)
						if firsts[t.text].HasE() {
							noe = firsts[t.text].NoE()
							changed = firsts[prod.name].Union(noe) || changed
							if i == len(prod.seq) - 1 {
								firsts[prod.name].Push(e)
							}
						} else {
							changed = firsts[prod.name].Union(firsts[t.text]) || changed
							break
						}
					}
				} else {
					changed = firsts[prod.name].Union(firsts[first.text]) || changed
				}
			}
			// fmt.Println(" ")
		}
		if !changed {break}
	}
	for _, w := range terms {
		word := w.(tok)
		firsts[word.text] = nil, false
	}
}

func parseGrammar() {	
	parseHeader()
	parseProductions()
}

func parseHeader() {
	word, err := nextWord()
	for err == nil && word.text != "%%" {
		if word.text == "%token" {parseTokenList()}
		word, err = nextWord()
	}
}

func parseTokenList() {
	for word, err := nextWord();
	    err == nil && word.ttype != newline;
		word, err = nextWord() {
	}
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
		prod.num = len(prods)
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
		case newline:
			word, err = nextWord()
			if word.ttype != newline && word.ttype != alternate && word.ttype != enddef{
				fmt.Println("Expected an alternation or end but got", word.text)
				return
			}
			goto HandleProductionToken
		case alternate:
			prod = &production{len(prods), prod.name, vector.Vector{}}
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
	switch word.ttype {
	case term:
		terms.Push(word)
	case nonterm:
		nonterms.Push(word)
	}
	return
}
