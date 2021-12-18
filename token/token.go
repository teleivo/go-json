package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	NUMBER = "NUMBER"

	COMMA = ","
	COLON = ":"
	QUOTE = "\""

	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "["
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
