package lexer

import (
	"testing"

	"github.com/teleivo/go-json/token"
)

func TestLexCompleteJSON(t *testing.T) {
	input := `{"cookies": 200, "ingredients": ["flour", "salt"], "fresh": true, "tasty": false}`

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
		{token.COMMA, ","},
		{token.STRING, "fresh"},
		{token.COLON, ":"},
		{token.TRUE, "true"},
		{token.COMMA, ","},
		{token.STRING, "tasty"},
		{token.COLON, ":"},
		{token.FALSE, "false"},
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
}

func TestLexStrings(t *testing.T) {
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
}

func TestLexNumbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{"200", "200"},
		{"200.3", "200.3"},
		{"0.31", "0.31"},
		{"-0.31", "-0.31"},
		{"-200.3", "-200.3"},
		{"-200.3", "-200.3"},
		{"-200.3e1", "-200.3e1"},
		{"-200.3e+1", "-200.3e+1"},
		{"-200.3e-1", "-200.3e-1"},
		{"-200.3E1", "-200.3E1"},
		{"-200.3E+12", "-200.3E+12"},
		{"-200.3E-12", "-200.3E-12"},
		{"0.31e100", "0.31e100"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != token.NUMBER {
			t.Fatalf("input %q - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, token.STRING)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %q - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}

func TestLexInvalidNumbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		description     string
	}{
		{"a200", "a", "number should only contain digits"},
		{"_200", "_", "number should only contain digits"},
		{"+200", "+", "number cannot be prefixed with +"},
		{"2-00", "-", "- sign is only allowed in first position or after exponent"},
		{"2+00", "+", "+ sign is only allowed after exponent"},
		{"e200", "e", "exponent needs to be preceded by digit"},
		{"E200", "E", "exponent needs to be preceded by digit"},
		{".200", ".", "fraction needs to be preceded by at least one digit"},
		{"-.200", "-", "fraction needs to be preceded by at least one digit"},
		{"+.200", "+", "fraction needs to be preceded by at least one digit"},
		{"1.", string(byte(0)), "fraction needs to be followed by at least one digit"},
		{"2.a34", "a", "fraction needs to be followed by at least one digit"},
		{"0.e+100", "e", "fraction needs to be followed by at least one digit"},
		{"0.e-100", "e", "fraction needs to be followed by at least one digit"},
		{"0.e100", "e", "fraction needs to be followed by at least one digit"},
		{"0.E+100", "E", "fraction needs to be followed by at least one digit"},
		{"0.E-100", "E", "fraction needs to be followed by at least one digit"},
		{"0.E100", "E", "fraction needs to be followed by at least one digit"},
		{"0.200+e100", "+", "exponent cannot be preceded by sign"},
		{"0.200-e100", "-", "exponent cannot be preceded by sign"},
		{"0.200+E100", "+", "exponent cannot be preceded by sign"},
		{"0.200-E100", "-", "exponent cannot be preceded by sign"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != token.ILLEGAL {
			t.Fatalf("input %q - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, token.ILLEGAL)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %q - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}

func TestLexTrueAndFalse(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedToken   token.TokenType
		description     string
	}{
		{`true,`, `true`, token.TRUE, "should lex true if followed by comma"},
		{`true;`, `true`, token.TRUE, "should lex true if followed by semicolon"},
		{`true'`, `true`, token.TRUE, "should lex true if followed by single-quote"},
		{`true"`, `true`, token.TRUE, "should lex true if followed by double-quote"},
		{`true{`, `true`, token.TRUE, "should lex true if followed by left brace"},
		{`true[`, `true`, token.TRUE, "should lex true if followed by left bracket"},
		{`true}`, `true`, token.TRUE, "should lex true if followed by right brace"},
		{`true]`, `true`, token.TRUE, "should lex true if followed by right bracket"},
		{` true	 `, `true`, token.TRUE, "should ignore spaces and tabs"},
		{` true
				`, `true`, token.TRUE, "should ignore newlines"},
		{` false	 `, `false`, token.FALSE, "should ignore spaces and tabs"},
		{` false
				`, `false`, token.FALSE, "should ignore newlines"},
		{`false,`, `false`, token.FALSE, "should lex false if followed by comma"},
		{`false;`, `false`, token.FALSE, "should lex false if followed by semicolon"},
		{`false'`, `false`, token.FALSE, "should lex false if followed by single-quote"},
		{`false"`, `false`, token.FALSE, "should lex false if followed by double-quote"},
		{`false{`, `false`, token.FALSE, "should lex false if followed by left brace"},
		{`false[`, `false`, token.FALSE, "should lex false if followed by left bracket"},
		{`false}`, `false`, token.FALSE, "should lex false if followed by right brace"},
		{`false]`, `false`, token.FALSE, "should lex false if followed by right bracket"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != tt.expectedToken {
			t.Fatalf("input %s - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, tt.expectedToken)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %s - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}

func TestLexInvalidTrueAndFalse(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedToken   token.TokenType
		description     string
	}{
		{`TRUE`, `T`, token.ILLEGAL, "uppercase true is invalid"},
		{`tru`, `tru`, token.ILLEGAL, "incomplete true"},
		{`tru,e`, `e`, token.ILLEGAL, "true interrupted by invalid character"},
		{`t`, `t`, token.ILLEGAL, "incomplete true"},
		{`FALSE`, `F`, token.ILLEGAL, "uppercase false is invalid"},
		{`f`, `f`, token.ILLEGAL, "incomplete false"},
		{`fals`, `fals`, token.ILLEGAL, "incomplete false"},
		{`fals,e`, `e`, token.ILLEGAL, "false interrupted by invalid character"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != tt.expectedToken {
			t.Fatalf("input %s - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, tt.expectedToken)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %s - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}

func TestLexNull(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedToken   token.TokenType
		description     string
	}{
		{`null,`, `null`, token.NULL, "should lex null if followed by comma"},
		{`null;`, `null`, token.NULL, "should lex null if followed by semicolon"},
		{`null'`, `null`, token.NULL, "should lex null if followed by single-quote"},
		{`null"`, `null`, token.NULL, "should lex null if followed by double-quote"},
		{`null[`, `null`, token.NULL, "should lex null if followed by right bracket"},
		{`null{`, `null`, token.NULL, "should lex null if followed by right brace"},
		{`null]`, `null`, token.NULL, "should lex null if followed by right bracket"},
		{`null}`, `null`, token.NULL, "should lex null if followed by right brace"},
		{` null	 `, `null`, token.NULL, "should ignore spaces and tabs"},
		{` null
				`, `null`, token.NULL, "should ignore newlines"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != tt.expectedToken {
			t.Fatalf("input %s - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, tt.expectedToken)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %s - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}

func TestLexInvalidNull(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedToken   token.TokenType
		description     string
	}{
		{`NULL`, `N`, token.ILLEGAL, "uppercase null is invalid"},
		{`nul`, `nul`, token.ILLEGAL, "incomplete null"},
		{`nul,l`, `l`, token.ILLEGAL, "null interrupted by invalid character"},
		{`n`, `n`, token.ILLEGAL, "incomplete null"},
	}

	for _, tt := range tests {
		l := New(tt.input)

		tok := l.NextToken()

		if tok.Type != tt.expectedToken {
			t.Fatalf("input %s - token type wrong. got=%s, want=%s",
				tt.input, tok.Type, tt.expectedToken)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("input %s - token literal wrong. got=%s, want=%s",
				tt.input, tok.Literal, tt.expectedLiteral)
		}
	}
}
