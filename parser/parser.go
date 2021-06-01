package parser

import (
	"fmt"
	"strconv"

	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/token"
)

//To parse expressions using pratt parsing. Based on what type of expression that is, we will define different functins. Broadly they will be in two categories:
type (
	infixParsefunc  func(ast.Expression) ast.Expression //It takes in the expression before the operator/token
	prefixParsefunc func() ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
	//Each token type will have some parse function associated with it.
	infixParsefuncns  map[token.TokenType]infixParsefunc
	prefixParsefuncns map[token.TokenType]prefixParsefunc
}

// These are the precedence of operators which would be passed in function call to specific parseExpression functions.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // ><
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X and !X
	CALL        // func(x)
)

//Parsing expressions
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParsefuncns[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.NextToken()
	p.NextToken()
	p.prefixParsefuncns = make(map[token.TokenType]prefixParsefunc)
	//registering parseExpressinoFunctions

	p.registerPrefixParse(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixParse(token.INTEGER, p.parseIntegerLiteral)
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
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("Expected token type %s. Got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		{
			return p.parseLetStatement()
		}
	case token.RETURN:
		{
			return p.parseReturnStatement()
		}
	default:
		{
			return p.parseExpressionStatement()
		}
	}
}

//parsing different types of statements.

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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retstmt := &ast.ReturnStatement{Token: p.currToken}

	for p.currToken.Type != token.SEMICOLON {
		p.NextToken()
	}
	return retstmt
}

//Parsing expressionns using pratt parser technique.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.expectPeek(token.SEMICOLON) { //Semicolon is not mandatory
		p.NextToken()
	}
	return stmt
}

//utilities
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	p.peekErrors(t)
	return false
}

func (p *Parser) expectCurr(t token.TokenType) bool {
	if p.currToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	return false
}

func (p *Parser) registerPrefixParse(t token.TokenType, f prefixParsefunc) {
	p.prefixParsefuncns[t] = f
}

func (p *Parser) registerInfixParse(t token.TokenType, f infixParsefunc) {
	p.infixParsefuncns[t] = f
}

// different types of parseExpressionfunc based on token type

func (p *Parser) parseIdentifier() ast.Expression { //For token.IDENT
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intexp := &ast.IntegerLiteral{Token: p.currToken}
	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as int64", intexp)
		p.errors = append(p.errors, msg)
		return nil
	}
	intexp.Value = val
	return intexp
}
