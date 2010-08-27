package main

import (
	"fmt"
	"os"
)

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
					if table[r][c] != -1 {
						fmt.Println("FIRST/FIRST CONFLICT")
						fmt.Println("Error between", prods[table[r][c]], "and", prod)
						os.Exit(-4)
					}
					table[r][c] = prodidx
				case firsts[prodidx].HasE() && follows[r].IndexOf(c) > -1:
					if table[r][c] != -1 {
						fmt.Println("FIRST/FOLLOW CONFLICT")
						fmt.Println("Error between", prods[table[r][c]], "and", prod)
						os.Exit(-4)
					}
					table[r][c] = prodidx
				}
			}
		}
	}
}

func computeFollows() {
	if len(nonterms) == 0 {return}
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
