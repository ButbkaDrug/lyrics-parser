package lyricsparser

import (
	"log/slog"
	"lp/internal/lexer"
	"lp/internal/parser"
	"lp/pkg/lyrics"
	"os"
)

// Receives an path to a file. Returns lyrics object
func Parse(p string) (lyrics.Lyrics, error) {

	f, err := os.ReadFile(p)

	if err != nil {
		return lyrics.Lyrics{}, err
	}

	logfile, err := os.Create("/tmp/lyricsparser.log")
	if err != nil {
		slog.Error("lyricsparser: failed to create log file.", "error", err.Error())
	}

	defer logfile.Close()

	slog.SetDefault(slog.New(
		slog.NewTextHandler(logfile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	))

	l := lexer.NewLexer(string(f))
	parser := parser.NewParser(l)

	return *parser.Parse(), nil
}
