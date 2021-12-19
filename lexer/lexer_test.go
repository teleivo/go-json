package lexer

import (
	"testing"

	"github.com/teleivo/go-template/token"
)

func TestNextToken(t *testing.T) {
	t.Run("LexCompleteJSON", func(t *testing.T) {
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
				t.Fatalf("character %q - token type wrong. got=%q, want=%q",
					l.ch, tok.Type, tt.expectedType)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("character %q - token literal wrong. got=%q, want=%q",
					l.ch, tok.Literal, tt.expectedLiteral)
			}
		}
	})
	t.Run("LexStrings", func(t *testing.T) {
		tests := []struct {
			input           string
			expectedLiteral string
		}{
			{`"fries"`, `fries`},
			{`"french fries"`, `french fries`},
			{`"french   fries"`, `french   fries`},
			{`"french\nfries"`, `french\nfries`},
			{`"french\tfries\r\n"`, `french\tfries\r\n`},
			{`"french\"fries\"`, `french\"fries\"`},
			{`"\/french\\fries\b"`, `\/french\\fries\b`},
		}

		for _, tt := range tests {
			l := New(tt.input)

			tok := l.NextToken()

			if tok.Type != token.STRING {
				t.Fatalf("input %s - token type wrong. got=%s, want=%s",
					tt.input, tok.Type, token.STRING)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("input %s - token literal wrong. got=%s, want=%s",
					tt.input, tok.Literal, tt.expectedLiteral)
			}
		}
	})
	// TODO add focused tests for number
}
