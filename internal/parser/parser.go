package parser

import (
	"fmt"
	"log"
	"log/slog"
	"lp/internal/lexer"
	"lp/internal/token"
	"lp/pkg/chords"
	"lp/pkg/lyrics"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	file, err := os.Create("parser.log")
	if err != nil {
		file = os.Stdout
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if t == p.peekToken.Type {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) expectPeekIs(t token.TokenType) bool {
	if t == p.peekToken.Type {
		return true
	}

	return false
}

func (p *Parser) parseSectionTitle() *lyrics.Section {
	slog.Debug("start parsing title", "token", p.curToken, "peek", p.peekToken)
	s := &lyrics.Section{}

	if p.curToken.Type != token.LSB {
		log.Fatalf("section tilte must start with %v got %v instead", token.LSB, p.curToken.Type)
	}

	//NOTE: Parsing section title
	if !p.expectPeek(token.STR) {
		log.Fatal("section name should start with a word")
	}

	var name string

	slog.Debug("start parsing title name", "token", p.curToken, "peek", p.peekToken)
	for p.curToken.Type == token.STR && len(p.curToken.Literal) > 2 {
		name += " " + p.curToken.Literal
		p.nextToken()
	}

	s.Name = strings.Trim(name, " ")

	if p.curToken.Type == token.RSB {
		p.nextToken()
		return s
	}

	p.skipNewLines()

	//NOTE: Parsing section number if specified

	slog.Debug("start parsing title number", "token", p.curToken, "peek", p.peekToken)
	if p.curToken.Type == token.INT {

		id, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			log.Fatal("failed to parse section number")
		}

		s.Number = id

		p.nextToken()
	}

	if p.curToken.Type == token.RSB {
		p.nextToken()
		return s
	}

	//NOTE: Parsing section Language if specified

	slog.Debug("start parsing title language", "token", p.curToken, "peek", p.peekToken)
	if p.curToken.Type == token.STR && len(p.curToken.Literal) < 3 {
		s.Lang = p.curToken.Literal
		p.nextToken()
	}

	//NOTE: Finishing with section header
	slog.Debug("finishing parsing title", "token", p.curToken, "peek", p.peekToken)
	if p.curToken.Type != token.RSB {
		log.Fatalf("failed to parse section. expected ']' got %v", p.curToken.Literal)
	}
	p.nextToken()

	slog.Debug("returning", "token", p.curToken, "peek", p.peekToken, "return", *s)
	return s
}

func (p *Parser) parseComments() []string {
	comments := []string{}
	for p.curToken.Type == token.COMMENT {
		comments = append(comments, p.curToken.Literal)

		if !p.expectPeek(token.NEWLINE) {
			log.Fatalf("failed to parse comments. expected new line got %q", p.curToken.Type)
		}

		if !p.expectPeekIs(token.COMMENT) {
			break
		}
		p.nextToken()
	}

	slog.Debug("finished with comments", "token", p.curToken, "peek", p.peekToken, "value", comments)

	return comments
}

func (p *Parser) parserChordsLine() []chords.Chord {
	slog.Debug("parsing chords section", "tok", p.curToken, "peek", p.peekToken)
	result := []chords.Chord{}

	if p.curToken.Type != token.INT && p.curToken.Type != token.STR {
		slog.Error("faild to parse chords line expected str or int.", "token", p.peekToken)
		return result
	}

	slog.Debug("true?", "newline", p.peekToken.Type != token.NEWLINE, "EOF", p.peekToken.Type != token.EOF)
	for p.curToken.Type != token.NEWLINE && p.curToken.Type != token.EOF {
		slog.Debug("we are in the looooooophole", "tok", p.curToken)

		result = append(result, chords.NewChord(p.curToken.Literal))
		p.nextToken()

	}

	slog.Debug("finished with a chord line", "tok", p.curToken, "peek", p.peekToken)

	return result

}

func (p *Parser) parseTextLine() string {
	slog.Debug("start parsing text line", "token", p.curToken)

	text := ""
	for p.curToken.Type == token.STR || p.curToken.Type == token.INT {
		text = fmt.Sprintf("%s %s", text, p.curToken.Literal)

		p.nextToken()
	}

	slog.Debug("exiting parse line", "token", p.curToken, "line", text)

	return strings.Trim(text, " ")
}

func (p *Parser) parseWholeLine() string {
	slog.Debug("start parsing whole line", "token", p.curToken)

	text := ""
	for p.curToken.Type != token.NEWLINE && p.curToken.Type != token.EOF {
		text = fmt.Sprintf("%s%s", text, p.curToken.Literal)
		slog.Debug("next token", "tok", p.peekToken)

		p.nextToken()
	}

	slog.Debug("exiting parse line", "token", p.curToken, "line", text)

	return strings.Trim(text, " ")

}

func (p *Parser) parseSection() *lyrics.Section {
	s := &lyrics.Section{}

	slog.Debug("parse section about to hit for loop", "token", p.curToken, "peek", p.peekToken)

	for p.peekToken.Type != token.LSB && p.curToken.Type != token.EOF {
		slog.Debug("iterating through tokens", "token", p.curToken)
		switch p.curToken.Type {
		case token.LSB:
			title := p.parseSectionTitle()
			s.Name = title.Name
			s.Number = title.Number
			s.Lang = title.Lang
		case token.COMMENT:
			slog.Debug("parsing comments", "token", p.curToken, "peek", p.peekToken)
			s.Comments = p.parseComments()
		case token.LCB:
			offset, length := p.parseOffset()
			p.skipNewLines()
			text := p.parseTextLine()
			s.Text = append(s.Text, lyrics.TextLine{
				Offset: offset,
				Length: length,
				Text:   text,
			})
		case token.HASH:
			p.nextToken()
			s.Chords = append(s.Chords, p.parserChordsLine())
			continue
		case token.STR:
			slog.Debug("parsing line", "token", p.curToken, "peek", p.peekToken)
			s.Text = append(s.Text, lyrics.TextLine{
				Offset: 0,
				Length: 0,
				Text:   p.parseTextLine(),
			})
		case token.NEWLINE:
			slog.Debug("parse section: skip new line", "token", p.curToken, "peek", p.peekToken)
			p.nextToken()
		case token.INT:
			slog.Debug("parse section: skip int", "token", p.curToken, "peek", p.peekToken)
			s.Text = append(s.Text, lyrics.TextLine{
				Offset: 0,
				Length: 0,
				Text:   p.parseTextLine(),
			})

		default:
			slog.Debug("parse section: skip undefined case", "token", p.curToken, "peek", p.peekToken)
			p.nextToken()

		}
	}

	slog.Debug("return from parseSection", "token", p.curToken, "peek", p.peekToken, "return", s)
	return s

}

func (p *Parser) skipNewLines() {
	for p.curToken.Type == token.NEWLINE {
		p.nextToken()
	}
}

func (p *Parser) parseOffset() (float64, float64) {
	slog.Debug("start parsing offset", "token", p.curToken, "peek", p.peekToken)
	if !p.expectPeek(token.FLOAT) && !p.expectPeek(token.INT) {
		log.Fatalf("failed to parse offset. expected float or int, got %v", p.curToken.Type)
	}

	slog.Debug("parsing offset", "token", p.curToken, "peek", p.peekToken)
	offset, err := strconv.ParseFloat(p.curToken.Literal, 10)
	if err != nil {
		log.Fatalf("failed to parse offset: %q", err)
	}

	if !p.expectPeek(token.COMMA) {
		log.Fatalf("failed to parse offset. expected %v got %v", token.COMMA, p.curToken.Type)
	}

	slog.Debug("parsing length", "token", p.curToken, "peek", p.peekToken)
	if !p.expectPeek(token.FLOAT) && !p.expectPeek(token.INT) {
		log.Fatalf("failed to parse length. expected %v got %v", token.FLOAT, p.curToken.Type)
	}

	length, err := strconv.ParseFloat(p.curToken.Literal, 10)
	if err != nil {
		log.Fatalf("failed to parse length %q", err)
	}

	if !p.expectPeek(token.RCB) {
		log.Fatalf("failed to parse offset/length. expected %v got %v", token.RCB, p.curToken.Type)
	}
	p.nextToken()

	slog.Debug("exiting parsing length", "token", p.curToken, "peek", p.peekToken, "offset", offset, "length", length)
	return offset, length
}

func (p *Parser) Parse() *lyrics.Lyrics {

	l := lyrics.NewLyrics()

	for p.curToken.Type != token.EOF {
		switch p.curToken.Type {
		case token.LSB:
			s := p.parseSection()
			if s != nil {
				l.Sections = append(l.Sections, *s)
			}
		}

		p.nextToken()
	}

	return l
}
