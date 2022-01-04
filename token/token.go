package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	NULL   = "NULL"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	STRING = "STRING"
	NUMBER = "NUMBER"

	COMMA = ","
	COLON = ":"

	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
