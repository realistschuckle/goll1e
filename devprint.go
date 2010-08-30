package main

import (
	"fmt"
)

func printTypedEntries() {
	for k, v := range typedEntries {
		fmt.Println(k, "is a", v)
	}
	fmt.Print("\n")
}

func printSet(s map[int]*set, getName func(i int) string) {
	for i := 0; i < len(s); i++ {
		fmt.Print(getName(i), "[")
		for j, t := range *s[i] {
			if j > 0 {
				fmt.Print(",")
			}
			token := t.(int)
			fmt.Print(translateTerm(token))
		}
		fmt.Print("]\n")
	}
	fmt.Print("\n")
}

func translateTerm(i int) (output string) {
	output = "<NOT FOUND>"
	if i == -1 {
		output = "$"
	}
	for k, v := range terms {
		if v != i {
			continue
		}
		output = k
		if len(output) == 0 {
			output = "<E>"
		}
		break
	}
	return
}

func translateNonterm(i int) (output string) {
	output = "<NOT FOUND>"
	for k, v := range nonterms {
		if v != i {
			continue
		}
		output = k
		break
	}
	return
}

func translateTerms(s *set) (output string) {
	output = "["
	for i, v := range *s {
		token := v.(int)
		if i > 0 {
			output += ","
		}
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
