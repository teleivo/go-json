package lexer

import (
	"github.com/teleivo/go-template/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (current char)
	readPosition int  // current reading position (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0 // ASCII code for "NUL" character
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0: // NUL byte
		tok.Literal = ""
		tok.Type = token.EOF
	case ' ', '\t', '\n', '\b', '\r', '\f': // eat up whitespace outside of strings
		l.readChar()
		return l.NextToken()
	default:
		if isNumber(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	// TODO read unicode \u1234
	l.readChar() // do not include the outer quotes in the string value
	pos := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peek() == '"' {
			// move two characters
			l.readChar()
			l.readChar()
		} else {
			l.readChar()
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func isNumber(ch byte) bool {
	// TODO its more complicated than that. find out all the valid numbers in
	// JSON. https://www.json.org/json-en.html
	return '0' <= ch && ch <= '9'
}
