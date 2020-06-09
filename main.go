package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mkc/eval"
	"mkc/lexer"
	obj "mkc/object"
	"mkc/parser"
	"mkc/repl"
	"os"
)

const VERSION = "0.1.0"

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	runFile(flag.Arg(0))
}

func runFile(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf(
			"Error in openning file %s: \n %s",
			fname, err.Error(),
		)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf(
			"Error in reading file %s: \n %s",
			fname, err.Error(),
		)
		return
	}

	l := lexer.New(string(contents))
	p := parser.New(l)
	program := p.ParseProgram()

	env := obj.NewEnvironment()

	evaluated := eval.Eval(program, env)
	fmt.Println(evaluated.Inspect())
}
