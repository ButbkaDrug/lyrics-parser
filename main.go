package main

import (
	"log"
	"os"
)

func main() {
	file := "/home/butbkadrug/Music/REAPER/Ever Be - Bethel Music - G - 69 BPM/lyrics.txt"

	f, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	lex := NewLexer(string(f))
	p := NewParser(lex)

	l := p.Parse()

	log.Printf("%+v\n", l.Sections)
}
