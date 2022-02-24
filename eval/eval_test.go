package eval

import (
	"fmt"
	"testing"

	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/obj"
	"github.com/Revolyssup/monkey/parser"
)

func TestEvalIntExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"100", 100},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evalObj := testEval(tt.input)
		testingIntObject(t, evalObj, tt.expected)
	}
}
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
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
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}
func testEval(input string) obj.Object {
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	env := obj.NewEnvironment()
	return Eval(prog, env)
}
func testBooleanObject(t *testing.T, object obj.Object, expected bool) bool {
	result, ok := object.(*obj.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", object, object)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testingIntObject(t *testing.T, object obj.Object, expected int64) bool {
	result, ok := object.(*obj.Integer)
	if !ok {
		t.Errorf("Object is not Integer. got=%T", object)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func testNullObject(t *testing.T, object obj.Object) bool {
	if object != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", object, object)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, object obj.Object, expected int64) bool {
	result, ok := object.(*obj.Integer)
	if !ok {
		fmt.Println("EX", expected, " EY ", object.DataType(), object.Inspect())
		t.Errorf("object is not Integer. got=%T (%+v)", object, object)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
		
			if (10 > 1) {
			return 10;
			}
			return 1;
		}`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: Integer + Bool",
		},
		{
			"5 + true; 5;",
			"type mismatch: Integer + Bool",
		},
		{
			"-true",
			"unknown operator: -Bool",
		},
		{
			"true + false;",
			"unknown operator: Bool + Bool",
		},
		{
			"5; true + false; 5",
			"unknown operator: Bool + Bool",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: Bool + Bool",
		},
		{
			`if (10 > 1) {
			if (10 > 1) {
			return true + false;
			}
			return 1;
			}
			`, "unknown operator: Bool + Bool",
		},
		{
			"foobar",
			"Undefined variable: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*obj.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.ErrMsg != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.ErrMsg)
		}
	}
}
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*obj.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Args) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Args)
	}
	if fn.Args[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Args[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())

	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
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
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*obj.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}
func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*obj.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
	input = `"Ashish Tiwari" - "ish T"`
	evaluated = testEval(input)
	str, ok = evaluated.(*obj.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Ashiwari" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}
func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "No string in arguments"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*obj.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.ErrMsg != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.ErrMsg)
			}

		}
	}
}
func TestArrIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a=[1,2,3,];a[1]", 2},
		{`let b=["Ashish",2,3,];b[0]`, "Ashish"},
	}
	for i, tt := range tests {
		eval := testEval(tt.input)
		if i == 0 && eval.DataType() != obj.INTEGER_OBJ {
			t.Errorf("Got datatype %T instead of integer", eval.DataType())
			continue
		}
		if i == 1 && eval.DataType() != obj.STRING_OBJ {
			t.Errorf("Got datatype %T instead of string", eval.DataType())
			continue
		}

	}
}
