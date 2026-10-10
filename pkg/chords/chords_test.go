package chords

import "testing"

func TestNewChord(t *testing.T) {

	input := []string{"1", "5sus4", "2m7", "1/3", "42", "5sus4/7"}

	tests := []struct {
		expectedRoot     int
		expectedRawChord string
		expectedMinor    bool
		expectedVariant  string
		expectedBass     int
	}{
		{
			expectedRawChord: "1",
			expectedRoot:     1,
			expectedMinor:    false,
			expectedVariant:  "",
			expectedBass:     0,
		},
		{
			expectedRawChord: "5sus4",
			expectedRoot:     5,
			expectedMinor:    false,
			expectedVariant:  "sus4",
			expectedBass:     0,
		},
		{
			expectedRawChord: "2m7",
			expectedRoot:     2,
			expectedMinor:    true,
			expectedVariant:  "7",
			expectedBass:     0,
		},
		{
			expectedRawChord: "1/3",
			expectedRoot:     1,
			expectedMinor:    false,
			expectedVariant:  "",
			expectedBass:     3,
		},
		{
			expectedRawChord: "42",
			expectedRoot:     4,
			expectedMinor:    false,
			expectedVariant:  "2",
			expectedBass:     0,
		},
		{
			expectedRawChord: "5sus4/7",
			expectedRoot:     5,
			expectedMinor:    false,
			expectedVariant:  "sus4",
			expectedBass:     7,
		},
	}

	for i, tt := range tests {
		c := NewChord(input[i])

		if c.RawChord != tt.expectedRawChord {
			t.Fatalf("tests[%d] RawChord: expected %s got %s", i, tt.expectedRawChord, c.RawChord)
		}

		if c.Root != tt.expectedRoot {
			t.Fatalf("tests[%d] failed RootChord: expected %d got %d", i, tt.expectedRoot, c.Root)
		}

		if c.Minor != tt.expectedMinor {
			t.Fatalf("tests[%d] failed. minor expected %v got %v", i, tt.expectedMinor, c.Minor)
		}

		if c.Variant != tt.expectedVariant {
			t.Fatalf("tests[%d] failed. major expected %v got %v", i, tt.expectedVariant, c.Variant)
		}
	}

}
