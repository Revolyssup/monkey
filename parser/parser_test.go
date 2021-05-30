package parser

import (
	"fmt"
	"testing"

	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/lexer"
)

//This function runs on the premise that there were no errors encountered in checkParserErrors
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letstmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letstmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letstmt.Name.Value)

		return false
	}

	if letstmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letstmt.Name)
		return false
	}
	return true
}

func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let  = 10;
		let foobar = 838383;
		`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		fmt.Println("####", stmt)
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Program has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("[Monke angry: ] %q ", msg)
	}
	t.FailNow()
}
