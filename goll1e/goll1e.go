package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

var s scanner

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
		out, err = os.Open(os.Args[2], os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
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
}

func parseGrammar() {
	word, err := s.nextWord()
	for err == nil {
		fmt.Println(word)
		word, err = s.nextWord()
	}
}

