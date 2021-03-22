package lexer

import (
	token "github.com/Revolyssup/monkey/token"
)

type Lexer struct {
	input    string
	lastRead int
	readPos  int
	ch       byte
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LEFT_BRACKET, l.ch)
	case ')':
		tok = newToken(token.RIGHT_BRACKET, l.ch)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.ch)
	case 0:
		tok = newToken(token.EOF, l.ch)
	default: //handling identifiers
		if l.isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IdentOrKeyword(tok.Literal)
			return tok
		} else if l.isNumber(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INTEGER
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.read()
	return tok
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.read()
	return l
}

//utilities

func (l *Lexer) readIdentifier() string {
	pos := l.lastRead
	for l.isLetter(l.ch) {
		l.read()
	}
	return l.input[pos:l.lastRead]
}

func (l *Lexer) readNumber() string {
	pos := l.lastRead
	for l.isNumber(l.ch) {
		l.read()
	}
	return l.input[pos:l.lastRead]
}

func (l *Lexer) read() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.lastRead = l.readPos
	l.readPos += 1
}

func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}

//currently only supporiting ASCII
func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
func (l *Lexer) isNumber(ch byte) bool {
	return 0 <= ch && ch <= 9
}
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.read()
	}
}
