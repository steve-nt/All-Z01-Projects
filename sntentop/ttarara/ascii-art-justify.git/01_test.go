package main

import (
	"errors"
	"os"
	"testing"
)

// Test for Args
func TestCheckArgs(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		wantErr error
	}{
		{"Arguments < 2", []string{"program"}, errors.New("invalid number of arguments")},
		{"Valid args (min)", []string{"program", "arg2"}, nil},
		{"Valid args (max)", []string{"program", "arg2", "arg3", "arg4"}, nil},
		{"Arguments > 4", []string{"program", "arg2", "arg3", "arg4", "arg5"}, errors.New("invalid number of arguments")},
		{"Empty args", []string{"program"}, errors.New("invalid number of arguments")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set os.Args to the test case arguments
			os.Args = tc.args

			// Call checkArgs and capture the result
			err := checkArgs()

			// Check if the error matches the expected outcome
			// First if statement: Detects cases where the presence or absence of an error itself was unexpected.
			if (err != nil && tc.wantErr == nil) || (err == nil && tc.wantErr != nil) {
				t.Errorf("checkArgs() = %v; want %v", err, tc.wantErr)
			}
			//Second if statement: Ensures that, when an error is expected and returned, its message matches the expectation.
			if err != nil && tc.wantErr != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("checkArgs() error = %v; want %v", err, tc.wantErr)

				//By including these checks, the test ensures checkArgs not only returns errors when expected, but also with the correct message.
			}
		})
	}

}

//-----------------------------------------------------------------------------------

// Test for Align flag

func TestCheckForAlign(t *testing.T) {
	// Test cases for checkForAlign
	testCases := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{"Valid args", []string{"program_name", "arg2"}, nil},
		{"Invalid align flag", []string{"program_name", "--align"}, errors.New("not a valid flag")},
		{"Invalid -align flag", []string{"program_name", "-align"}, errors.New("not a valid flag")},
		{"Invalid -align=value flag", []string{"program_name", "-align=center"}, errors.New("not a valid flag")},
		{"Valid args without align", []string{"program_name", "validArg"}, nil},
	}

	// Loop through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set os.Args to simulate command-line arguments
			os.Args = tc.args

			// Call checkForAlign
			err := checkForAlign()

			// Check if the returned error matches the expected error
			if err != nil && tc.expectedErr == nil {
				t.Errorf("For args %v, expected no error, but got %v", tc.args, err)
			} else if err == nil && tc.expectedErr != nil {
				t.Errorf("For args %v, expected error %v, but got nil", tc.args, tc.expectedErr)
			} else if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("For args %v, expected error %v, but got %v", tc.args, tc.expectedErr, err)
			}
		})
	}
}

//-----------------------------------------------------------------------------------

// Test for alignment values

func TestIsValidAlignment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		align         string
		expectedError error
	}{
		{"Valid left align", leftAlign, nil},
		{"Valid center align", centerAlign, nil},
		{"Valid right align", rightAlign, nil},
		{"Valid justify align", justifyAlign, nil},
		{"Invalid align value", "invalid", errors.New("invalid value")},
	}

	// Loop through all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the isValidAlignment function
			err := isValidAlignment(tc.align)

			// Check if the error matches the expected error
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("For align '%s', expected error '%v', but got '%v'", tc.align, tc.expectedError, err)
			} else if err == nil && tc.expectedError != nil {
				t.Errorf("For align '%s', expected error '%v', but got nil", tc.align, tc.expectedError)
			}
		})
	}
}

//-----------------------------------------------------------------------------------

// Test for ascii string input

func TestIsASCII(t *testing.T) {
	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{"Valid ASCII input", []string{"Hello, World!"}, false},
		{"Input with non-ASCII character", []string{"Hello, 世界!"}, true},
		{"Empty input", []string{""}, false},
		{"Single non-ASCII character", []string{"€"}, true},
	}

	// Loop through all test cases
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Call isASCII with the test args
			err := isASCII(test.args[0]) // Passing the string as input

			// Check if the result matches the expected error state
			if (err != nil) != test.expectErr {
				t.Errorf("For input %q, expected error: %v, got: %v", test.args[0], test.expectErr, err)
			}
		})
	}
}

//-----------------------------------------------------------------------------------

// Test for banners

func TestGetBannerType(t *testing.T) {
	testBanners := []struct {
		name          string
		args          []string
		expectedValue string
	}{
		{"No arguments", []string{"program_name"}, "standard"},
		{"Shadow", []string{"program_name", "shadow"}, "shadow"},
		{"Thinkertoy", []string{"program_name", "thinkertoy"}, "thinkertoy"},
		{"Invalid Banner", []string{"program_name", "invalid"}, "standard"},
		{"Shadow at end", []string{"program_name", "something", "shadow"}, "shadow"},
		{"Thinkertoy at end", []string{"program_name", "something", "thinkertoy"}, "thinkertoy"},
	}

	// Loop through all test cases
	for _, test := range testBanners {
		t.Run(test.name, func(t *testing.T) {
			// Call getBannerType with the test args
			result := getBannerType(test.args)

			// Check if the result matches the expected value
			if result != test.expectedValue {
				t.Errorf("For args %v, expected bannerType to be '%s', but got '%s'", test.args, test.expectedValue, result)
			}
		})
	}
}

//-----------------------------------------------------------------------------------
