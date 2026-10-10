package comments

import (
	"strings"
)

type CommentType string

const (
	GENERIC CommentType = "generic"
	BASS    CommentType = "bass"
	DRUMS   CommentType = "drums"
	PIANO   CommentType = "piano"
	GUITAR  CommentType = "guitar"
)

type Comment struct {
	//Id of the section comment belongs too
	//Type of the comment.
	//It can be:
	//1. generic
	//2. bass
	//3. drums
	//4. piano
	//5. guitar
	Type CommentType
	//main content of the comment
	Text string
}

func (c Comment) String() string {
	return c.Text
}

func FilterComments(c []Comment, t CommentType) []Comment {
	result := []Comment{}

	for _, e := range c {
		if e.Type == t {
			result = append(result, e)
		}
	}

	return result
}

func GetCommentType(s string) CommentType {

	if s == "" {
		return GENERIC
	}

	types := []CommentType{
		GENERIC,
		BASS,
		DRUMS,
		PIANO,
		GUITAR,
	}

	for _, t := range types {
		if strings.ToLower(s) == string(t) {
			return t
		}
	}

	return GENERIC
}

// returns Comment struct with Type = GENERIC
func NewComment() Comment {
	return Comment{Type: GENERIC}
}

// returns Comment struct with Type = BASS
func NewBassComment() Comment {
	return Comment{
		Type: BASS,
	}
}

// returns Comment struct with Type = DRUMS
func NewDrumsComment() Comment {
	return Comment{
		Type: DRUMS,
	}
}

// returns Comment struct with Type = PIANO
func NewPianoComment() Comment {
	return Comment{
		Type: PIANO,
	}
}

// returns Comment struct with Type = GUITAR
func NewGuitarComment() Comment {
	return Comment{
		Type: GUITAR,
	}
}

func NewCommentText(text string) Comment {
	return Comment{Type: GENERIC, Text: text}
}

func NewBassCommentText(text string) Comment {
	return Comment{
		Type: BASS,
		Text: text,
	}
}

func NewDrumsCommentText(text string) Comment {
	return Comment{
		Type: DRUMS,
		Text: text,
	}
}

func NewPianoCommentText(text string) Comment {
	return Comment{
		Type: PIANO,
		Text: text,
	}
}

func NewGuitarCommentText(text string) Comment {
	return Comment{
		Type: GUITAR,
		Text: text,
	}
}
