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
	MOD			// %
	POWER		// **
	PREFIX		// -X or !X
	CALL		// myFunction(X)
)

var precedenceTable = map[tk.TokenType]pRank{
	tk.PLUS:     SUM,
	tk.MINUS:    SUM,
	tk.ASTRICK:  PRODUCT,
	tk.SLASH:    PRODUCT,
	tk.MOD:    	 MOD,
	tk.DASTRICK: POWER,
	tk.EQ:       EQUALS,
	tk.NOTEQ:    EQUALS,
	tk.LT:       LESSGREATER,
	tk.LTEQ:     LESSGREATER,
	tk.GT:       LESSGREATER,
	tk.GTEQ:     LESSGREATER,
	tk.LPAREN:   CALL,
}
