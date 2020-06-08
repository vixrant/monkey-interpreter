package parser

import "mkc/token"

type pRank uint8

const (
	_ pRank = iota
	LOWEST
	EQUALS		// ==
	LESSGREATER // >, <, <=, >=
	SUM			// + -
	PRODUCT		// * /
	POWER		// **
	PREFIX		// -X or !X
	CALL		// myFunction(X)
)

var precedenceTable = map[token.TokenType]pRank{
	token.EQ: 		EQUALS,
	token.NOTEQ: 	EQUALS,
	token.PLUS: 	SUM,
	token.MINUS: 	SUM,
	token.ASTRICK: 	PRODUCT,
	token.SLASH: 	PRODUCT,
	token.DASTRICK: POWER,
	token.LT: 		LESSGREATER,
	token.LTEQ: 	LESSGREATER,
	token.GT: 		LESSGREATER,
	token.GTEQ: 	LESSGREATER,
	token.LPAREN: 	CALL,
}
