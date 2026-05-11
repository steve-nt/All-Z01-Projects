package utils

import (
	"errors"
	"testing"
)

func TestInternalCheckFileIntegrity(t *testing.T) {
	testCases := []struct {
		name           string
		mockCreateHash func(string) (string, error)
		expectedErr    bool
	}{
		{
			name: "Matching Hash",
			mockCreateHash: func(string) (string, error) {
				return hash, nil
			},
			expectedErr: false,
		},
		{
			name: "Mismatching Hash",
			mockCreateHash: func(string) (string, error) {
				return "wrong hash", nil
			},
			expectedErr: true,
		},
		{
			name: "CreateHash error",
			mockCreateHash: func(string) (string, error) {
				return "", errors.New("createHash Error")
			},
			expectedErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := checkFileIntegrity(testCase.mockCreateHash)
			if (err != nil) != testCase.expectedErr {
				t.Errorf("CheckFileIntegrity() error: %v, expected error: %v", err, testCase.expectedErr)
			}
		})
	}
}

func TestCreateHash(t *testing.T) {
	testCases := []struct {
		name         string
		filePath     string
		expectedHash string
		expectedErr  bool
	}{
		{
			name:         "Valid File",
			filePath:     mapPath,
			expectedHash: hash,
			expectedErr:  false,
		},
		{
			name:         "Invalid File",
			filePath:     "path/to/invalid/file",
			expectedHash: "",
			expectedErr:  true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultHash, err := createHash(testCase.filePath)
			if resultHash != testCase.expectedHash {
				t.Errorf("Error Hash: %v, expectedHash: %v", resultHash, testCase.expectedHash)
			}
			if (err != nil) != testCase.expectedErr {
				t.Errorf("createHash() error: %v, expectedErr: %v", err, testCase.expectedErr)
			}
		})
	}
}

func TestCheckFileIntegrity(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr bool
	}{
		{
			name:        "Valid Operation",
			expectedErr: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := CheckFileIntegrity()
			if (err != nil) != testCase.expectedErr {
				t.Errorf("CheckFIleIntegrity() error: %v, expectedErr: %v", err, testCase.expectedErr)
			}
		})
	}
}
