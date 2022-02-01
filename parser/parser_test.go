package parser

import (
	"fmt"
	"testing"

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
	test := []struct {
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

	for _, tt := range test {
		t.Run("ParseValidArray", func(t *testing.T) {
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
	t.Run("ParseInvalidArray", func(t *testing.T) {
		input := `[  "fantastic",]`

		l := lexer.New(input)
		p := New(l)

		p.ParseJSON()

		errors := p.Errors()
		if want := 1; len(errors) != want {
			t.Fatalf("got %d errors but want %d", len(errors), want)
		}
		// TODO adapt error handling. I do not want to assert on the error
		// message
		if want := fmt.Sprintf("expected next token to be one of '%s', got '%s' instead", []token.TokenType{token.TRUE, token.FALSE, token.NULL, token.NUMBER, token.STRING}, token.RBRACKET); errors[0] != want {
			t.Errorf("got %q, want %q", errors[0], want)
		}
	})
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
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
