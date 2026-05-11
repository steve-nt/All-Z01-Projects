package utils

import (
	"testing"
)

var testMap = asciiMap{
	content: map[rune][]string{
		'1': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'2': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'3': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'4': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'5': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'6': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'7': {"1", "2", "3", "4", "5", "6", "7", "8"},
		'8': {"1", "2", "3", "4", "5", "6", "7", "8"},
	},
}

func TestAsciiOutput(t *testing.T) {
	testCases := []struct {
		name           string
		inputMap       asciiMap
		input          string
		expectedOutput []string
	}{
		{
			name:           "Empty string",
			inputMap:       testMap,
			input:          "",
			expectedOutput: []string{},
		},
		{
			name:           "Valid string",
			inputMap:       testMap,
			input:          "12345678",
			expectedOutput: []string{"11111111", "22222222", "33333333", "44444444", "55555555", "66666666", "77777777", "88888888"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := asciiOutput(testCase.input, &testCase.inputMap)
			if len(output) != len(testCase.expectedOutput) {
				t.Errorf("Error output has lenght: %v, expected lenght: %v", len(output), len(testCase.expectedOutput))
			}
			for elem := range output {
				if output[elem] != testCase.expectedOutput[elem] {
					t.Errorf("Error output: %v, expected: %v", output, testCase.expectedOutput)
				}
			}
		})
	}
}

var mockasciiOutput = func(string, *asciiMap) []string {
	return []string{"11111111", "22222222", "33333333", "44444444", "55555555", "66666666", "77777777", "88888888"}
}

func TestInternalFormatAsciiArt(t *testing.T) {
	testCases := []struct {
		name           string
		inputMap       asciiMap
		mockFunc       func(string, *asciiMap) []string
		expectedOutput string
	}{
		{
			name:     "Valid input",
			inputMap: asciiMap{input: "1235678", printContent: ""},
			mockFunc: mockasciiOutput,
			expectedOutput: `11111111
22222222
33333333
44444444
55555555
66666666
77777777
88888888
`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.inputMap.formatAsciiArt(testCase.mockFunc)
			if testCase.inputMap.printContent != testCase.expectedOutput {
				t.Errorf("Error output: %v, expected: %v", testCase.inputMap.printContent, testCase.expectedOutput)
			}
		})
	}
}

func TestFormatAsciiArt(t *testing.T) {
	testCases := []struct {
		name           string
		inputMap       asciiMap
		expectedOutput string
	}{
		{
			name:     "Valid Input",
			inputMap: asciiMap{input: "12345678", printContent: "", content: testMap.content},
			expectedOutput: `11111111
22222222
33333333
44444444
55555555
66666666
77777777
88888888
`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.inputMap.FormatAsciiArt()
			if testCase.inputMap.printContent != testCase.expectedOutput {
				t.Errorf("Error output: %v, expectedOutput: %v", []byte(testCase.inputMap.printContent), []byte(testCase.expectedOutput))
			}
		})
	}
}
