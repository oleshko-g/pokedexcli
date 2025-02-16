package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input         string
		expectedWords []string
	}{
		{
			input:         "",
			expectedWords: []string{},
		},
		{
			input:         "	Hello, World!",
			expectedWords: []string{"hello,", "world!"},
		},
	}

	for i, c := range cases {
		actualWords := cleanInput(c.input)
		if len(actualWords) != len(c.expectedWords) {
			t.Errorf("Failed at test case: %d. Expected return length: %d. Actual return length: %d", i, len(c.expectedWords), len(actualWords))
			t.FailNow()
		}

		for j := 0; j < len(actualWords); j++ {
			if actualWords[j] != c.expectedWords[j] {
				t.Errorf("Failed at test case: %d. Expected word at index [%d]: %#v. Actual word at index %d: %#v.", i, j, c.expectedWords[j], j, actualWords[j])
				t.FailNow()
			}
		}
	}
}
