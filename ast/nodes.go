package ast

import (
	"mkc/token"
)

/////////////
// Structs //
/////////////

// Root node of AST

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	program := &Program{
		Statements: []Statement{},
	}
	return program
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// identifier node

type Identifier struct {
	Token	token.Token
	Value 	string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}

// let statement node

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) statementNode() {}

// return statement

type ReturnStatement struct {
	Token		token.Token
	ReturnValue	Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) statementNode() {}
