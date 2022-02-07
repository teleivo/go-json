package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/teleivo/go-json/ast"
	"github.com/teleivo/go-json/lexer"
	"github.com/teleivo/go-json/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []error{}}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []error {
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
	p.errors = append(p.errors, &ParseError{Expected: tt, Actual: p.peekToken})
}

func (p *Parser) parseElement() ast.Element {
	switch p.curToken.Type {
	case token.STRING:
		return p.parseString()
	case token.TRUE, token.FALSE:
		return p.parseBoolean()
	case token.NULL:
		return p.parseNull()
	case token.NUMBER:
		return p.parseNumber()
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

func (p *Parser) parseNumber() *ast.Number {
	nr := &ast.Number{Token: p.curToken}

	vl, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("failed to parse number: %w", err))
		return nil
	}
	nr.Value = vl

	return nr
}

func (p *Parser) parseArray() *ast.Array {
	ar := &ast.Array{Token: p.curToken, Elements: make([]ast.Element, 0)}

	// array should either be closed or contain an element
	if !p.expectPeek(token.RBRACKET, token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET) {
		return nil
	}
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		el := p.parseElement()
		ar.Elements = append(ar.Elements, el)

		if !p.expectPeek(token.COMMA, token.RBRACKET) {
			return nil
		}
		// if curToken is a comma, then peekToken should be an element
		if p.curTokenIs(token.COMMA) {
			if !p.expectPeek(token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET) {
				return nil
			}
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

type ParseError struct {
	Expected []token.TokenType
	Actual   token.Token
}

func (pe *ParseError) Error() string {
	var sb strings.Builder
	sb.WriteString("expected")
	if len(pe.Expected) > 1 {
		sb.WriteString(" one of tokens ")
		for i, t := range pe.Expected {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(string(t))
		}
	} else if len(pe.Expected) == 1 {
		sb.WriteString(" token ")
		sb.WriteString(string(pe.Expected[0]))
	}
	sb.WriteString(" got ")
	sb.WriteString(pe.Actual.Literal)
	sb.WriteString(" instead")
	return sb.String()
}
