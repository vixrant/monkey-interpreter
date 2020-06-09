package parser

import (
	"fmt"
	"mkc/token"
	"testing"

	"mkc/ast"
	"mkc/lexer"
)

///////////////////////
// Utility functions //
///////////////////////

func getAST(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParserErrors(t, p)

	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

////////////////////////
// Test smaller nodes //
////////////////////////

func assertLetStatement(t *testing.T, s ast.Statement, name string, value int64) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Fatalf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	
	if !assertLiteralExpression(t, letStmt.Value, value) {
		return false
	}

	return true
}

func assertIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf(
			"exp not *ast.IntegerLiteral. got=%T",
			exp,
		)
		return false
	}

	if integ.Token.Type != tk.INT {
		t.Fatalf(
			"ident.Token is not token.INT, got=%T",
			integ.Token,
		)
		return false
	}

	if integ.Value != value {
		t.Fatalf(
			"integ.Value not %d. got=%d",
			value, integ.Value,
		)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf(
			"integ.TokenLiteral not %d. got=%s",
			value, integ.TokenLiteral(),
		)
		return false
	}

	return true
}

func assertBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bl, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Fatalf(
			"exp not *ast.BooleanLiteral. got=%T",
			exp,
		)
		return false
	}

	if bl.Token.Type != tk.TRUE && bl.Token.Type != tk.FALSE {
		t.Fatalf("ident.Token is not token.TRUE or token.FALSE, got=%T",
			bl.Token,
		)
	}

	if bl.Value != value {
		t.Fatalf(
			"integ.Value not %t. got=%t",
			value, bl.Value,
		)
		return false
	}

	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf(
			"integ.TokenLiteral not %t. got=%s",
			value, bl.TokenLiteral(),
		)
		return false
	}

	return true
}

func assertIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifiter, got=%T",
			exp,
		)
		return false
	}

	if ident.Token.Type != tk.IDENTIFIER {
		t.Fatalf("ident.Token is not ast.IDENTIFER, got=%T",
			ident.Token,
		)
		return false
	}

	if ident.Value != value {
		t.Fatalf("ident.Value is not %s, got=%s",
			value, ident.Value,
		)
		return false
	}

	return true
}

func assertLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return assertIntegerLiteral(t, exp, int64(v))
	case int64:
		return assertIntegerLiteral(t, exp, v)
	case bool:
		return assertBooleanLiteral(t, exp, v)
	case string:
		return assertIdentifier(t, exp, v)
	default:
		t.Fatalf("type of exp not handled. got=%T", exp)
		return false
	}
}

func assertInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !assertLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !assertLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

///////////
// Tests //
///////////

// let statement

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	program := getAST(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements),
		)
	}

	tests := []struct {
		expectedIdentifier string
		expectedValue      int64
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !assertLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}

// return statement

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 10201;
	`;

	program := getAST(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements),
		)
	}

	tests := []struct {
		expectedValue int64
	}{
		{5},
		{10},
		{10201},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.Token.Type != tk.RETURN {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

		if !assertLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

// identifier expression statement

func TestIdentifierExpressionStatement(t *testing.T) {
	input := "foobar;"

	program := getAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got=%d",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0],
		)
	}

	if !assertIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}

// integer literal statement

func TestIntegerLiteralStatement(t *testing.T) {
	input := "5;"

	program := getAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got=%d",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
			program.Statements[0],
		)
	}

	if !assertIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

// boolean literal statement

func TestBooleanLiteralStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		program := getAST(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have enough statements, got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T",
				program.Statements[0],
			)
		}

		if !assertBooleanLiteral(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

// prefix expression parser

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"-5", "-", 5},
		{"!15", "!", 15},
		{"+20643", "+", 20643},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		program := getAST(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"stmt is not ast.PrefixExpression. got=%T",
				stmt.Expression,
			)
		}

		if exp.Operator != tt.operator {
			t.Fatalf(
				"exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator,
			)
		}

		if !assertLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

// infix expression

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 ** 5;", 5, "**", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := getAST(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		if !assertInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

// operator precedence

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"!+a", "(!(+a))"},
		{"a+b+c", "((a + b) + c)"},
		{"a+b-c", "((a + b) - c)"},
		{"a*b*c", "((a * b) * c)"},
		{"a*b/c", "((a * b) / c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3+4;-5*5", "(3 + 4)((-5) * 5)"},
		{"5>4==3<4", "((5 > 4) == (3 < 4))"},
		{"5<4!=3>4", "((5 < 4) != (3 > 4))"},
		{"3+4*5==3*1+4*5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"2**3", "(2 ** 3)"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))",},
		{"-(5 + 5)", "(-(5 + 5))", },
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))","add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		program := getAST(t, tt.input)

		actual := program.String()
		if actual != tt.expected {
			t.Fatalf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// if statements

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := getAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression,
		)
	}

	if !assertInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf(
			"consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements),
		)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0],
		)
	}

	if !assertIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0],
		)
	}

	if !assertIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// function literals

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y }`

	program := getAST(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf(
			"stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression,
		)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf(
			"function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters),
		)
	}
	assertLiteralExpression(t, function.Parameters[0], "x")
	assertLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements),
		)
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0],
		)
	}

	assertInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input			string
		expectedParams 	[]string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		program := getAST(t, tt.input)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fl := stmt.Expression.(*ast.FunctionLiteral)

		if len(fl.Parameters) != len(tt.expectedParams) {
			t.Fatalf(
				"length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(fl.Parameters),
			)
		}

		for i, ident := range tt.expectedParams {
			assertLiteralExpression(t, fl.Parameters[i], ident)
		}
	}
}

// function calls

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	program := getAST(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	ce, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf(
			"stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression,
		)
	}

	if !assertIdentifier(t, ce.Function, "add") {
		return
	}

	if len(ce.Arguments) != 3 {
		t.Fatalf(
			"ce.Arguments does not contain %d expressions, got=%d",
			3, len(ce.Arguments),
		)
	}

	assertLiteralExpression(t, ce.Arguments[0], 1)
	assertInfixExpression(t, ce.Arguments[1], 2, "*", 3)
	assertInfixExpression(t, ce.Arguments[2], 4, "+", 5)
}
