package tetris

import (
	"testing"
)

func TestCheckValidTestCoordinates(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"../samples/badexample01.txt", false},
		{"../samples/badexample02.txt", false},
		{"../samples/badexample04.txt", false},
		{"../samples/goodexample00.txt", true},
		{"../samples/goodexample01.txt", true},
		{"../samples/goodexample02.txt", true},
		{"../samples/goodexample03.txt", true},
		{"../samples/goodexample04.txt", true},
		{"../samples/goodexample05.txt", true},
		{"../samples/hardexam.txt", true},
	}
	for _, test := range tests {
		f := ReadSample(test.input)

		got := true
		for _, tile := range f {
			var coordinates [][]int
			for y, line := range tile {
				for x, symbol := range line {
					if symbol == "#" {
						coordinates = append(coordinates, []int{x, y})
					}
				}
			}
			got = checkValidCoordinates(coordinates)
			if got == false {
				break
			}
		}

		if got != test.want {
			t.Errorf("For %s expected: %v, got %v\n", test.input, test.want, got)
		}
	}
}
