package repl

import (
	"bufio"
	"fmt"
	"io"
	"mkc/lexer"
	"mkc/parser"
)

const PROMPT = "-> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Println("Parser errors:")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
