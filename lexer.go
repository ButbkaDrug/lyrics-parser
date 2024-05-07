package main

import "unicode"

const (
	LSB     = "["
	RSB     = "]"
	IDENT   = "IDENTIFYER"
	NEWLINE = "NEWLINE"
	STR     = "STRING"
	INT     = "INTAGER"
	EOF     = "EOF"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

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

	for unicode.IsDigit(rune(l.ch)) {
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

func (l *Lexer) nextToken() Token {
	var token Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		token.Type = LSB
		token.Literal = string(l.ch)
	case ']':
		token.Type = RSB
		token.Literal = string(l.ch)
	case '\n':
		token.Type = NEWLINE
		token.Literal = string(l.ch)
	case 0:
		token.Type = EOF
		token.Literal = ""
	default:
		if isLetter(l.ch) {
			token.Type = STR
			token.Literal = l.readIdentifyer()
			return token
		} else if unicode.IsDigit(rune(l.ch)) {
			token.Type = INT
			token.Literal = l.readNumber()
			return token
		}
	}

	l.readChar()
	return token
}

func isLetter(ch byte) bool {
	if unicode.IsDigit(rune(ch)) {
		return false
	}

	return ch != '[' && ch != ']' && ch != 0 && ch != '\n' && ch != ' '
}
