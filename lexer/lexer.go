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
func New(input string) *Lexer {
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
func newToken(tokenType tk.TokenType, ch byte) tk.Token {
	return tk.Token{Type: tokenType, Literal: string(ch)}
}

// Function to return a new token from a byte array
func newTokenString(tokenType tk.TokenType, s string) tk.Token {
	return tk.Token{Type: tokenType, Literal: s}
}

///////////////
// NextToken //
///////////////

// Returns next token in input stream
func (l *Lexer) NextToken() tk.Token {
	var tok tk.Token

	l.eatWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(tk.PLUS, l.ch)

	case '-':
		tok = newToken(tk.MINUS, l.ch)

	case '*':
		tok = newToken(tk.ASTRICK, l.ch)

		if l.peekChar() == '*' {
			s := l.readString(2)
			tok = newTokenString(tk.DASTRICK, s)
		}

	case '/':
		tok = newToken(tk.SLASH, l.ch)

	case '%':
		tok = newToken(tk.MOD, l.ch)

	case '=':
		tok = newToken(tk.ASSIGN, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(tk.EQ, s)
		}

	case '!':
		tok = newToken(tk.BANG, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(tk.NOTEQ, s)
		}

	case '<':
		tok = newToken(tk.LT, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(tk.LTEQ, s)
		}

	case '>':
		tok = newToken(tk.GT, l.ch)

		if l.peekChar() == '=' {
			s := l.readString(2)
			tok = newTokenString(tk.GTEQ, s)
		}

	case '(':
		tok = newToken(tk.LPAREN, l.ch)

	case ')':
		tok = newToken(tk.RPAREN, l.ch)

	case '{':
		tok = newToken(tk.LBRACE, l.ch)

	case '}':
		tok = newToken(tk.RBRACE, l.ch)

	case ',':
		tok = newToken(tk.COMMA, l.ch)

	case ';':
		tok = newToken(tk.SEMICOLON, l.ch)

	case 0:
		tok = newTokenString(tk.EOF, "")

	default:
		if isLetter(l.ch) {
			s := l.readIdentifier()
			t := tk.LookupIdent(s)
			tok := newTokenString(t, s)
			return tok
		}

		if isDigit(l.ch) {
			n := l.readNumber()
			tok := newTokenString(tk.INT, n)
			return tok
		}

		tok = newToken(tk.ILLEGAL, l.ch)
	}

	l.readChar()

	return tok
}


