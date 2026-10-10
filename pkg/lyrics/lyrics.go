package lyrics

import (
	"fmt"
	"io"
	"lp/pkg/chords"
	"lp/pkg/comments"
	"strings"
)

type Offset struct {
	Offset float64
	Length float64
}

type OffsetMap []Offset

func NewOffsetMap(s Section) OffsetMap {
	offsets := make([]Offset, len(s.Text))

	for i, line := range s.Text {
		offsets[i].Offset = line.Offset
		offsets[i].Length = line.Length
	}

	return offsets
}

func (o OffsetMap) Empty() bool {

	for _, offset := range o {
		if offset.Offset != offset.Length {
			return false
		}
	}

	return true
}

// struct that represents the song lyrics/notes/chords
type Lyrics struct {
	//path to a source file
	Path string
	//never found use of it yet
	Name string
	//content of the file in go format
	Sections []Section
}

// writes lyrics structure in text format to a given writer
func (l Lyrics) Write(w io.Writer) (int, error) {

	return w.Write([]byte(l.String()))

}

// return a chords for a section with sepcified id. if not found return empty struct
func (l Lyrics) GetChordsWithID(ID int) chords.Chords {
	for _, section := range l.Sections {
		if section.Number == ID {
			return section.Chords
		}
	}

	return chords.Chords{}
}

// Searches for the chords in the other section with of a same type, going from the closesed to futhers in simillarity.
func (l Lyrics) GetChordsForSection(s Section) chords.Chords {
	chords := s.Chords

	if !chords.Empty() {
		return chords
	}

	alt := l.GetAlternativeSection(s)
	chords = alt.Chords
	if !chords.Empty() {
		return chords
	}

	rel := l.GetRelativeSection(s)
	chords = rel.Chords
	if !chords.Empty() {
		return chords
	}

	relalt := l.GetRelativeAlternativeSection(s)
	chords = relalt.Chords
	if !chords.Empty() {
		return chords
	}

	all := l.GetAllRelativeSections(s)

	for _, relative := range all {
		chords = relative.Chords
		if !chords.Empty() {
			return chords
		}
	}

	return chords
}

// loades a note where note number matches ID parametr. if not found return empty slice. function ment to be used with the notes loaded of the file. since "number" section of the Name in the file represents ID of associated region. Unlike usually.
func (l Lyrics) GetNotesWithNumber(ID int) []comments.Comment {

	for _, note := range l.Sections {
		if note.Number == ID {
			return note.Comments
		}
	}

	return []comments.Comment{}
}

func (l Lyrics) GetNotesWithID(ID int) []comments.Comment {

	for _, s := range l.Sections {
		if s.ID == ID {
			return s.Comments
		}
	}

	return []comments.Comment{}
}

// looking for the first section with the same name and number, but not the language and returns it
// if no matches found returns original section
func (l Lyrics) GetAlternativeSection(s Section) Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number == alt.Number && s.Lang != alt.Lang {
			return alt
		}
	}

	return s
}

// looking for the first section with the same name and language but not number returns it
// if no matches found returns original section
func (l Lyrics) GetRelativeSection(s Section) Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number != alt.Number && s.Lang == alt.Lang {
			return alt
		}
	}

	return s
}

// lokking for a section with a same name but not the number and not the language
func (l Lyrics) GetRelativeAlternativeSection(s Section) Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number != alt.Number && s.Lang != alt.Lang {
			return alt
		}
	}

	return s
}

// looking for all section with a matching name
func (l Lyrics) GetAllRelativeSections(s Section) []Section {
	result := []Section{}
	for _, section := range l.Sections {
		if section.Name == s.Name {
			result = append(result, section)
		}
	}

	return result
}

// get section by iteration over Lyrics.Section and will retrun first section with matching id. If not found return an empty Section struct
func (l Lyrics) GetSectionByID(id int) Section {
	for _, sec := range l.Sections {
		if sec.ID == id {
			return sec
		}
	}

	return Section{}
}

// get section with a specified Section.Number by iteratin through all section. It will return the first one with a requested Number
func (l Lyrics) GetSectionByNumber(num int) Section {
	for _, sec := range l.Sections {
		if sec.Number == num {
			return sec
		}
	}

	return Section{}
}

func (l Lyrics) GetSectionByFullName(s Section) Section {
	for _, sec := range l.Sections {
		if sec.GetFullName() == s.GetFullName() {
			return sec
		}
	}

	return s
}

func (l Lyrics) GetOffsets(s Section) OffsetMap {
	offsets := NewOffsetMap(s)

	if !offsets.Empty() {
		return offsets
	}

	exact := l.GetSectionByID(s.ID)

	if exact.Empty() {
		exact = s
	}

	offsets = NewOffsetMap(exact)
	if !offsets.Empty() {
		return offsets
	}

	matchByName := l.GetSectionByFullName(s)
	offsets = NewOffsetMap(matchByName)
	if !offsets.Empty() {
		return offsets
	}

	alt := l.GetAlternativeSection(s)

	offsets = NewOffsetMap(alt)
	if !offsets.Empty() {
		return offsets
	}

	rel := l.GetRelativeSection(s)
	offsets = NewOffsetMap(rel)
	if !offsets.Empty() {
		return offsets
	}

	relalt := l.GetRelativeAlternativeSection(s)
	offsets = NewOffsetMap(relalt)
	if !offsets.Empty() {
		return offsets
	}

	all := l.GetAllRelativeSections(s)

	for _, relative := range all {
		offsets := NewOffsetMap(relative)

		if !offsets.Empty() {
			return offsets
		}
	}

	return offsets

}

type LyricLine struct {
	ID       int
	Position float64
	Length   float64
	Chords   chords.ChordLine
	Text     string
}

func (l LyricLine) Chordless() string {
	l.Text = strings.ReplaceAll(l.Text, ">", "")
	return strings.ReplaceAll(l.Text, "^", "")
}

func (l LyricLine) WithChordsString() string {
	if len(l.Text) < 1 {
		return l.Chords.String()
	}
	l.Text = strings.ReplaceAll(l.Text, ">", "^...")

	n := strings.Count(l.Text, "^")

	for i := range n {
		chord := ""
		if i < len(l.Chords) {
			chord = l.Chords[i].String()
		}
		l.Text = strings.Replace(l.Text, "^", chord, 1)
	}

	return l.Text
}

// it is so wrong that I have html in here. I should separate the logic and html
func (l LyricLine) WithChordsHTMLKey(key string) string {
	if len(l.Text) < 1 {
		return fmt.Sprintf("<span class=\"textless\">%s</span>", l.Chords.HTMLKey(key))
	}

	l.Text = strings.ReplaceAll(l.Text, ">", "^...")

	n := strings.Count(l.Text, "^")

	if n == 0 {

		l.Text = fmt.Sprintf("<span class=\"positionless\">%s</span>%s", l.Chords.HTMLKey(key), l.Text)
	}

	for i := range n {
		chord := ""
		if i < len(l.Chords) {
			chord = fmt.Sprintf("<span class=\"relative\"><span class=\"chord\">%s</span></span>", l.Chords[i].HTMLKey(key))
		}
		l.Text = strings.Replace(l.Text, "^", chord, 1)
	}

	reminder := fmt.Sprintf("<span class=\"relative\"><span class=\"chord\">")
	if diff := len(l.Chords) - n; diff > 0 && n > 0 {
		for i := n; i < len(l.Chords); i++ {
			reminder += l.Chords[i].HTMLKey(key) + " "

		}
	}

	reminder += "</span></span>"

	l.Text += reminder

	return l.Text
}

type Lines []LyricLine

func (l Lines) Empty() bool {
	return len(l) == 0
}

func (l Lyrics) String() string {
	var text string
	for _, s := range l.Sections {

		text = fmt.Sprintf("%s\n\n%s", text, s)

	}

	return text
}

func (l *Lyrics) Empty() bool {
	return len(l.Sections) == 0
}

type Section struct {
	ID       int
	Name     string
	Number   int
	Lang     string
	Comments []comments.Comment
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

// returns a plain text representation of a section
func (s Section) String() string {
	text := fmt.Sprintf("[ %s ]", s.GetFullName())

	for _, comment := range s.Comments {
		t := ""
		if comment.Type != comments.GENERIC {
			t = string(comment.Type) + ": "
		}
		text = fmt.Sprintf("%s\n--%s%s", text, t, comment)

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
