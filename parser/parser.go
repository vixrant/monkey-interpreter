package parser

import (
	"fmt"
	"mkc/ast"
	"mkc/lexer"
	"mkc/token"
	"strconv"
)

type (
	tPrefixParseFn 	func () ast.Expression
	tInfixParseFn 	func(ast.Expression) ast.Expression

	prefixParserTable 	map[token.TokenType]tPrefixParseFn
	infixParserTable 	map[token.TokenType]tInfixParseFn
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string

	prefixParseFns 	prefixParserTable
	infixParseFns 	infixParserTable
}

////////////////
// Definition //
////////////////

// Returns a new parser
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},

		prefixParseFns: make(prefixParserTable),
		infixParseFns:	make(infixParserTable),
	}

	// All prefix operators
	p.registerPrefix(token.IDENTIFIER,	p.parseIdentifier)
	p.registerPrefix(token.INT,			p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE,		p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE,		p.parseBooleanLiteral)
	p.registerPrefix(token.BANG,		p.parsePrefixExpression)
	p.registerPrefix(token.MINUS,		p.parsePrefixExpression)
	p.registerPrefix(token.PLUS,		p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN,		p.parseGroupedExpression)
	p.registerPrefix(token.IF,			p.parseIfExpression)
	p.registerPrefix(token.FUNCTION,	p.parseFunctionLiteral)

	// All infix operators
	p.registerInfix(token.EQ, 			p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, 		p.parseInfixExpression)
	p.registerInfix(token.PLUS, 		p.parseInfixExpression)
	p.registerInfix(token.MINUS, 		p.parseInfixExpression)
	p.registerInfix(token.ASTRICK, 		p.parseInfixExpression)
	p.registerInfix(token.SLASH, 		p.parseInfixExpression)
	p.registerInfix(token.DASTRICK, 	p.parseInfixExpression)
	p.registerInfix(token.LT, 			p.parseInfixExpression)
	p.registerInfix(token.LTEQ, 		p.parseInfixExpression)
	p.registerInfix(token.GT, 			p.parseInfixExpression)
	p.registerInfix(token.GTEQ, 		p.parseInfixExpression)
	// Call arguments are like IDENTIFIER ( ARGUMENTS
	p.registerInfix(token.LPAREN,		p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Registers prefix parse function for a token
func (p *Parser) registerPrefix(tokenType token.TokenType, fn tPrefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Registers infix parse function for a token
func (p *Parser) registerInfix(tokenType token.TokenType, fn tInfixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Go to next token
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// If next token is passed type, then move ahead
func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.peekTokenIs(t) {
		p.peekError(t)
		return false
	}

	p.nextToken()
	return true
}

// Check if current is passed type
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// Check if peek is passed type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
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
func (p *Parser) peekError(t token.TokenType) {
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
func (p* Parser) noPrefixParseFnError(t token.Token) {
	msg := fmt.Sprintf(
		"No prefix parse function for token %+v",
		t,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for unregistered infix parse function
func (p* Parser) noInfixParseFnError(t token.Token) {
	msg := fmt.Sprintf(
		"No infix parse function for token %+v",
		t,
	)
	p.errors = append(p.errors, msg)
}

// Adds error for no if condition
func (p* Parser) wrongBracketError(t token.Token, e string) {
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

	for p.currToken.Type != token.EOF {
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
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// let IDENTIFIER = EXPRESSION;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// return EXPRESSION;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// EXPRESSION
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// { STATEMENTS[] }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	be := &ast.BlockStatement{Token: p.currToken, Statements: []ast.Statement{}}

	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
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

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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
	il := &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}

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

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// if CONDITION { CONSEQUENT } else { ALTERNATIVE }
func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		p.wrongBracketError(p.peekToken, token.LPAREN)
		return nil
	}

	p.nextToken()
	ie.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		p.wrongBracketError(p.peekToken, token.RPAREN)
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		p.wrongBracketError(p.peekToken, token.LBRACE)
		return nil
	}

	ie.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			p.wrongBracketError(p.peekToken, token.LBRACE)
			return nil
		}

		ie.Alternative = p.parseBlockStatement()
	}

	return ie
}

// fn (PARAMETERS) { BODY }

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Go to identifier
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		p.wrongBracketError(p.peekToken, token.RPAREN)
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
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return exps
	}

	p.nextToken()
	arg := p.parseExpression(LOWEST)
	exps = append(exps, arg)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Skip comma
		p.nextToken() // Go to expression
		arg := p.parseExpression(LOWEST)
		exps = append(exps, arg)
	}

	if !p.expectPeek(token.RPAREN) {
		p.wrongBracketError(p.peekToken, token.RPAREN)
		return nil
	}

	return exps
}
