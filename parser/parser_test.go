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

	checkParserErrors(t, input, p)

	tf := prefixTestPrint(t, input, t.Fatalf)
	te := prefixTestPrint(t, input, t.Errorf)
	if j == nil {
		tf("returned nil")
	}
	if j.Element == nil {
		tf("returned with no element")
	}
	if !testString(te, j.Element, "broccoli") {
		return
	}
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

			checkParserErrors(t, tt.input, p)

			tf := prefixTestPrint(t, tt.input, t.Fatalf)
			te := prefixTestPrint(t, tt.input, t.Errorf)
			if j == nil {
				tf("returned nil")
			}
			if j.Element == nil {
				tf("returned with no element")
			}
			if !testBoolean(te, j.Element, tt.want) {
				return
			}
		})
	}
}

func TestNull(t *testing.T) {
	input := `null`

	l := lexer.New(input)
	p := New(l)

	j := p.ParseJSON()

	checkParserErrors(t, input, p)

	tf := prefixTestPrint(t, input, t.Fatalf)
	te := prefixTestPrint(t, input, t.Errorf)
	if j == nil {
		tf("returned nil")
	}
	if j.Element == nil {
		tf("returned with no element")
	}
	if !testNull(te, j.Element) {
		return
	}
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
		{
			desc:  "Nested",
			input: `[  "fantastic", [ true, [ "banana" ], null ], [], "carrot"]`,
			ast: []astAssertion{
				assertString("fantastic"),
				assertArray(
					assertBoolean(true),
					assertArray(
						assertString("banana"),
					),
					assertNull(),
				),
				assertArray(),
				assertString("carrot"),
			},
		},
	}
	for _, tt := range vt {
		t.Run("ParseValidArray"+tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			j := p.ParseJSON()

			checkParserErrors(t, tt.input, p)

			tf := prefixTestPrint(t, tt.input, t.Fatalf)
			te := prefixTestPrint(t, tt.input, t.Errorf)
			if j == nil {
				tf("returned nil")
			}
			if j.Element == nil {
				tf("returned JSON with no element")
			}
			ar, ok := j.Element.(*ast.Array)
			if !ok {
				tf("j.Element not *ast.Array. got=%T", j.Element)
			}

			for i, at := range tt.ast {
				if i >= len(ar.Elements) {
					tf("no element is left in array. got %d, want %d elements.", len(ar.Elements), len(tt.ast))
				}
				if !at(te, ar.Elements[i]) {
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
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING, token.RBRACKET, token.LBRACKET},
		},
		{
			input:    `[  "fantastic",]`,
			actual:   token.RBRACKET,
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET},
		},
		{
			input:    `[  "fantastic",`,
			actual:   token.EOF,
			expected: []token.TokenType{token.FALSE, token.TRUE, token.NULL, token.NUMBER, token.STRING, token.LBRACKET},
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
			f := prefixTestPrint(t, tt.input, t.Fatalf)
			e := prefixTestPrint(t, tt.input, t.Errorf)
			if want := 1; len(errs) != want {
				f("got %d errors but want %d", len(errs), want)
			}
			err, ok := errs[0].(*ParseError)
			if !ok {
				f("err not *ParseError got=%T", errs[0])
			}
			if want := tt.actual; string(err.Actual.Type) != want {
				f("got err.Actual %q, expected %q", err.Actual.Type, want)
			}
			opt := cmpopts.SortSlices(func(a, b token.TokenType) bool {
				return a < b
			})
			if diff := cmp.Diff(tt.expected, err.Expected, opt); diff != "" {
				e("err.Expected mismatch (-want, +got): %s\n", diff)
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
				t.Errorf("ParseError.String() mismatch (-want +got): %s\n", diff)
			}
		})
	}
}

func prefixTestPrint(t *testing.T, input string, prn func(format string, args ...interface{})) func(format string, args ...interface{}) {
	pf := fmt.Sprintf("ParseJSON(%q): ", input)
	return func(format string, args ...interface{}) {
		prn(pf + fmt.Sprintf(format, args...))
	}
}

func checkParserErrors(t *testing.T, input string, p *Parser) {
	te := prefixTestPrint(t, input, t.Errorf)
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	te("parser has %d errors.", len(errors))
	for _, err := range errors {
		te("parser error: %q", err)
	}
	t.FailNow()
}

type astAssertion func(te func(format string, args ...interface{}), el ast.Element) bool

func assertNull() astAssertion {
	return func(te func(format string, args ...interface{}), el ast.Element) bool {
		return testNull(te, el)
	}
}

func assertBoolean(want bool) astAssertion {
	return func(te func(format string, args ...interface{}), el ast.Element) bool {
		return testBoolean(te, el, want)
	}
}

func assertString(want string) astAssertion {
	return func(te func(format string, args ...interface{}), el ast.Element) bool {
		return testString(te, el, want)
	}
}

// TODO do I want to pass in an ast? or a []astAssertion
func assertArray(want ...astAssertion) astAssertion {
	return func(te func(format string, args ...interface{}), el ast.Element) bool {
		return testArray(te, el, want)
	}
}

func testNull(te func(format string, args ...interface{}), el ast.Element) bool {
	if want := "null"; el.TokenLiteral() != want {
		te("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	_, ok := el.(*ast.Null)
	if !ok {
		te("j not *ast.Null. got=%T", el)
		return false
	}
	return true
}

func testBoolean(te func(format string, args ...interface{}), el ast.Element, want bool) bool {
	if want := fmt.Sprintf("%t", want); el.TokenLiteral() != want {
		te("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	b, ok := el.(*ast.Boolean)
	if !ok {
		te("el not *ast.Boolean. got=%T", el)
		return false
	}
	if b.Value != want {
		te("got %t, want %t", b.Value, want)
		return false
	}
	return true
}

func testString(te func(format string, args ...interface{}), el ast.Element, want string) bool {
	if el.TokenLiteral() != want {
		te("got %q, want %q", el.TokenLiteral(), want)
		return false
	}
	str, ok := el.(*ast.String)
	if !ok {
		te("str not *ast.String. got=%T", el)
		return false
	}
	if str.Value != want {
		te("got %q, want %q", str.Value, want)
		return false
	}
	return true
}

func testArray(te func(format string, args ...interface{}), el ast.Element, want []astAssertion) bool {
	// TODO should I also test the TokenLiteral() here?
	ar, ok := el.(*ast.Array)
	if !ok {
		te("j.Element not *ast.Array. got=%T", el)
		return false
	}

	for i, at := range want {
		if i >= len(ar.Elements) {
			te("no element is left in array. got %d, want %d elements.", len(ar.Elements), len(want))
			return false
		}
		if !at(te, ar.Elements[i]) {
			return false
		}
	}
	return true
}
