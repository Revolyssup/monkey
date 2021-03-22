package lexer

import (
	"testing"

	"github.com/Revolyssup/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
	x + y;
	};
	let result = add(five, ten);
	
	
	
	 `
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_BRACKET, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_BRACKET, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_BRACKET, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_BRACKET, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := New(input)

	for i, tt := range tests {
		tok := lex.NextToken()

		if tok.Type != tt.Type {
			t.Fatalf("test[%d]: Wrong Type Token. Expected %q--Got %q", i, tt.Type, tok.Type)

		}
		if tok.Literal != tt.Literal {
			t.Fatalf("test[%d]: Wrong Literal. Expected %q--Got %q", i, tt.Literal, tok.Literal)

		}

	}
}
