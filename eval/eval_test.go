package eval

import (
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
	return Eval(prog)
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
