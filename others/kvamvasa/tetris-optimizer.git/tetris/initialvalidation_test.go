package tetris

import (
	"testing"
)

func TestInitialValidation(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"../samples/badformat.txt", false},
		{"../samples/badformat01.txt", false},
		{"../samples/badformat02.txt", false},
		{"../samples/badformat03.txt", false},
		{"../samples/badexample00.txt", false},
		{"../samples/badexample01.txt", true}, // true for now, should be invalidated later
		{"../samples/badexample02.txt", true}, // true for now, should be invalidated later
		{"../samples/badexample03.txt", false},
		{"../samples/badexample04.txt", true}, // true for now, should be invalidated later
		{"../samples/goodexample00.txt", true},
		{"../samples/goodexample01.txt", true},
		{"../samples/goodexample02.txt", true},
		{"../samples/goodexample03.txt", true},
		{"../samples/goodexample04.txt", true},
		{"../samples/goodexample05.txt", true},
		{"../samples/hardexam.txt", true},
		{"../samples/example.txt", true},
	}
	for _, test := range tests {
		f := ReadSample(test.input)

		if got := initialValidation(f); got != test.want {
			t.Errorf("For %s expected: %v, got %v\n", test.input, test.want, got)
		}
	}
}
