package lexer

import (
	"testing"

	"github.com/Revolyssup/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){};,`
	tests := []token.Token{
		{token.COMMA, ","},
		{token.EOF, ""},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.LEFT_BRACKET, "("},
		{token.RIGHT_BRACKET, ")"},
		{token.SEMICOLON, ";"},
	}
}
