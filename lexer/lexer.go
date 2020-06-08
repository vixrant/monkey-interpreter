package lexer

import "mkc/token"

type Lexer struct {
	input			string
	position 		int // current position
	readPosition 	int // after current position
	ch 				byte
}

////////////////
// Definition //
////////////////

// Creates a lexer for given input string
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Reads the next character in input
func (l *Lexer) readChar() {
	l.ch = 0 // eof

	if l.readPosition < len(l.input) {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// Returns the next character in input
// Doesn't affect pointer
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// Returns bunch of characters, starting with the current one
// Provide length of string
func (l *Lexer) readString(n int) string {
	ch := l.ch
	ba := []byte{ ch }
	n -= 1
	for i := 0; i < n; i++ {
		l.readChar()
		ba = append(ba, l.ch)
	}
	return string(ba)
}

// Returns numerical literal
func (l *Lexer) readNumber() string {
	startPosition := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

// Returns identifier, that is, the alphanumeric sequence till next end
func (l* Lexer) readIdentifier() string {
	startPosition := l.position
	for isLegalIdentChar(l.ch) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

// Skips over all whitespace
func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

///////////////
// Utilities //
///////////////

// Function to return a new token from a byte
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Function to return a new token from a byte array
func newTokenString(tokenType token.TokenType, s string) token.Token {
	return token.Token{Type: tokenType, Literal: s}
}

///////////////
// NextToken //
///////////////

// Returns next token in input stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch)

	case '-':
		tok = newToken(token.MINUS, l.ch)

	case '*':
		tok = newToken(token.ASTRICK, l.ch)

		if l.peekChar() == '*' {
			s := l.readString(2)
			tok = newTokenString(token.DASTRICK, s)
		}

	case '/':
		tok = newToken(token.SLASH, l.ch)

	case '=':
		tok = newToken(token.ASSIGN, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(token.EQ, s)
		}

	case '!':
		tok = newToken(token.BANG, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(token.NOTEQ, s)
		}

	case '<':
		tok = newToken(token.LT, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(token.LTEQ, s)
		}

	case '>':
		tok = newToken(token.GT, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(token.GTEQ, s)
		}

	case '(':
		tok = newToken(token.LPAREN, l.ch)

	case ')':
		tok = newToken(token.RPAREN, l.ch)

	case '{':
		tok = newToken(token.LBRACE, l.ch)

	case '}':
		tok = newToken(token.RBRACE, l.ch)

	case ',':
		tok = newToken(token.COMMA, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)

	case 0:
		tok = newTokenString(token.EOF, "")

	default:
		if isLetter(l.ch) {
			s := l.readIdentifier()
			t := token.LookupIdent(s)
			tok := newTokenString(t, s)
			return tok
		}

		if isDigit(l.ch) {
			n := l.readNumber()
			tok := newTokenString(token.INT, n)
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()

	return tok
}


