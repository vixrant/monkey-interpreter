package eval

import (
	"fmt"
	"mkc/ast"
	obj "mkc/object"
)

// Fixed values
var (
	OTRUE	= &obj.Boolean{Value: true}
	OFALSE	= &obj.Boolean{Value: false}
	ONULL	= &obj.Null{}
)

///////////////
// Evaluator //
///////////////

func Eval(node ast.Node, env *obj.Environment) obj.Object {
	switch node := node.(type) {
	// --- start evaluating ---

	// >> Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) { return val }
		env.Set(node.Name.Value, val)

	// >> Expressions

	// data types
	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	// expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) { return left }
		right := Eval(node.Right, env)
		if isError(right) { return right }
		return evalInfixExpression(node.Operator, left, right)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) { return val }
		return &obj.ReturnValue{Value: val}

	// identifiers
	case *ast.Identifier:
		return evalIdentifier(node, env)

	// block constructs
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		body := node.Body
		params := node.Parameters
		return &obj.Function{Body: body, Parameters: params, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) { return function }

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) { return args[0] }

		return applyFunction(function, args)


	// --- end evaluating ---
	default:
		fmt.Printf("Unknown node %+v \n", node)
	}

	return nil
}

///////////////
// Utilities //
///////////////

// Checks if its an error object
func isError(o obj.Object) bool {
	if o == nil {
		return false
	}
	return o.Type() == obj.ERROR_OBJ
}

// Returns the same address of boolean object
func nativeBoolToBooleanObject(input bool) *obj.Boolean {
	if input {
		return OTRUE
	}

	return OFALSE
}

////////////////
// Statements //
////////////////

// Evaluates complete program
func evalProgram(program *ast.Program, env *obj.Environment) obj.Object {
	var result obj.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *obj.ReturnValue:
			return result.Value
		case *obj.Error:
			return result
		}
	}
	return result
}

// Loops over all statements
func evalBlockStatement(block *ast.BlockStatement, env *obj.Environment) obj.Object {
	var result obj.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil  {
			rt := result.Type()
			if rt == obj.RETURN_OBJ || rt == obj.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// Matches operator with required function call
func evalPrefixExpression(operator string, right obj.Object) obj.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	default:
		return newOErrorUnknownPrefixOp(operator, right)
	}
}

//////////////////
// Prefix Exprs //
//////////////////

// Matches right object with supported data types
func evalBangOperatorExpression(right obj.Object) obj.Object {
	switch right {
	case OTRUE:
		return OFALSE
	case OFALSE:
		return OTRUE
	default:
		return newOErrorInvalidOperand("!", right)
	}
}

// Returns - of given right expression
func evalMinusPrefixOperatorExpression(right obj.Object) obj.Object {
	if right.Type() != obj.INTEGER_OBJ {
		return newOErrorInvalidOperand("-", right)
	}

	value := - right.(*obj.Integer).Value
	return &obj.Integer{Value: value}
}

// Returns + of given right expression
func evalPlusPrefixOperatorExpression(right obj.Object) obj.Object {
	if right.Type() != obj.INTEGER_OBJ {
		return newOErrorInvalidOperand("+", right)
	}

	return right
}

/////////////////
// Infix Exprs //
/////////////////

// Passes infix expression to respective handlers
func evalInfixExpression(operator string, left obj.Object, right obj.Object) obj.Object {
	switch {
	case left.Type() == obj.INTEGER_OBJ && right.Type() == obj.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() != right.Type():
		return newOErrorTypeMismatch(left, operator, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return newOErrorUnknownInfixOp(left, operator, right)
	}
}

// Evaluates an arithmetic operation
func evalIntegerInfixExpression(operator string, left obj.Object, right obj.Object) obj.Object {
	lval := left.(*obj.Integer).Value
	rval := right.(*obj.Integer).Value

	switch operator {
	case "+":
		return &obj.Integer{Value: lval + rval}

	case "-":
		return &obj.Integer{Value: lval - rval}

	case "*":
		return &obj.Integer{Value: lval * rval}

	case "/":
		return &obj.Integer{Value: lval / rval}

	case "%":
		return &obj.Integer{Value: lval % rval}

	case "**":
		res := int64(1)
		for i := int64(0); i < rval; i++ {
			res = res * lval
		}
		return &obj.Integer{Value: res}

	case "<":
		return nativeBoolToBooleanObject(lval < rval)

	case "<=":
		return nativeBoolToBooleanObject(lval <= rval)

	case ">":
		return nativeBoolToBooleanObject(lval > rval)

	case ">=":
		return nativeBoolToBooleanObject(lval >= rval)

	case "==":
		return nativeBoolToBooleanObject(lval == rval)

	case "!=":
		return nativeBoolToBooleanObject(lval != rval)

	default:
		return newOErrorUnknownInfixOp(left, operator, right)
	}
}

////////////
// Others //
////////////

// Handles an if else expression
func evalIfExpression(ie *ast.IfExpression, env *obj.Environment) obj.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) { return condition }

	if condition == OTRUE {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return ONULL
	}
}

// Returns identifier object from environment
func evalIdentifier(ie *ast.Identifier, env *obj.Environment) obj.Object {
	val, ok := env.Get(ie.Value)
	if !ok {
		return newOIdentifierError(ie.Value)
	}

	return val
}

// Loops over multiple expressions
func evalExpressions(exps []ast.Expression, env *obj.Environment) []obj.Object {
	var result []obj.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) { return []obj.Object{evaluated} }
		result = append(result, evaluated)
	}

	return result
}

// Evaluates a function
func applyFunction(fnObj obj.Object, args []obj.Object) obj.Object {
	function, ok := fnObj.(*obj.Function)
	if !ok {
		return newOFunctionError(fnObj)
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

// Extends the env with function arguments and returns wrapped env
func extendFunctionEnv(fn *obj.Function, args []obj.Object) *obj.Environment {
	env := obj.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// Unwraps the return value
func unwrapReturnValue(o obj.Object) obj.Object {
	if returnValue, ok := o.(*obj.ReturnValue); ok {
		return returnValue.Value
	}

	return o
}
