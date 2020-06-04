package parser

import (
	"fmt"
	"mkc/ast"
	"mkc/lexer"
	"mkc/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

////////////////
// Definition //
////////////////

// Returns a new parser
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
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
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// Check if peek is passed type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

///////////////////
// Parser Errors //
//////////////////

// Returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Returns error for peek not being same type as expected
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
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
		return nil
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

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// return EXPRESSION;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	// TODO: We're skipping the expressions until we // encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
