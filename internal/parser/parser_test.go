package parser

import (
	"log/slog"
	"lp/internal/lexer"
	"lp/internal/lyrics"
	"testing"
)

func TestParseOffset(t *testing.T) {
	input := `{2.333, 3.333}`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	offset, length := p.parseOffset()
	expectedOffset := 2.333
	expectedLength := 3.333

	if offset != expectedOffset {
		t.Fatalf("wrong offset expected %f got %f", expectedOffset, offset)
	}

	if length != expectedLength {
		t.Fatalf("wrong length. expected %f got %f", expectedLength, length)
	}

}

func TestParseSectionTitle(t *testing.T) {
	input := `[Verse]
    [Chorus]
    [Tag]
    [Verse 2]
    [Instrumental]
    [Bridge R]
    [Verse 2 R]
    `

	l := lexer.NewLexer(input)
	p := NewParser(l)

	tests := []struct {
		expectedName   string
		expectedNumber int
		expectedLang   string
	}{
		{
			expectedName:   "Verse",
			expectedNumber: 0,
			expectedLang:   "",
		},
		{
			expectedName:   "Chorus",
			expectedNumber: 0,
			expectedLang:   "",
		},
	}

	output := p.Parse()

	for i, tt := range tests {
		name := output.Sections[i].Name
		if tt.expectedName != name {
			t.Fatalf("tests[%d] failed. expected name %s got %s", i, tt.expectedName, name)
		}
	}

}
func TestParseSection(t *testing.T) {
	input := `[ Verse 1 R ] --Drums In
    {0, 3.333} I want to screem it out
    {3.333, 3.333} From Every Mountain top
    {6.666, 3.333} Your goodness knows no bounds
    {9.999, 3.333} Your goodness never stops

    [Chorus] --All In []{}, 1234
    [Tag]
    --Breakdown
    --Base comes in on the second half
    [Instrumental]
    --Пианино и скрипка
    [Verse 2 R]
    Хочу я прокречать,
    С высоких гор - вершин
    Что Благодать твоя!
    Безмерна. Безконца?

    `
	slog.SetLogLoggerLevel(slog.LevelDebug)
	tests := []lyrics.Section{
		{
			Name:     "Verse",
			Number:   1,
			Lang:     "R",
			Comments: []string{"Drums In"},
			Text: []lyrics.TextLine{
				{Offset: 0, Length: 3.333, Text: "I want to screem it out"},
				{Offset: 3.333, Length: 3.333, Text: "From Every Mountain top"},
				{Offset: 6.666, Length: 3.333, Text: "Your goodness knows no bounds"},
				{Offset: 9.999, Length: 3.333, Text: "Your goodness never stops"},
			},
		},
		{
			Name:     "Chorus",
			Comments: []string{"All In []{}, 1234"},
		},
		{
			Name:     "Tag",
			Comments: []string{"Breakdown", "Base comes in on the second half"},
		},
		{
			Name:     "Instrumental",
			Comments: []string{"Пианино и скрипка"},
		},
		{
			Name:   "Verse",
			Number: 2,
			Lang:   "R",
			Text: []lyrics.TextLine{
				{Offset: 0, Length: 0, Text: "Хочу я прокречать,"},
				{Offset: 0, Length: 0, Text: "С высоких гор - вершин"},
				{Offset: 0, Length: 0, Text: "Что Благодать твоя!"},
				{Offset: 0, Length: 0, Text: "Безмерна. Безконца?"},
			},
		},
	}

	testParseSection(input, tests, t)

}

func testParseSection(input string, tests []lyrics.Section, t *testing.T) {

	l := lexer.NewLexer(input)
	p := NewParser(l)

	lx := p.Parse()

	for i, tt := range tests {
		s := lx.Sections[i]
		if tt.Name != s.Name {
			t.Fatalf("tests[%d] wrong name. expected %q got %q", i, tt.Name, s.Name)
		}

		if tt.Number != s.Number {
			t.Fatalf("tests[%d] wrong number. expected %d got %d", i, tt.Number, s.Number)
		}

		if tt.Lang != s.Lang {
			t.Fatalf("tests[%d] wrong lang. expected %q got %q", i, tt.Lang, s.Lang)
		}

		if len(tt.Comments) != len(s.Comments) {
			t.Fatalf("test[%d] number of comments doesn't match. expected %d got %d", i, len(tt.Comments), len(s.Comments))
		}

		for j, tc := range tt.Comments {
			comment := s.Comments[j]

			if tc != comment {
				t.Fatalf("test[%d] comment[%d] doen't match. expected %s got %s",
					i,
					j,
					tc,
					comment,
				)
			}
		}

		if len(tt.Text) != len(s.Text) {
			t.Fatalf("tests[%d] number of text line doesn't match. expected %d got %d.\n\n%v\n\n%v",
				i, len(tt.Text),
				len(s.Text),
				tt.Text,
				s.Text,
			)
		}

		for j, tl := range tt.Text {
			line := s.Text[j]

			if tl.Offset != line.Offset {
				t.Fatalf("tests[%d]. Line[%d] offset doesn't match. expected %v got %v",
					i,
					j,
					tl.Offset,
					line.Offset,
				)
			}
			if tl.Length != line.Length {
				t.Fatalf("tests[%d]. Line[%d] length doesn't match. expected %v got %v",
					i,
					j,
					tl.Length,
					line.Length,
				)
			}
			if tl.Text != line.Text {
				t.Fatalf("tests[%d]. Line[%d] text doesn't match. expected %v got %v",
					i,
					j,
					tl.Text,
					line.Text,
				)
			}
		}
	}

}
