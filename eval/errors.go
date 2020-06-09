package eval

import (
	"fmt"
	obj "mkc/object"
)

func newError(format string, a ...interface{}) *obj.Error {
	return &obj.Error{Message: fmt.Sprintf(format, a...)}
}

func newOErrorUnknownPrefixOp(operator string, right obj.Object) *obj.Error {
	return newError("unknown operator: %s%s", operator, right.Type())
}

func newOErrorUnknownInfixOp(left obj.Object, operator string, right obj.Object) *obj.Error {
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func newOErrorTypeMismatch(left obj.Object, operator string, right obj.Object) *obj.Error {
	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func newOErrorInvalidOperand(operator string, right obj.Object) *obj.Error {
	return newError("invalid operand: %s%s", operator, right.Type())
}

func newOIdentifierError(ident string) *obj.Error {
	return newError("identifier not found: %s", ident)
}

func newOFunctionError(function obj.Object) *obj.Error {
	return newError("not a function: %s", function.Type())
}
