package lexer

import (
	"testing"

	"github.com/teleivo/go-template/token"
)

func TestNextToken(t *testing.T) {
	input := `{"cookies": 200, "ingredients": ["flour", "salt"]}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.STRING, "cookies"},
		{token.COLON, ":"},
		{token.NUMBER, "200"},
		{token.COMMA, ","},
		{token.STRING, "ingredients"},
		{token.COLON, ":"},
		{token.LBRACKET, "["},
		{token.STRING, "flour"},
		{token.COMMA, ","},
		{token.STRING, "salt"},
		{token.RBRACKET, "]"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for _, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("character %q - tokentype wrong. got=%q, want=%q",
				l.ch, tok.Type, tt.expectedType)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("character %q - tokentype wrong. got=%q, want=%q",
				l.ch, tok.Literal, tt.expectedLiteral)
		}
	}

	// TODO add focused tests for string and number
}
