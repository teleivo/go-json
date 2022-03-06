package lexer

import (
	"errors"
	"fmt"

	"github.com/teleivo/go-json/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (current char)
	readPosition int  // current reading position (after current char)
	ch           byte // current char under examination
}

var charToKeyword = map[byte]string{
	't': "true",
	'f': "false",
	'n': "null",
}

var keywordToToken = map[string]token.TokenType{
	"true":  token.TRUE,
	"false": token.FALSE,
	"null":  token.NULL,
}

const NUL = 0 // ASCII code for "NUL" character

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = NUL
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return NUL
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

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
		lit, err := l.readString()
		tok.Literal = lit
		if err != nil {
			tok.Type = token.ILLEGAL
		} else {
			tok.Type = token.STRING
		}
	case NUL:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isNumber(l.ch) {
			lit, err := l.readNumber()
			tok.Literal = lit
			if err != nil {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = token.NUMBER
			}
			return tok
		}
		if isKeyword(l.ch) {
			lit, err := l.readKeyword()
			tok.Literal = lit
			if err != nil {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = keywordToToken[lit]
			}
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func isWhitespace(ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\b' || ch == '\r' || ch == '\f' {
		return true
	}
	return false
}

func (l *Lexer) readString() (string, error) {
	// TODO read unicode \u1234
	var closing bool
	l.readChar() // do not include the outer quotes in the string value
	pos := l.position
	for l.ch != '"' && l.ch != NUL {
		if l.ch == '\\' && l.peekChar() == '"' {
			closing = true
			// move two characters
			l.readChar()
			l.readChar()
		} else {
			l.readChar()
		}
	}
	if !closing && l.ch != '"' {
		return string(l.ch), errors.New("missing closing quotes \"")
	}
	return l.input[pos:l.position], nil
}

func (l *Lexer) readNumber() (string, error) {
	pos := l.position
	for isNumber(l.ch) || l.ch == '.' || l.ch == '-' || l.ch == '+' || l.ch == 'e' || l.ch == 'E' {
		if l.peekChar() == '.' && !isDigit(l.ch) {
			return string(l.ch), errors.New("invalid number token: '.' needs to be preceded by a digit")
		}
		if l.ch == '.' && !isDigit(l.peekChar()) {
			return string(l.peekChar()), errors.New("invalid number token: '.' needs to be followed by a digit")
		}
		if (l.peekChar() == '+' || l.peekChar() == '-') && (l.ch != 'e' && l.ch != 'E') {
			return string(l.peekChar()), errors.New("invalid number token: '+' or '-' needs to be preceded by 'e' or 'E' for exponent")
		}
		l.readChar()
	}
	return l.input[pos:l.position], nil
}

func (l *Lexer) readKeyword() (string, error) {
	k := charToKeyword[l.ch]
	pos := l.position
	for l.ch != NUL && l.position-pos < len(k) {
		l.readChar()
		if l.input[pos:l.position] != k[0:l.position-pos] {
			return string(l.ch), fmt.Errorf("invalid token %q: expect %q", k, k)
		}
	}
	if l.input[pos:l.position] != k {
		return l.input[pos:l.position], fmt.Errorf("invalid token %q: expect %q", k, k)
	}
	return l.input[pos:l.position], nil
}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func isNumber(ch byte) bool {
	return isDigit(ch) || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isKeyword(ch byte) bool {
	_, ok := charToKeyword[ch]
	return ok
}
