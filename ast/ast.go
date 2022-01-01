package ast

import "github.com/teleivo/go-json/token"

type Node interface {
	TokenLiteral() string
}

type Element interface {
	Node
	elementNode()
}

type JSON struct {
	Element Element
}

func (j *JSON) TokenLiteral() string {
	if j.Element != nil {
		return j.Element.TokenLiteral()
	}
	return ""
}

type String struct {
	Token token.Token // the token.STRING token
	Value string
}

func (s *String) elementNode() {}

func (s *String) TokenLiteral() string {
	return s.Token.Literal
}

type Boolean struct {
	Token token.Token // the token.TRUE or token.FALSE
	Value bool
}

func (b *Boolean) elementNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

type Null struct {
	Token token.Token // the token.NULL token
}

func (n *Null) elementNode() {}

func (n *Null) TokenLiteral() string {
	return n.Token.Literal
}
