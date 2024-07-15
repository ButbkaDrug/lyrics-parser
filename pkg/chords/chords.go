package chords

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// represents a chord in a number system. Has helper function to represent chord in any key. Chord should know how print itself. And also how to transpose itself to any given key.
type Chord struct {
	//Chord as it's been passed befor parsing. Can be empty.
	RawChord string
	//bass note if different from the root
	Bass int
	//root note of the chord
	Root int
	// true if major
	Minor bool
	//can be minor, suspended, seventh or bunch of other chords
	Variant string
}

// return a stringified version of a chord
func (c Chord) String() string {
	str := fmt.Sprintf("%d", c.Root)

	if c.Minor {
		str += fmt.Sprint("m")
	}

	if c.Variant != "" {
		str += fmt.Sprintf("%s", c.Variant)
	}

	if c.Bass != 0 {
		str += fmt.Sprintf("/%d", c.Bass)
	}

	return str
}

// return HTML version of a chord. whre .Variant put in <sup></sup> tag
func (c Chord) HTML() string {

	html := fmt.Sprintf("%d", c.Root)

	if c.Minor {
		html += "m"
	}

	if c.Variant != "" {
		html += fmt.Sprintf("<sup>%s</sup>", c.Variant)
	}

	if c.Bass > 0 {
		html += fmt.Sprintf("/%d", c.Bass)
	}

	return html

}

// returns html version of a chord in a given key. if key doen't exist return results of a .HTML() function
func (c Chord) HTMLKey(key string) string {
	keys := map[string]map[int]string{
		"A":  {1: "A", 2: "B", 3: "C#", 4: "D", 5: "E", 6: "F#", 7: "G#"},
		"Bb": {1: "Bb", 2: "C", 3: "D", 4: "Eb", 5: "F", 6: "G", 7: "A"},
		"B":  {1: "B", 2: "C#", 3: "D#", 4: "E", 5: "F#", 6: "G#", 7: "A#"},
		"C":  {1: "C", 2: "D", 3: "E", 4: "F", 5: "G", 6: "A", 7: "B"},
		"Db": {1: "Db", 2: "Eb", 3: "F", 4: "Gb", 5: "Ab", 6: "Bb", 7: "C"},
		"D":  {1: "D", 2: "E", 3: "F#", 4: "G", 5: "A", 6: "B", 7: "C#"},
		"Eb": {1: "Eb", 2: "F", 3: "G", 4: "Ab", 5: "Bb", 6: "C", 7: "D"},
		"E":  {1: "E", 2: "F#", 3: "G#", 4: "A", 5: "B", 6: "C#", 7: "D#"},
		"F":  {1: "F", 2: "G", 3: "A", 4: "Bb", 5: "C", 6: "D", 7: "E"},
		"Gb": {1: "Gb", 2: "Ab", 3: "Bb", 4: "Cb", 5: "Db", 6: "Eb", 7: "F"},
		"G":  {1: "G", 2: "A", 3: "B", 4: "C", 5: "D", 6: "E", 7: "F#"},
		"Ab": {1: "Ab", 2: "Bb", 3: "C", 4: "Db", 5: "Eb", 6: "F", 7: "G"},
	}

	if keys[key] == nil {
		return c.HTML()
	}

	html := fmt.Sprintf("%s", keys[key][c.Root])

	if c.Minor {
		html += "m"
	}

	if c.Variant != "" {
		html += fmt.Sprintf("<sup>%s</sup>", c.Variant)
	}

	if c.Bass > 0 {
		html += fmt.Sprintf("/%s", keys[key][c.Bass])
	}

	return html
}

func NewChord(chord string) Chord {
	c := Chord{
		RawChord: chord,
	}
	start := 1
	end := len(chord)

	if len(chord) < 1 {
		return c
	}

	rootString := string(chord[0])

	root, err := strconv.Atoi(rootString)

	if err != nil {
		slog.Error("failed to parse chord. root is not a number", "chord", chord)
		return c
	}

	c.Root = root
	idx := strings.IndexByte(chord, '/')

	if idx > 0 && len(chord) > idx {
		end = idx
		bassString := string(chord[end+1])
		bass, err := strconv.Atoi(bassString)

		if err != nil {
			slog.Error("failed to parse chord bass . bass is not a number", "chord", chord)
		} else {
			c.Bass = bass
		}
	}

	varianString := chord[start:end]

	if varianString == "m" || varianString == "m7" {
		c.Minor = true
		start = 2
	}

	c.Variant = chord[start:end]

	return c
}

// represents a slice of chords. It should have a String() method to be easily convertable to a string line if needed
type ChordLine []Chord
type Chords []ChordLine

func (c Chords) Empty() bool {
	return len(c) == 0
}

func (c Chords) String() string {
	text := ""
	for _, chord := range c {
		text += chord.String()
	}

	return text
}

func (c Chords) Pack() ChordLine {

	line := ChordLine{}

	for _, chord := range c {

		line = append(line, chord...)

	}

	return line

}

func (c ChordLine) String() string {
	str := ""

	for _, chord := range c {
		str += fmt.Sprintf("%s ", chord.String())
	}

	return strings.Trim(str, " ")
}

func (c ChordLine) HTML() string {
	str := ""

	for _, chord := range c {
		str += fmt.Sprintf("%s ", chord.HTML())
	}

	return strings.Trim(str, " ")
}

func (c ChordLine) HTMLKey(key string) string {
	str := ""

	for _, chord := range c {
		str += fmt.Sprintf("%s ", chord.HTMLKey(key))
	}

	return strings.Trim(str, " ")
}
