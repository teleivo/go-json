package parser

import (
	"github.com/teleivo/go-json/ast"
	"github.com/teleivo/go-json/lexer"
	"github.com/teleivo/go-json/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseJSON() *ast.JSON {
	j := &ast.JSON{}

	for p.curToken.Type != token.EOF {
		el := p.parseElement()
		j.Element = el
		p.nextToken()
	}
	return j
}

func (p *Parser) parseElement() ast.Element {
	switch p.curToken.Type {
	case token.STRING:
		return p.parseString()
	default:
		return nil
	}
}

func (p *Parser) parseString() *ast.String {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}
}
