package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	for name, cc := range lexerTestCases {
		l := NewLexer(cc.input)

		for i, tt := range cc.expect {
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Fatalf("tests %s [%d] - token Type wrong. expected = %q. got = %q, literal=%s",
					name, i, tt.expectedType, tok.Type, tok.Literal)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests %s [%d] - token Literal wrong. expected = %q. got = %q",
					name, i, tt.expectedLiteral, tok.Literal)
			}
		}
	}
}
