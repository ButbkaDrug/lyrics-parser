package main

type Lyrics struct {
	Name     string
	Sections []Section
}

type Section struct {
	Name   string
	Number int
	Lang   string
	Text   []string
}

func NewLyrics() *Lyrics {
	return &Lyrics{}
}
