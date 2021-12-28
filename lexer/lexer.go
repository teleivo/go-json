package lexer

import (
	"errors"

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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // ASCII code for "NUL" character
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
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0: // NUL byte
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
		if isTrue(l.ch) {
			lit, err := l.readTrue()
			tok.Literal = lit
			if err != nil {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = token.TRUE
			}
			return tok
		}
		if isFalse(l.ch) {
			lit, err := l.readFalse()
			tok.Literal = lit
			if err != nil {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = token.FALSE
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

func (l *Lexer) readString() string {
	// TODO read unicode \u1234
	l.readChar() // do not include the outer quotes in the string value
	pos := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() == '"' {
			// move two characters
			l.readChar()
			l.readChar()
		} else {
			l.readChar()
		}
	}
	return l.input[pos:l.position]
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

func (l *Lexer) readTrue() (string, error) {
	// TODO handle not being followed by COMMA
	pos := l.position
	for l.ch != 0 && l.ch != ',' && !isWhitespace(l.ch) {
		l.readChar()
	}
	t := l.input[pos:l.position]
	if t != "true" {
		return t, errors.New("invalid token true: expects 'true'")
	}
	return t, nil
}

func (l *Lexer) readFalse() (string, error) {
	pos := l.position
	for l.ch != 0 && l.ch != ',' && l.ch != '}' && !isWhitespace(l.ch) {
		l.readChar()
	}
	t := l.input[pos:l.position]
	if t != "false" {
		return t, errors.New("invalid token false: expects 'false'")
	}
	return t, nil
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

func isTrue(ch byte) bool {
	return ch == 't'
}

func isFalse(ch byte) bool {
	return ch == 'f'
}
