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

func init() {
	firsts = make(map[string]*set)
}

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
	// if s == nil {return false}
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
	// if self == nil {return false}
	for _, w := range *self {
		word := w.(tok)
		if word.ttype == e.ttype {return true}
	}
	return false
}

type production struct {
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
	
	for i, r := range prods {
		fmt.Println(i, ":", r)
	}
	
	fmt.Println(" ")
	
	for k, v := range firsts {
		fmt.Println(k, ":", v)
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

			// fmt.Println("Processing", prod, firsts[prod.name])

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
	word, err := s.nextWord()
	for err == nil && word.text != "%%" {
		if word.text == "%token" {parseTokenList()}
		word, err = s.nextWord()
	}
}

func parseTokenList() {
	for word, err := s.nextWord();
	    err == nil && word.ttype != newline;
		word, err = s.nextWord() {
		terms.Push(word)
	}
}

func parseProductions() {
	for word, err := s.nextWord();
		err == nil && word.text != "%%";
		word, err = s.nextWord() {
		if word.ttype == newline {continue}
		if word.ttype != nonterm {
			fmt.Println("Expected non-terminal but got", word.text)
			return
		}
		nonterms.Push(word)
		prod := new(production)
		prods.Push(prod)
		prod.name = word.text
		word, err = s.nextWord()
		if word.ttype != begindef {
			fmt.Println("Expected ':' but got", word.text)
			return
		}
		parseProduction(prod)
	}
}

func parseProduction(prod *production) {
	for word, err := s.nextWord();
		err == nil && word.ttype != enddef;
		word, err = s.nextWord() {
		HandleProductionToken:
		switch word.ttype {
		case newline:
			word, err = s.nextWord()
			if word.ttype != newline && word.ttype != alternate && word.ttype != enddef{
				fmt.Println("Expected an alternation or end but got", word.text)
				return
			}
			goto HandleProductionToken
		case alternate:
			prod = &production{prod.name, vector.Vector{}}
			prods.Push(prod)
		case enddef:
			return
		default:
			prod.seq.Push(word)
		}
	}
}
