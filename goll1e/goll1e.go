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
var table [][]int

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

	printTermMap(terms)
	printTermMap(nonterms)
	printProductions()

	computeFirsts()
	computeFollows()
	
	printSet(firsts, func(i int) string {return prods[i].(*production).name})
	printSet(follows, func(i int) string {return translateNonterm(i)})
	
	computeTable()
	
	printTable()
	printTableRaw()
}

func computeTable() {
	table = make([][]int, len(nonterms))
	for i := 0; i < len(nonterms); i++ {
		table[i] = make([]int, len(terms))
		for j := 0; j < len(terms); j++ {
			table[i][j] = -1
		}
	}
	for nonterm, r := range nonterms {
		for _, c := range terms {
			for prodidx, p := range prods {
				prod := p.(*production)
				if prod.name != nonterm {continue}
				switch {
				case firsts[prodidx].IndexOf(c) > -1:
					table[r][c] = prodidx
				case firsts[prodidx].HasE() && follows[r].IndexOf(c) > -1:
					table[r][c] = prodidx
				}
			}
		}
	}
}

func computeFollows() {
	for i := 0; i < len(nonterms); i++ {
		follows[i] = &set{}
	}
	follows[0].Push(-1)
	for true {
		changed := false
		for _, p := range prods {
			prod := p.(*production)
			for sidx, s := range prod.seq {
				word := s.(tok)
				if word.ttype != nonterm {continue}
				wordidx := nonterms[word.text]
				after := prod.seq[sidx + 1:]
				if len(after) == 0 {
					fs := followsFor(prod.name)
					changed = follows[wordidx].Union(fs.NoE()) || changed
					goto NextItem
				}
				for seqidx, t := range after {
					next := t.(tok)
					last := seqidx == len(after) - 1
					switch next.ttype {
					case term:
						changed = follows[wordidx].Push(terms[next.text]) || changed
						goto NextItem
					case nonterm:
						fs := firstsFor(next.text)
						changed = follows[wordidx].Union(fs.NoE()) || changed
						if !fs.HasE() {goto NextItem}
						if last {
							fs = followsFor(prod.name)
							changed = follows[wordidx].Union(fs.NoE()) || changed
						}
					}
				}
				NextItem:
			}
		}
		if !changed {break}
	}
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
				goto NextProd
			}
			for wordidx, w := range prod.seq {
				word := w.(tok)
				switch word.ttype {
				case term:
					changed = firsts[prodidx].Push(terms[word.text]) || changed
					goto NextProd
				case nonterm:
					fs := firstsFor(word.text)
					changed = firsts[prodidx].Union(fs.NoE()) || changed
					if !fs.HasE() {goto NextProd}
					if wordidx == len(prod.seq) - 1 {
						changed = firsts[prodidx].Push(0) || changed
					}
				}
			}
			NextProd:
		}
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


func followsFor(name string) *set {
	return follows[nonterms[name]]
}





func printSet(s map[int]*set, getName func(i int) string) {
	for i := 0; i < len(s); i++ {
		fmt.Print(getName(i), "[")
		for j, t := range *s[i] {
			if j > 0 {fmt.Print(",")}
			token := t.(int)
			fmt.Print(translateTerm(token))
		}
		fmt.Print("]\n")
	}
	fmt.Print("\n")
}

func translateTerm(i int) (output string) {
	output = "<NOT FOUND>"
	if i == -1 {output = "$"}
	for k, v := range terms {
		if v != i {continue}
		output = k
		if len(output) == 0 {output = "<E>"}
		break
	}
	return
}

func translateNonterm(i int) (output string) {
	output = "<NOT FOUND>"
	for k, v := range nonterms {
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

func printTable() {
	fmt.Print("             ")
	for c := 1; c < len(terms); c++ {
		fmt.Printf("%6s ", translateTerm(c))
	}
	fmt.Print("\n")
	for r := 0; r < len(nonterms); r++ {
		fmt.Printf("%12s ", translateNonterm(r))
		for c := 1; c < len(terms); c++ {
			fmt.Printf("%6d ", table[r][c])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func printTableRaw() {
	fmt.Print("             ")
	for c := 1; c < len(terms); c++ {
		fmt.Printf("%6d ", c)
	}
	fmt.Print("\n")
	for r := 0; r < len(nonterms); r++ {
		fmt.Printf("%12d ", r)
		for c := 1; c < len(terms); c++ {
			fmt.Printf("%6d ", table[r][c])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
