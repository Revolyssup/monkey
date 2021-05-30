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
		let y = 10;
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
		fmt.Println(i+1, stmt.TokenLiteral(), "statement")
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return ass();
	return 1;
	return 2;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	for _, stmt := range program.Statements {
		rstmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not ast.Returnstatement, got=%s", stmt)
			continue
		}
		if rstmt.TokenLiteral() != "return" {
			t.Errorf("statement not ast.Returnstatement, got=%q", rstmt.TokenLiteral())
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
