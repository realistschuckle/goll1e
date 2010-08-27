package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"container/vector"
)

var s scanner
var terms map[string]int
var nonterms map[string]int
var prods vector.Vector
var firsts map[int]*set
var follows map[int]*set

func init() {
	firsts = make(map[int]*set)
	follows = make(map[int]*set)
	terms = make(map[string]int)
	nonterms = make(map[string]int)
	terms[""] = 0
}

type production struct {
	name string	
	seq vector.Vector
}

func (self *production) String() (output string) {
	output = self.name + "["
	for i, w := range self.seq {
		word := w.(tok)
		if i > 0 {output += " "}
		output += word.text
	}
	output += "]"
	return
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
	
	printTermMap(terms)
	printTermMap(nonterms)
	printProductions()
	printFirsts()
}

func printFirsts() {
	for i := 0; i < len(firsts); i++ {
		prod := prods[i].(*production)
		fmt.Print(prod.name, "[")
		for j, t := range *firsts[i] {
			if j > 0 {fmt.Print(",")}
			token := t.(int)
			fmt.Print(translateTerm(token))
		}
		fmt.Print("]\n")
	}
	fmt.Println(" ")
	for k, _ := range nonterms {
		firstFor := firstsFor(k)
		fmt.Print(k, "[")
		for j, t := range *firstFor {
			if j > 0 {fmt.Print(",")}
			token := t.(int)
			fmt.Print(translateTerm(token))
		}
		fmt.Print("]\n")
	}
}

func translateTerm(i int) (output string) {
	output = "<NOT FOUND>"
	for k, v := range terms {
		if v != i {continue}
		output = k
		break
	}
	return
}

func translateTerms(s *set) (output string) {
	output = "["
	for i, v := range *s {
		token := v.(int)
		if i > 0 {output += ","}
		output += translateTerm(token)
	}
	output += "]"
	return
}

func printProductions() {
	for i, p := range prods {
		fmt.Println(i, p)
	}
	fmt.Println(" ")
}

func printTermMap(terms map[string]int) {
	sort := make([]string, len(terms))
	for k, v := range terms {
		sort[v] = k
	}
	for i, v := range sort {
		fmt.Println(i, v)
	}
	fmt.Println(" ")
}

func computeFollows() {
}

func computeFirsts() {
	for i := 0; i < len(prods); i++ {
		firsts[i] = &set{}
	}
	for true {
		changed := false
		for prodidx, p := range prods {
			prod := p.(*production)
			if len(prod.seq) == 0 {
				changed = firsts[prodidx].Push(terms[""]) || changed
				// fmt.Println("Adding <EMPTY> to", prod.name, "and changed?", changed)
				goto NextProd
			}
			for _, w := range prod.seq {
				word := w.(tok)
				switch word.ttype {
				case term:
					changed = firsts[prodidx].Push(terms[word.text]) || changed
					// fmt.Println("Adding", word.text, "to", prod.name, "and changed?", changed)
					goto NextProd
				case nonterm:
					fs := firstsFor(word.text)
					changed = firsts[prodidx].Union(fs.NoE()) || changed
					// fmt.Println("Adding firsts of", word.text, translateTerms(fs), "to", prod.name, "and changed?", changed)
					if !fs.HasE() {goto NextProd}
				}
			}
			NextProd:
		}
		// fmt.Println(" ")
		if !changed {break}
	}
}

func parseGrammar() {	
	parseProductions()
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
		case newline:
			word, err = nextWord()
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

func nextWord() (word tok, err os.Error) {
	word, err = s.nextWord()
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
