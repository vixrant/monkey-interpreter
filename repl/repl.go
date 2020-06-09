package repl

import (
	"io"
	"mkc/eval"
	"mkc/lexer"
	obj "mkc/object"
	"mkc/parser"
)

const PROMPT = "-> "

func Start(in io.Reader, out io.Writer) {
	rio := SetupIO(in, out)
	env := obj.NewEnvironment()

	for {
		rio.Write(PROMPT)

		line := rio.Read()

		switch line {
		case "":
			continue
		case ".exit":
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			rio.Write("Parser errors: \n")
			for _, msg := range p.Errors() {
				rio.Write("\t" + msg + "\n")
			}
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			rio.Write(evaluated.Inspect())
			rio.Write("\n")
		}
	}
}
