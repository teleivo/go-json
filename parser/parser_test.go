package parser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/teleivo/go-json/ast"
	"github.com/teleivo/go-json/lexer"
	"github.com/teleivo/go-json/token"
)

func TestString(t *testing.T) {
	input := `"broccoli"`

	l := lexer.New(input)
	p := New(l)

	j := p.ParseJSON()

	checkParserErrors(t, p)

	if j == nil {
		t.Fatal("ParseJSON() returned nil")
	}
	if j.Element == nil {
		t.Fatal("ParseJSON() returned JSON with no element")
	}
	if !testString(t, j.Element, "broccoli") {
		return
	}
}

func testString(t *testing.T, el ast.Element, want string) bool {
	if el.TokenLiteral() != want {
		t.Errorf("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	str, ok := el.(*ast.String)
	if !ok {
		t.Errorf("str not *ast.String. got=%T", el)
		return false
	}
	if str.Value != want {
		t.Errorf("got %q, want %q", str.Value, want)
		return false
	}
	return true
}

func TestBoolean(t *testing.T) {
	test := []struct {
		input string
		want  bool
	}{
		{`true`, true},
		{`false`, false},
	}

	for _, tt := range test {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			j := p.ParseJSON()

			checkParserErrors(t, p)

			if j == nil {
				t.Fatal("ParseJSON() returned nil")
			}
			if j.Element == nil {
				t.Fatal("ParseJSON() returned JSON with no element")
			}
			if !testBoolean(t, j.Element, tt.want) {
				return
			}
		})
	}
}

func testBoolean(t *testing.T, el ast.Element, want bool) bool {
	if want := fmt.Sprintf("%t", want); el.TokenLiteral() != want {
		t.Errorf("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	b, ok := el.(*ast.Boolean)
	if !ok {
		t.Errorf("el not *ast.Boolean. got=%T", el)
		return false
	}
	if b.Value != want {
		t.Errorf("got %t, want %t", b.Value, want)
		return false
	}
	return true
}

func TestNull(t *testing.T) {
	input := `null`

	l := lexer.New(input)
	p := New(l)

	j := p.ParseJSON()

	checkParserErrors(t, p)

	if j == nil {
		t.Fatal("ParseJSON() returned nil")
	}
	if j.Element == nil {
		t.Fatal("ParseJSON() returned JSON with no element")
	}
	if !testNull(t, j.Element) {
		return
	}
}

func testNull(t *testing.T, el ast.Element) bool {
	if want := "null"; el.TokenLiteral() != want {
		t.Errorf("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	_, ok := el.(*ast.Null)
	if !ok {
		t.Errorf("j not *ast.Null. got=%T", el)
		return false
	}
	return true
}

func TestArray(t *testing.T) {
	vt := []struct {
		desc  string
		input string
		ast   []astAssertion
	}{
		{
			desc:  "Empty",
			input: `[  ]`,
		},
		{
			desc:  "Simple",
			input: `[  "fantastic", true, null, "carrot"]`,
			ast: []astAssertion{
				assertString("fantastic"),
				assertBoolean(true),
				assertNull(),
				assertString("carrot"),
			},
		},
	}
	for _, tt := range vt {
		t.Run("ParseValidArray"+tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			j := p.ParseJSON()

			checkParserErrors(t, p)

			if j == nil {
				t.Fatal("ParseJSON() returned nil")
			}
			if j.Element == nil {
				t.Fatal("ParseJSON() returned JSON with no element")
			}
			ar, ok := j.Element.(*ast.Array)
			if !ok {
				t.Fatalf("j.Element not *ast.Array. got=%T", j.Element)
			}

			for i, at := range tt.ast {
				if i >= len(ar.Elements) {
					t.Fatalf("no element is left in array. got %d, want %d elements.", len(ar.Elements), len(tt.ast))
				}
				if !at(t, ar.Elements[i]) {
					return
				}
			}
		})
	}

	ivt := []struct {
		input    string
		actual   string
		expected []token.TokenType
	}{
		{
			input:    `[ `,
			actual:   token.EOF,
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING, token.RBRACKET},
		},
		{
			input:    `[  "fantastic",]`,
			actual:   token.RBRACKET,
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING},
		},
		{
			input:    `[  "fantastic",`,
			actual:   token.EOF,
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING},
		},
		{
			input:    `[  "fantastic"`,
			actual:   token.EOF,
			expected: []token.TokenType{token.COMMA, token.RBRACKET},
		},
	}
	for _, tt := range ivt {
		t.Run("ParseInvalidArray", func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			p.ParseJSON()

			errs := p.Errors()
			if want := 1; len(errs) != want {
				t.Fatalf("ParseJSON(%q): got %d errors but want %d", tt.input, len(errs), want)
			}
			err, ok := errs[0].(*ParseError)
			if !ok {
				t.Fatalf("ParseJSON(%q): err not *ParseError got=%T", tt.input, errs[0])
			}
			if want := tt.actual; string(err.Actual.Type) != want {
				t.Fatalf("ParseJSON(%q): got err.Actual %q, expected %q", tt.input, err.Actual.Type, want)
			}
			opt := cmpopts.SortSlices(func(a, b token.TokenType) bool {
				return a < b
			})
			if diff := cmp.Diff(tt.expected, err.Expected, opt); diff != "" {
				t.Errorf("ParseJSON(%q): err.Expected mismatch (-want, +got): %s\n", tt.input, diff)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	test := []struct {
		desc  string
		input ParseError
		want  string
	}{
		{
			desc: "ExpectsSingleToken",
			input: ParseError{
				Actual:   token.Token{Type: token.LBRACE, Literal: token.LBRACE},
				Expected: []token.TokenType{token.COLON},
			},
			want: "expected token : got { instead",
		},
		{
			desc: "ExpectsMultipleTokens",
			input: ParseError{
				Actual:   token.Token{Type: token.LBRACE, Literal: token.LBRACE},
				Expected: []token.TokenType{token.COLON, token.FALSE},
			},
			want: "expected one of tokens :, FALSE got { instead",
		},
	}
	for _, tt := range test {
		t.Run(tt.desc, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.input.Error()); diff != "" {
				t.Errorf("ParseError() mismatch (-want +got): %s\n", diff)
			}
		})
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}
	t.FailNow()
}

type astAssertion func(t *testing.T, el ast.Element) bool

func assertNull() astAssertion {
	return func(t *testing.T, el ast.Element) bool {
		return testNull(t, el)
	}
}

func assertBoolean(want bool) astAssertion {
	return func(t *testing.T, el ast.Element) bool {
		return testBoolean(t, el, want)
	}
}

func assertString(want string) astAssertion {
	return func(t *testing.T, el ast.Element) bool {
		return testString(t, el, want)
	}
}
