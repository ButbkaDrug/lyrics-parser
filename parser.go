package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.nextToken()
}

func (p *Parser) parseSection() Section {

	s := Section{}

	p.nextToken()

	if p.curToken.Type != STR {
		log.Fatal("section name should start with a word")
	}

	s.Name = p.curToken.Literal
	p.nextToken()

	for p.curToken.Type == STR && len(p.curToken.Literal) > 2 {
		s.Name = fmt.Sprintf("%s %s", s.Name, p.curToken.Literal)
		p.nextToken()
	}

	if p.curToken.Type == INT {

		id, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			log.Fatal("failed to parse section number")
		}

		s.Number = id

		p.nextToken()
	}

	if p.curToken.Type == STR && len(p.curToken.Literal) == 1 {
		s.Lang = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type != RSB {
		log.Fatal("failed to parse section. expected ']' got ", p.curToken.Literal)
	}
	p.nextToken()

	if p.curToken.Type == NEWLINE {
		p.nextToken()
	}

	var line string
	for p.peekToken.Type != LSB && p.peekToken.Type != EOF {

		if len(line) > 0 && p.curToken.Type == NEWLINE {
			s.Text = append(s.Text, strings.Trim(line, " "))
			p.nextToken()
			line = ""
			continue
		}

		line = fmt.Sprintf("%s %s", line, p.curToken.Literal)

		p.nextToken()
	}

	return s
}

func (p *Parser) Parse() *Lyrics {

	lyrics := NewLyrics()

	for p.curToken.Type != EOF {
		s := Section{}

		if p.curToken.Type == LSB {
			s = p.parseSection()
		}

		lyrics.Sections = append(lyrics.Sections, s)
		p.nextToken()
	}

	return lyrics
}
