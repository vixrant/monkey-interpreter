package ast

////////////////
// Interfaces //
////////////////

// An AST node
type Node interface {
	TokenLiteral()	string
	String() 		string
}

// Node type - Statement
type Statement interface {
	Node
	statementNode()
}

// Node type - Expression
type Expression interface {
	Node
	expressionNode()
}
