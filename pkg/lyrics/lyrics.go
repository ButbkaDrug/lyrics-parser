package lyrics

import (
	"fmt"
	"lp/pkg/chords"
)

type Lyrics struct {
	Name     string
	Sections []Section
}

func (l Lyrics) String() string {
	var text string
	for _, s := range l.Sections {

		text = fmt.Sprintf("%s\n\n%s", text, s)

	}

	return text
}

func (l *Lyrics) Empty() bool {
	return len(l.Sections) > 0
}

type Section struct {
	ID       int
	Name     string
	Number   int
	Lang     string
	Comments []string
	Chords   chords.Chords
	Text     []TextLine
}

func (s *Section) Empty() bool {
	return len(s.Comments) == len(s.Chords) && len(s.Chords) == len(s.Text) && len(s.Text) == 0 && s.Number == 0 && s.Lang == s.Name && s.Name == ""
}

func (s *Section) GetFullName() string {
	name := s.Name

	if s.Number > 0 {
		name = fmt.Sprintf("%s %d", name, s.Number)
	}

	if s.Lang != "" {
		name = fmt.Sprintf("%s %s", name, s.Lang)
	}

	return name
}

func (s Section) String() string {
	text := fmt.Sprintf("[ %s ]", s.GetFullName())

	for _, comment := range s.Comments {
		text = fmt.Sprintf("%s\n--%s", text, comment)

	}

	for i, line := range s.Text {
		if len(s.Chords) > i {
			text += fmt.Sprintf("\n#%s", s.Chords[i])
		}
		text = fmt.Sprintf("%s\n%s", text, line)
	}

	if len(s.Text) == 0 && len(s.Chords) > 0 {
		for _, line := range s.Chords {
			text += fmt.Sprintf("\n#%s", line)
		}
	}

	text = fmt.Sprintf("%s", text)

	return text
}

type TextLine struct {
	Offset float64
	Length float64
	Text   string
}

func (t TextLine) String() string {
	if t.Offset != 0 || t.Length != 0 {
		return fmt.Sprintf("{%v, %v}\n%s", t.Offset, t.Length, t.Text)
	}

	return fmt.Sprintf("%s", t.Text)
}

func NewLyrics() *Lyrics {
	return &Lyrics{}
}
