package lexer

import "mkc/token"

type expectations struct {
	expectedType    tk.TokenType
	expectedLiteral string
}

type testCase struct {
	input	string
	expect 	[]expectations
}

///////////
// Tests //
///////////

var lexerTestCases = map[string]testCase {
	"algebra": {
		input: `
			let five = 5;
			let ten = 10;
			let add = fn(x, y) {
				x + y;
			};
			let result = add(five, ten);
		`,
		expect: []expectations{
			{tk.LET, "let"},
			{tk.IDENTIFIER, "five"},
			{tk.ASSIGN, "="},
			{tk.INT, "5"},
			{tk.SEMICOLON, ";"},
			{tk.LET, "let"},
			{tk.IDENTIFIER, "ten"},
			{tk.ASSIGN, "="},
			{tk.INT, "10"},
			{tk.SEMICOLON, ";"},
			{tk.LET, "let"},
			{tk.IDENTIFIER, "add"},
			{tk.ASSIGN, "="},
			{tk.FUNCTION, "fn"},
			{tk.LPAREN, "("},
			{tk.IDENTIFIER, "x"},
			{tk.COMMA, ","},
			{tk.IDENTIFIER, "y"},
			{tk.RPAREN, ")"},
			{tk.LBRACE, "{"},
			{tk.IDENTIFIER, "x"},
			{tk.PLUS, "+"},
			{tk.IDENTIFIER, "y"},
			{tk.SEMICOLON, ";"},
			{tk.RBRACE, "}"},
			{tk.SEMICOLON, ";"},
			{tk.LET, "let"},
			{tk.IDENTIFIER, "result"},
			{tk.ASSIGN, "="},
			{tk.IDENTIFIER, "add"},
			{tk.LPAREN, "("},
			{tk.IDENTIFIER, "five"},
			{tk.COMMA, ","},
			{tk.IDENTIFIER, "ten"},
			{tk.RPAREN, ")"},
			{tk.SEMICOLON, ";"},
			{tk.EOF, ""},
		},
	},

	"logic": {
		input: `
			!-/*5;
			5 < 10 > 5;
			true;false;
		`,
		expect: []expectations{
			{tk.BANG, "!"},
			{tk.MINUS, "-"},
			{tk.SLASH, "/"},
			{tk.ASTRICK, "*"},
			{tk.INT, "5"},
			{tk.SEMICOLON, ";"},
			{tk.INT, "5"},
			{tk.LT, "<"},
			{tk.INT, "10"},
			{tk.GT, ">"},
			{tk.INT, "5"},
			{tk.SEMICOLON, ";"},
			{tk.TRUE, "true"},
			{tk.SEMICOLON, ";"},
			{tk.FALSE, "false"},
			{tk.SEMICOLON, ";"},
		},
	},

	"if-else": {
		input: `
			if (5 < 10) {
				return true;
			} else {
				return false;
			}
		`,
		expect: []expectations{
			{tk.IF, "if"},
			{tk.LPAREN, "("},
			{tk.INT, "5"},
			{tk.LT, "<"},
			{tk.INT, "10"},
			{tk.RPAREN, ")"},
			{tk.LBRACE, "{"},
			{tk.RETURN, "return"},
			{tk.TRUE, "true"},
			{tk.SEMICOLON, ";"},
			{tk.RBRACE, "}"},
			{tk.ELSE, "else"},
			{tk.LBRACE, "{"},
			{tk.RETURN, "return"},
			{tk.FALSE, "false"},
			{tk.SEMICOLON, ";"},
			{tk.RBRACE, "}"},
		},
	},

	"double-symbols": {
		input: `
			if x == 5
				return x != 7
		`,
		expect: []expectations{
			{tk.IF, "if"},
			{tk.IDENTIFIER, "x"},
			{tk.EQ, "=="},
			{tk.INT, "5"},
			{tk.RETURN, "return"},
			{tk.IDENTIFIER, "x"},
			{tk.NOTEQ, "!="},
			{tk.INT, "7"},
		},
	},
}
