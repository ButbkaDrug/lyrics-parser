package parser

import (
	"lp/internal/lexer"
	"lp/pkg/lyrics"
	"testing"
)

func TestTypicalInput(t *testing.T) {
	input := `[ Verse 1 ]
--22.857143
{0, 4.571429}
Water You turned into wine
{4.571429, 2.857143}
Opened the eyes of the blind
{7.428572, 4.58}
There's no one like You
{12, 4}
None like You
None like You`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	expectedSection := lyrics.Section{
		Name:     "Verse",
		Number:   1,
		Comments: []string{"--22.857143"},
		Text: []lyrics.TextLine{
			{Offset: 0, Length: 4.571429, Text: "Water You turned into wine"},
			{Offset: 4.571429, Length: 2.857143, Text: "Opened the eyes of the blind"},
			{Offset: 7.428572, Length: 4.58, Text: "There's no one like You"},
			{Offset: 12, Length: 4, Text: "None like You"},
			{Offset: 0, Length: 0, Text: "None like You"},
		},
	}

	output := p.Parse()

	out := output.Sections[0]

	if expectedSection.Name != out.Name {
		t.Fatalf("failty .Name field. expected %v got %v", expectedSection.Name, output.Name)
	}
	if expectedSection.Number != out.Number {
		t.Fatalf("failty .Number field. expected %v got %v", expectedSection.Number, out.Number)
	}

	for i, tt := range expectedSection.Text {

		if tt.Text != out.Text[i].Text {
			t.Fatalf("failty Text[%d].Text field. expected %v got %v", i, tt.Text, out.Text[i].Text)

		}
		if tt.Offset != out.Text[i].Offset {
			t.Fatalf("failty Text[%d].Offset field. expected %v got %v", i, tt.Offset, out.Text[i].Offset)

		}
		if tt.Length != out.Text[i].Length {
			t.Fatalf("failty Text[%d].Length field. expected %v got %v", i, tt.Length, out.Text[i].Length)

		}

	}

}
