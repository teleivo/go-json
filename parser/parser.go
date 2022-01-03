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

	for !p.curTokenIs(token.EOF) {
		el := p.parseElement()
		j.Element = el
		p.nextToken()
	}
	return j
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) parseElement() ast.Element {
	switch p.curToken.Type {
	case token.STRING:
		return p.parseString()
	case token.TRUE, token.FALSE:
		return p.parseBoolean()
	case token.NULL:
		return p.parseNull()
	case token.LBRACKET:
		return p.parseArray()
	default:
		return nil
	}
}

func (p *Parser) parseString() *ast.String {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() *ast.Boolean {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Literal == "true"}
}

func (p *Parser) parseNull() *ast.Null {
	return &ast.Null{Token: p.curToken}
}

func (p *Parser) parseArray() *ast.Array {
	ar := &ast.Array{Token: p.curToken, Elements: make([]ast.Element, 0)}

	p.nextToken()
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		el := p.parseElement()
		ar.Elements = append(ar.Elements, el)
		// TODO handle errors with missing comma, missing element after comma
		p.nextToken() // skip comma
		p.nextToken() // move onto next element
	}
	return ar
}
