package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestPrintAsciiMapCharacters(t *testing.T) {
	testCases := []struct {
		name           string
		input          asciiMap
		expectedOutput string
	}{
		{
			name:           "test print",
			input:          asciiMap{printContent: "ABCDEFG"},
			expectedOutput: "ABCDEFG",
		},
	}
	for _, testCase := range testCases {
		r, w, _ := os.Pipe()
		originalStdout := os.Stdout
		os.Stdout = w
		defer func() { os.Stdout = originalStdout }()
		testCase.input.PrintAsciiMapCharacters()
		w.Close()
		var buf bytes.Buffer
		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Fatal("Error reading from pipe:", err)
		}
		actualOutput := buf.String()
		if actualOutput != testCase.expectedOutput {
			t.Errorf("Expected output: %s, but got: %s", testCase.expectedOutput, actualOutput)
		}
	}
}
