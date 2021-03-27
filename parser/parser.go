package parser

import (
	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.NextToken()
	p.NextToken()
	return p
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil { //we will get some sort of parsed statement
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		{
			return p.parseLetStatement()
		}
	default:
		{
			return nil
		}
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letstmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	letstmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for p.currToken.Type != token.SEMICOLON {
		p.NextToken()
	}

	return letstmt
}

//utilities
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	return false
}

func (p *Parser) expectCurr(t token.TokenType) bool {
	if p.currToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	return false
}
