package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	var tcases = []testCase{
		testAlgebra,
		testLogic,
		testIfElse,
		testDoubleSymbols,
	}

	for cidx, cc := range tcases {
		l := NewLexer(cc.input)

		for i, tt := range cc.expect {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests %d [%d] - token Type wrong. expected = %q. got = %q, literal=%s",
					cidx, i, tt.expectedType, tok.Type, tok.Literal)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests %d [%d] - token Literal wrong. expected = %q. got = %q",
					cidx, i, tt.expectedLiteral, tok.Literal)
			}
		}
	}
}
