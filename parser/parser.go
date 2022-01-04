package parser

import (
	"fmt"

	"github.com/teleivo/go-json/ast"
	"github.com/teleivo/go-json/lexer"
	"github.com/teleivo/go-json/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
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

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(tt ...token.TokenType) {
	var msg string
	if len(tt) == 1 {
		msg = fmt.Sprintf("expected next token to be '%s', got '%s' instead", tt[0], p.peekToken.Type)
	} else {
		msg = fmt.Sprintf("expected next token to be one of '%s', got '%s' instead", tt, p.peekToken.Type)
	}
	p.errors = append(p.errors, msg)
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

	// TODO peekToken should either be RBRACKET or an element
	p.nextToken()
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		el := p.parseElement()
		ar.Elements = append(ar.Elements, el)

		if !p.expectPeek(token.COMMA, token.RBRACKET) {
			return nil
		}
		// if curToken is a comma, then peekToken should be an element
		if p.curTokenIs(token.COMMA) {
			// TODO an array inside an array is also allowed
			if !p.expectPeek(token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING) {
				return nil
			}
		} else {
			p.nextToken()
		}
	}
	return ar
}

func (p *Parser) expectPeek(tt ...token.TokenType) bool {
	for _, t := range tt {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}
	p.peekError(tt...)
	return false
}
