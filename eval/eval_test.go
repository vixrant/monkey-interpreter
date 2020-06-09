package eval

import (
	"mkc/lexer"
	obj "mkc/object"
	"mkc/parser"
	"testing"
)

///////////////
// Utilities //
///////////////

func runEval(t *testing.T, input string) obj.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}

		t.FailNow()
	}

	env := obj.NewEnvironment()

	return Eval(program, env)
}

func assertOInteger(t *testing.T, o obj.Object, expected int64) bool {
	result, ok := o.(*obj.Integer)
	if !ok {
		t.Errorf(
			"object is not Integer. got=%T (%+v)",
			o, o,
		)
		return false
	}

	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%d, want=%d",
			result.Value, expected,
		)
		return false
	}

	return true
}

func assertOBoolean(t *testing.T, o obj.Object, expected bool) bool {
	result, ok := o.(*obj.Boolean)
	if !ok {
		t.Errorf(
			"object is not Boolean. got=%T (%+v)",
			o, o,
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func assertNullObject(t *testing.T, o obj.Object) bool {
	if o != ONULL {
		t.Errorf(
			"object is not NULL. got=%T (%+v)",
			o, o,
		)
		return false
	}
	return true
}

///////////
// Tests //
///////////

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 -10", 10},
		{"2 * 2 * 2 * 2 * 2", 32}, {"-50 + 100 + -50", 0}, {"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5 ** 3", 125},
		{"5 ** 3 * 2 + 6", 256},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		assertOInteger(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for idx, tt := range tests {
		evaluated := runEval(t, tt.input)
		res := assertOBoolean(t, evaluated, tt.expected)
		if !res {
			t.Log("Test ", idx)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		assertOBoolean(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests :=  []struct {
		input		string
		expected	interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			assertOInteger(t, evaluated, int64(integer))
		} else {
			assertNullObject(t, evaluated)
		}
	}
}

func TestReturnExpression(t *testing.T) {
	tests :=  []struct {
		input		string
		expected	int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"2 + 5; return 10;", 10},
		{"2 ** 4; return 10; 9;", 10},
		{"return 5*2; 9;", 10},
		{"return (2*9-8);", 10},
		{`
			if (10 > 1) {
     			if (10 > 1) {
					return 10;
				}
				return 1;
			}`, 10},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)

		if !assertOInteger(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input 		string
		expected	string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"invalid operand: -BOOLEAN"},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobarhoohaa",
			"identifier not found: foobarhoohaa",
		},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		errObj, ok := evaluated.(*obj.Error)

		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf(
				"wrong error message. expected=%q, got=%q",
				tt.expected, errObj.Message,
			)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		assertOInteger(t, evaluated, tt.expected)
	}
}


func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := runEval(t, input)

	fn, ok := evaluated.(*obj.Function)
	if !ok {
		t.Fatalf(
			"object is not Function. got=%T (%+v)",
			evaluated, evaluated,
		)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf(
			"function has wrong parameters. Parameters=%+v",
			fn.Parameters,
		)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf(
			"parameter is not 'x'. got=%q",
			fn.Parameters[0],
		)
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf(
			"body is not %q. got=%q",
			expectedBody, fn.Body.String(),
		)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		evaluated := runEval(t, tt.input)
		assertOInteger(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) {
     		fn(y) { x + y };
		};
		let addTwo = newAdder(2); addTwo(2);
	`

	evaluated := runEval(t, input)
	assertOInteger(t, evaluated, 4)
}
