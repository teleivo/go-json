package parser

import (
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
	if want := "broccoli"; j.Element.TokenLiteral() != want {
		t.Fatalf("got %q, want %q", j.Element.TokenLiteral(), want)
	}
	str, ok := j.Element.(*ast.String)
	if !ok {
		t.Fatalf("j not *ast.String. got=%T", j.Element)
	}
	if want := "broccoli"; str.Value != want {
		t.Fatalf("got %q, want %q", str.Value, want)
	}
}

func TestTrue(t *testing.T) {
	input := `true`

	l := lexer.New(input)
	p := New(l)

	j := p.ParseJSON()

	if j == nil {
		t.Fatal("ParseJSON() returned nil")
	}
	if j.Element == nil {
		t.Fatal("ParseJSON() returned JSON with no element")
	}
	if want := "true"; j.Element.TokenLiteral() != want {
		t.Fatalf("got %q, want %q", j.Element.TokenLiteral(), want)
	}
	str, ok := j.Element.(*ast.Boolean)
	if !ok {
		t.Fatalf("j not *ast.String. got=%T", j.Element)
	}
	if want := true; str.Value != want {
		t.Fatalf("got %t, want %t", str.Value, want)
	}
}
