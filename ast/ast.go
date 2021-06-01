package ast

import (
	"bytes"

	"github.com/Revolyssup/monkey/token"
)

//Everything is a node in AST and has to implement a TokenLiteral method
type Node interface {
	TokenLiteral() string
	String() string // Return the exact string of code. Useful for debugging
}

//There are two types of node. Expression and Statement.

type Expression interface {
	Node
	expNode()
}

type Statement interface {
	Node
	stateNode()
}

//Our program is essentially a slice of statements.

//Root node
type Program struct {
	Statements []Statement
}

//Like other nodes, root node also implements a token literal method
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral() //Every further Node(statement/exp) will implement its tokenliteral
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

//Identifiers are token which hold some string like x,y,z...
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

/***LET STATEMENT****/
type LetStatement struct {
	Token token.Token //LET token
	Name  *Identifier
	Value Expression
}

//every statement has a method stateNode.
func (ls *LetStatement) stateNode() {}

//every statement is also a node and hence implements token literal method.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String() + " = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String() + ";")
	}

	return out.String()
}

/*****RETURN STATEMENT*******/

type ReturnStatement struct {
	Token       token.Token //RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) stateNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String() + ";")
	}

	return out.String()
}

/*************Expression Statement*******/

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) stateNode() {}
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer
	out.WriteString(es.TokenLiteral() + " " + es.Expression.TokenLiteral())

	return out.String()
}
