package token

// TODO true, false, null

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	STRING = "STRING"
	NUMBER = "NUMBER"

	COMMA = ","
	COLON = ":"

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
