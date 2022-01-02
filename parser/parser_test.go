package parser

import (
	"fmt"
	"testing"

	"github.com/teleivo/go-json/ast"
	"github.com/teleivo/go-json/lexer"
)

func TestString(t *testing.T) {
	input := `"broccoli"`

	l := lexer.New(input)
	p := New(l)

	j := p.ParseJSON()

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
	}
}
