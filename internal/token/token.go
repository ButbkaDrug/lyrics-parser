package token

const (
	LSB   = "["
	RSB   = "]"
	LCB   = "{"
	RCB   = "}"
	COMMA = ","

	COMMENT = "COMMENT"
	IDENT   = "IDENTIFYER"
	NEWLINE = "NEWLINE"

	STR   = "STRING"
	INT   = "INTAGER"
	FLOAT = "FLOAT"

	NOTES = "NOTES"

	EOF = "EOF"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
