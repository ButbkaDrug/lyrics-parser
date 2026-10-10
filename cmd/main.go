package main

import (
	"log"
	"log/slog"
	"lp/internal/lexer"
	"lp/internal/parser"
	"os"
)

func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	file := "/home/butbkadrug/Music/REAPER/Our God - Chris Tomlin - B - 105 BPM/lyrics.txt"

	f, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(string(f))
	p := parser.NewParser(l)

	lyrics := p.Parse()

	log.Println(lyrics)

}
