package lexer

import (
	"lp/internal/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `[ Verse 1 R ]
    --Drums In
    {0, 3.333}
    Line of text 1
    Line of text 2
    Drums In`

	l := NewLexer(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LSB, "["},
		{token.STR, "Verse"},
		{token.INT, "1"},
		{token.STR, "R"},
		{token.RSB, "]"},
		{token.NEWLINE, "\n"},
		{token.COMMENT, "Drums In"},
		{token.NEWLINE, "\n"},
		{token.LCB, "{"},
		{token.INT, "0"},
		{token.COMMA, ","},
		{token.FLOAT, "3.333"},
		{token.RCB, "}"},
		{token.NEWLINE, "\n"},
		{token.STR, "Line"},
		{token.STR, "of"},
		{token.STR, "text"},
		{token.INT, "1"},
		{token.NEWLINE, "\n"},
		{token.STR, "Line"},
		{token.STR, "of"},
		{token.STR, "text"},
		{token.INT, "2"},
		{token.NEWLINE, "\n"},
		{token.STR, "Drums"},
		{token.STR, "In"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] wrong tokentype. expected %q got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] wrong literal. expected %q got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
