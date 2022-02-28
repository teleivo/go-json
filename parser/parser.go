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

func (p *Parser) ParseJSON() (*ast.JSON, error) {
	j := &ast.JSON{}

	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.ILLEGAL) {
		el, err := p.parseElement()
		if err != nil {
			return j, err
		}
		j.Element = el
		p.nextToken()
	}
	if p.curTokenIs(token.ILLEGAL) {
		return j, &ParseError{Actual: p.curToken}
	}
	return j, nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseElement() (ast.Element, error) {
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
		return nil, nil
	}
}

func (p *Parser) parseString() (*ast.String, error) {
	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseBoolean() (*ast.Boolean, error) {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Literal == "true"}, nil
}

func (p *Parser) parseNull() (*ast.Null, error) {
	return &ast.Null{Token: p.curToken}, nil
}

func (p *Parser) parseNumber() (*ast.Number, error) {
	nr := &ast.Number{Token: p.curToken}

	vl, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number: %w", err)
	}
	nr.Value = vl

	return nr, nil
}

func (p *Parser) parseArray() (*ast.Array, error) {
	ar := &ast.Array{Token: p.curToken, Elements: make([]ast.Element, 0)}

	// array should either be closed or contain an element
	if err := p.expectPeek(token.RBRACKET, token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET); err != nil {
		return nil, err
	}
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		el, err := p.parseElement()
		if err != nil {
			return nil, err
		}
		ar.Elements = append(ar.Elements, el)

		if err := p.expectPeek(token.COMMA, token.RBRACKET); err != nil {
			return nil, err
		}
		// if curToken is a comma, then peekToken should be an element
		if p.curTokenIs(token.COMMA) {
			if err := p.expectPeek(token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET); err != nil {
				return nil, err
			}
		}
	}
	return ar, nil
}

func (p *Parser) expectPeek(tt ...token.TokenType) error {
	for _, t := range tt {
		if p.peekTokenIs(t) {
			p.nextToken()
			return nil
		}
	}
	return &ParseError{Expected: tt, Actual: p.peekToken}
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
