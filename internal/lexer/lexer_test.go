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
    Drums In
    [Pre Chorus]
    [Pre Chorus R]
    `

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
		{token.NEWLINE, "\n"},
		{token.LSB, "["},
		{token.STR, "Pre"},
		{token.STR, "Chorus"},
		{token.RSB, "]"},
		{token.NEWLINE, "\n"},
		{token.LSB, "["},
		{token.STR, "Pre"},
		{token.STR, "Chorus"},
		{token.STR, "R"},
		{token.RSB, "]"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] wrong tokentype. expected %q got %q\n%v != %v", i, tt.expectedType, tok.Type, tt.expectedLiteral, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] wrong literal. expected %q got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestReadString(t *testing.T) {
	input := `5sus4/7 2m7 4sus2`

	l := NewLexer(input)

	tests := []string{
		"5sus4/7",
		"2m7",
	}

	for i, tt := range tests {
		str := l.readString()

		if str != tt {
			t.Fatalf("test[%d] failed expected string %s got %s\n", i, tt, str)
		}

		l.skipWhitespace()
	}

}

func TestStringIsDigit(t *testing.T) {
	input := []string{"5sus4/7", "123", "", "\n"}

	tests := []bool{false, true, false, false}

	for i, tt := range tests {
		if res := stringIsNumber(input[i]); res != tt {
			t.Fatalf("test[%d] failed expted %v got %v", i, tt, res)
		}
	}
}

func TestChordsLexing(t *testing.T) {
	input := `[ Verse 1 ]
    #1 5sus4/7 2m7 4sus2 1
    `

	l := NewLexer(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LSB, "["},
		{token.STR, "Verse"},
		{token.INT, "1"},
		{token.RSB, "]"},
		{token.NEWLINE, "\n"},
		{token.HASH, "#"},
		{token.INT, "1"},
		{token.STR, "5sus4/7"},
		{token.STR, "2m7"},
		{token.STR, "4sus2"},
		{token.INT, "1"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] wrong tokentype. expected %q got %q\n%s != %s", i, tt.expectedType, tok.Type, tt.expectedLiteral, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] wrong literal. expected %q got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
