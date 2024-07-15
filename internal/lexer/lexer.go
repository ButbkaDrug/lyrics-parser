package lexer

import (
	"lp/internal/token"
	"slices"
	"unicode"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}

	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) prevChar() byte {
	if l.position < 1 {
		return 0
	}

	return l.input[l.position-1]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]

}

func (l *Lexer) readIdentifyer() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for unicode.IsDigit(rune(l.ch)) || l.ch == '.' || l.ch == '-' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readUntilChar(char byte) string {
	position := l.position

	for l.ch != char && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

// reads from current position until the end of line char '\n'
func (l *Lexer) readUntilEOL() string {
	position := l.position

	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		tok = newToken(token.LSB, l.ch)
	case ']':
		tok = newToken(token.RSB, l.ch)
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
	case '{':
		tok = newToken(token.LCB, l.ch)
	case '}':
		tok = newToken(token.RCB, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '#':
		tok = newToken(token.HASH, l.ch)
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			l.readChar()
			tok.Type = token.COMMENT
			tok.Literal = l.readUntilEOL()
			return tok
		}

		if unicode.IsDigit(rune(l.peekChar())) {
			tok.Type = token.FLOAT
			tok.Literal = l.readNumber()
			return tok
		}

		tok.Type = token.STR
		tok.Literal = l.readIdentifyer()
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:

		if !l.peekNumber() {

			tok.Type = token.STR
			tok.Literal = l.readString()
			return tok
		}

		num := l.readNumber()
		tok.Literal = num

		if isFloat(tok.Literal) {
			tok.Type = token.FLOAT
			return tok

		}

		tok.Type = token.INT
		return tok
	}

	l.readChar()
	return tok
}

func isFloat(s string) bool {

	for i := range s {
		if s[i] == '.' {
			return true
		}
	}

	return false

}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    t,
		Literal: string(ch),
	}
}

func (l *Lexer) peekNumber() bool {
	pos := l.position

	blacklist := []byte{
		' ',
		'\n',
		0,
		'{',
		'}',
		'[',
		']',
		',',
	}

	for !slices.Contains(blacklist, l.input[pos]) {

		if !unicode.IsDigit(rune(l.input[pos])) && l.input[pos] != '.' {
			return false
		}

		pos++

		if pos == len(l.input) {
			return true
		}
	}

	return true

}

func (l *Lexer) readString() string {
	position := l.position

	blacklist := []byte{
		' ',
		'\n',
		0,
		'{',
		'}',
		'[',
		']',
	}

	for !slices.Contains(blacklist, l.ch) {

		l.readChar()
	}

	return l.input[position:l.position]
}

func stringIsNumber(s string) bool {

	if len(s) < 1 {
		return false
	}

	for _, ch := range s {
		if !unicode.IsDigit(ch) && ch != '.' {
			return false
		}
	}

	return true
}

func isLetter(ch byte) bool {
	if unicode.IsDigit(rune(ch)) {
		return false
	}

	return ch != '[' && ch != ']' && ch != 0 && ch != '\n' && ch != ' '
}
