package ast

import "github.com/Revolyssup/monkey/token"

//Everything is a node in AST and has to implement a TokenLiteral method
type Node interface {
	TokenLiteral() string
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

//Identifiers are token which hold some string like x,y,z...
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

/***LET STATEMENT****/
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

//every statement has a method stateNode.
func (ls *LetStatement) stateNode() {}

//every statement is also a node and hence implements token literal method.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
