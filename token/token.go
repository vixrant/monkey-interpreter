package token

// Data structures

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// Vocabulary

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// Arithmetic
	ASSIGN		= "="
	PLUS		= "+"
	MINUS		= "-"
	SLASH		= "/"
	BANG		= "!"
	ASTRICK		= "*"
	DASTRICK	= "**"

	// Relational
	LT    = "<"
	GT    = ">"
	LTEQ  = "<="
	GTEQ  = ">="
	EQ    = "=="
	NOTEQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION	= "FUNCTION"
	LET			= "LET"
	IF			= "IF"
	ELSE		= "ELSE"
	FOR			= "FOR"
	RETURN		= "RETURN"
	TRUE		= "TRUE"
	FALSE		= "FALSE"
)

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"for": FOR,
	"return": RETURN,
}

// Checks if supposed identifier is a keyword
// If true, returns the keyword token
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
