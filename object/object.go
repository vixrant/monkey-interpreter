package object

import (
	"bytes"
	"fmt"
	"mkc/ast"
	"strings"
)

////////////////
// Interfaces //
////////////////

type ObjectType string

// Foundation of our object system
type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ 	= "INTEGER"
	BOOLEAN_OBJ 	= "BOOLEAN"
	NULL_OBJ 		= "NULL"
	ERROR_OBJ		= "ERROR"
	RETURN_OBJ		= "RETURN"
	FUNCTION_OBJ	= "FUNCTION"
)

/////////////
// Objects //
/////////////

// Data types

type Integer struct {
	Value	int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }


type Boolean struct {
	Value	bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }


type Null struct {}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Functions

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body		*ast.BlockStatement
	Env 		*Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string  {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Errors

type Error struct {
	Message	string
}

func (e *Error) Inspect() string { return "Error" + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
