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

func TestBoolean(t *testing.T) {
	test := []struct {
		input       string
		wantLiteral string
		want        bool
	}{
		{`true`, "true", true},
		{`false`, "false", false},
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
			if j.Element.TokenLiteral() != tt.wantLiteral {
				t.Fatalf("got %q, want %q", j.Element.TokenLiteral(), tt.wantLiteral)
			}
			b, ok := j.Element.(*ast.Boolean)
			if !ok {
				t.Fatalf("j not *ast.Boolean. got=%T", j.Element)
			}
			if b.Value != tt.want {
				t.Fatalf("got %t, want %t", b.Value, tt.want)
			}
		})
	}
}
