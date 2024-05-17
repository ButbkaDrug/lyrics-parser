package lyricsparser

import (
	"lp/internal/lexer"
	"lp/internal/lyrics"
	"lp/internal/parser"
)

// Receives path to a file. Returns lyrics object
func Parse(path string) lyrics.Lyrics {
	l := lexer.NewLexer(path)
	p := parser.NewParser(l)

	return *p.Parse()
}
