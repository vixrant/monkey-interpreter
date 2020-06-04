package lexer

import token "mkc/tokens"

type testExpects struct {
	expectedType	token.TokenType
	expectedLiteral string
}

type testCase struct {
	input	string
	expect 	[]testExpects
}

///////////
// Tests //
///////////

var testAlgebra = testCase{
	input: `
		let five = 5;
		let ten = 10;
		let add = fn(x, y) {
			x + y;
		};
		let result = add(five, ten);
	`,
	expect: []testExpects{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	},
}

var testLogic = testCase{
	input: `
		!-/*5;
		5 < 10 > 5;
		true;false;
	`,
	expect: []testExpects{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTRICK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
	},
}

var testIfElse = testCase{
	input: `
		if (5 < 10) {
			return true;
		} else {
			return false;
		}
	`,
	expect: []testExpects{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
	},
}

var testDoubleSymbols = testCase {
	input: `
		if x == 5
			return x != 7
	`,
	expect: []testExpects{
		{token.IF, "if"},
		{token.IDENTIFIER, "x"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "x"},
		{token.NOTEQ, "!="},
		{token.INT, "7"},
	},
}
