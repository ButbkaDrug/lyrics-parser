package lyrics

type Lyrics struct {
	Name     string
	Sections []Section
}

type Section struct {
	Name     string
	Number   int
	Lang     string
	Comments []string
	Text     []TextLine
}

type TextLine struct {
	Offset float64
	Length float64
	Text   string
}

func NewLyrics() *Lyrics {
	return &Lyrics{}
}
