package repl

import (
	"bufio"
	"fmt"
	"io"
	"mkc/lexer"
	token "mkc/tokens"
)

const PROMPT = "-> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)

		scan := scanner.Scan()
		if !scan {
			return
		}

		line := scanner.Text()

		l := lexer.NewLexer(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
