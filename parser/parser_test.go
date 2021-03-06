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

	j, err := p.ParseJSON()

	checkParserErrors(t, input, err)

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

			j, err := p.ParseJSON()

			checkParserErrors(t, tt.input, err)

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

	j, err := p.ParseJSON()

	checkParserErrors(t, input, err)

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

func TestNumber(t *testing.T) {
	vt := []struct {
		desc  string
		input string
		want  float64
	}{
		{
			desc:  "Integer",
			input: `1501245569`,
			want:  1501245569,
		},
		{
			desc:  "Float",
			input: `2.34`,
			want:  2.34,
		},
		{
			desc:  "Exponent",
			input: `-3.146e7`,
			want:  -3.146e+07,
		},
	}
	for _, tt := range vt {
		t.Run("ParseValidNumber"+tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			j, err := p.ParseJSON()

			checkParserErrors(t, tt.input, err)

			tf := prefixTestPrint(t, tt.input, t.Fatalf)
			te := prefixTestPrint(t, tt.input, t.Errorf)
			if j == nil {
				tf("returned nil")
			}
			if j.Element == nil {
				tf("returned with no element")
			}
			if !testNumber(te, j.Element, tt.want) {
				return
			}
		})
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

			j, err := p.ParseJSON()

			checkParserErrors(t, tt.input, err)

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

			_, err := p.ParseJSON()

			tf := prefixTestPrint(t, tt.input, t.Fatalf)
			te := prefixTestPrint(t, tt.input, t.Errorf)
			if err == nil {
				tf("got no error but want one")
			}
			pe, ok := err.(*ParseError)
			if !ok {
				tf("err not *ParseError got=%T", err)
			}
			if want := tt.actual; string(pe.Actual.Type) != want {
				tf("got err.Actual %q, expected %q", pe.Actual.Type, want)
			}
			opt := cmpopts.SortSlices(func(a, b token.TokenType) bool {
				return a < b
			})
			if diff := cmp.Diff(tt.expected, pe.Expected, opt); diff != "" {
				te("err.Expected mismatch (-want, +got): %s\n", diff)
			}
		})
	}
}

func TestParseIllegal(t *testing.T) {
	input := `2.a34`
	l := lexer.New(input)
	p := New(l)

	_, err := p.ParseJSON()

	if err == nil {
		t.Fatal("expected error but got none")
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

func checkParserErrors(t *testing.T, input string, err error) {
	if err == nil {
		return
	}

	tf := prefixTestPrint(t, input, t.Fatalf)
	tf("parser error: %q", err)
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

func testNumber(te func(format string, args ...interface{}), el ast.Element, want float64) bool {
	// if el.TokenLiteral() != want {
	// 	te("got %q, want %q", el.TokenLiteral(), want)
	// 	return false
	// }
	nr, ok := el.(*ast.Number)
	if !ok {
		te("nr not *ast.Number. got=%T", el)
		return false
	}
	if nr.Value != want {
		te("got %f, want %f", nr.Value, want)
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
