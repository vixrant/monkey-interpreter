package parser

import (
	"fmt"
	"mkc/ast"
	"mkc/lexer"
	tk "mkc/token"
	"strconv"
)

type (
	tPrefixParseFn 	func () ast.Expression
	tInfixParseFn 	func(ast.Expression) ast.Expression

	prefixParserTable 	map[tk.TokenType]tPrefixParseFn
	infixParserTable 	map[tk.TokenType]tInfixParseFn
)

type Parser struct {
	l *lexer.Lexer

	currToken tk.Token
	peekToken tk.Token

	errors []string

	prefixParseFns 	prefixParserTable
	infixParseFns 	infixParserTable
}

////////////////
// Definition //
////////////////

// Returns a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},

		prefixParseFns: make(prefixParserTable),
		infixParseFns:	make(infixParserTable),
	}

	// All prefix operators
	p.registerPrefix(tk.IDENTIFIER,	p.parseIdentifier)
	p.registerPrefix(tk.INT,		p.parseIntegerLiteral)
	p.registerPrefix(tk.TRUE,		p.parseBooleanLiteral)
	p.registerPrefix(tk.FALSE,		p.parseBooleanLiteral)
	p.registerPrefix(tk.BANG,		p.parsePrefixExpression)
	p.registerPrefix(tk.MINUS,		p.parsePrefixExpression)
	p.registerPrefix(tk.PLUS,		p.parsePrefixExpression)
	p.registerPrefix(tk.LPAREN,		p.parseGroupedExpression)
	p.registerPrefix(tk.IF,			p.parseIfExpression)
	p.registerPrefix(tk.FUNCTION,	p.parseFunctionLiteral)

	// All infix operators
	p.registerInfix(tk.EQ, 			p.parseInfixExpression)
	p.registerInfix(tk.NOTEQ, 		p.parseInfixExpression)
	p.registerInfix(tk.PLUS, 		p.parseInfixExpression)
	p.registerInfix(tk.MINUS, 		p.parseInfixExpression)
	p.registerInfix(tk.ASTRICK, 	p.parseInfixExpression)
	p.registerInfix(tk.SLASH, 		p.parseInfixExpression)
	p.registerInfix(tk.MOD, 		p.parseInfixExpression)
	p.registerInfix(tk.DASTRICK, 	p.parseInfixExpression)
	p.registerInfix(tk.LT, 			p.parseInfixExpression)
	p.registerInfix(tk.LTEQ, 		p.parseInfixExpression)
	p.registerInfix(tk.GT, 			p.parseInfixExpression)
	p.registerInfix(tk.GTEQ, 		p.parseInfixExpression)
	// Call arguments are like IDENTIFIER ( ARGUMENTS
	p.registerInfix(tk.LPAREN,		p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Registers prefix parse function for a token
func (p *Parser) registerPrefix(tokenType tk.TokenType, fn tPrefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Registers infix parse function for a token
func (p *Parser) registerInfix(tokenType tk.TokenType, fn tInfixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Go to next token
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// If next token is passed type, then move ahead
func (p *Parser) expectPeek(t tk.TokenType) bool {
	if !p.peekTokenIs(t) {
		p.peekError(t)
		return false
	}

	p.nextToken()
	return true
}

// Check if current is passed type
func (p *Parser) currTokenIs(t tk.TokenType) bool {
	return p.currToken.Type == t
}

// Check if peek is passed type
func (p *Parser) peekTokenIs(t tk.TokenType) bool {
	return p.peekToken.Type == t
}

// Returns precedence of current token
func (p *Parser) currPrecedence() pRank {
	if p, ok := precedenceTable[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Returns precedence of peek token
func (p *Parser) peekPrecedence() pRank {
	if p, ok := precedenceTable[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

///////////////////
// Parser Errors //
///////////////////

// Returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Adds error for peek not being same type as expected
func (p *Parser) peekError(t tk.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for parsing
func (p* Parser) typeError(s string, t string) {
	msg := fmt.Sprintf(
		"Cannot parse %s into %s",
		s, t,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for unregistered prefix parse function
func (p* Parser) noPrefixParseFnError(t tk.Token) {
	msg := fmt.Sprintf(
		"No prefix parse function for token %+v",
		t,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for unregistered infix parse function
func (p* Parser) noInfixParseFnError(t tk.Token) {
	msg := fmt.Sprintf(
		"No infix parse function for token %+v",
		t,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for no if condition
func (p* Parser) wrongBracketError(t tk.Token, e string) {
	msg := fmt.Sprintf(
		"expected %s, got: %+v",
		e, t,
	)
	p.errors = append(p.errors, msg)
}

///////////////////
// Parse Program //
///////////////////

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()

	for p.currToken.Type != tk.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

//////////////////
// Parse Blocks //
//////////////////

// Call respective parse procedures
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case tk.LET:
		return p.parseLetStatement()
	case tk.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// let IDENTIFIER = EXPRESSION;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(tk.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(tk.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// return EXPRESSION;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// EXPRESSION
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// { STATEMENTS[] }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	be := &ast.BlockStatement{Token: p.currToken, Statements: []ast.Statement{}}

	p.nextToken()

	for !p.currTokenIs(tk.RBRACE) && !p.currTokenIs(tk.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			continue
		}
		be.Statements = append(be.Statements, stmt)
		p.nextToken()
	}

	return be
}

////////////////////////////////
// Expression Parse Functions //
////////////////////////////////

// Parses any expression
func (p *Parser) parseExpression(precedence pRank) ast.Expression {
	prefix, ok := p.prefixParseFns[p.currToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.currToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(tk.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			p.noInfixParseFnError(p.peekToken)
			return nil
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// IDENTIFIER;
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

// INTEGER
func (p *Parser) parseIntegerLiteral() ast.Expression {
	il := &ast.IntegerLiteral{ Token: p.currToken }

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.typeError(p.currToken.Literal, "integer")
		return nil
	}

	il.Value = value

	return il
}

// BOOLEAN
func (p *Parser) parseBooleanLiteral() ast.Expression {
	il := &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(tk.TRUE)}

	return il
}

// OPERATOR EXPRESSION
func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)

	return pe
}

// EXPRESSION OPERATOR EXPRESSION
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
		Left: left,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	ie.Right = p.parseExpression(precedence)

	return ie
}

// (EXPRESSION)
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(tk.RPAREN) {
		return nil
	}

	return exp
}

// if CONDITION { CONSEQUENT } else { ALTERNATIVE }
func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeek(tk.LPAREN) {
		p.wrongBracketError(p.peekToken, tk.LPAREN)
		return nil
	}

	p.nextToken()
	ie.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tk.RPAREN) {
		p.wrongBracketError(p.peekToken, tk.RPAREN)
		return nil
	}

	if !p.expectPeek(tk.LBRACE) {
		p.wrongBracketError(p.peekToken, tk.LBRACE)
		return nil
	}

	ie.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(tk.ELSE) {
		p.nextToken()

		if !p.expectPeek(tk.LBRACE) {
			p.wrongBracketError(p.peekToken, tk.LBRACE)
			return nil
		}

		ie.Alternative = p.parseBlockStatement()
	}

	return ie
}

// fn (PARAMETERS) { BODY }

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(tk.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tk.LBRACE) {
		return nil
	}
	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier
	if p.peekTokenIs(tk.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(tk.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Go to identifier
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(tk.RPAREN) {
		p.wrongBracketError(p.peekToken, tk.RPAREN)
		return nil
	}

	return identifiers
}

// FUNCTIONLITERAL ( ARGUMENTS )

func (p *Parser) parseCallExpression(fl ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: fl}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var exps []ast.Expression
	if p.peekTokenIs(tk.RPAREN) {
		p.nextToken()
		return exps
	}

	p.nextToken()
	arg := p.parseExpression(LOWEST)
	exps = append(exps, arg)

	for p.peekTokenIs(tk.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Go to expression
		arg := p.parseExpression(LOWEST)
		exps = append(exps, arg)
	}

	if !p.expectPeek(tk.RPAREN) {
		p.wrongBracketError(p.peekToken, tk.RPAREN)
		return nil
	}

	return exps
}
