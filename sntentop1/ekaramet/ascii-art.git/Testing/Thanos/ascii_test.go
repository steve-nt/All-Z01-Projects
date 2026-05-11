package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadBanner(t *testing.T) {
	// Create a test banner file
	bannerFile := "test_banner.txt"
	bannerContent := " ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n\n"
	err := ioutil.WriteFile(bannerFile, []byte(bannerContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(bannerFile)

	// Load the banner
	bannerMap, err := loadBanner(bannerFile)
	if err != nil {
		t.Errorf("loadBanner returned error: %v", err)
	}

	// Check the banner map
	if len(bannerMap) != 1 {
		t.Errorf("expected 1 character in banner map, got %d", len(bannerMap))
	}
	for char, art := range bannerMap {
		if char != ' ' {
			t.Errorf("expected character ' ', got '%c'", char)
		}
		if len(art) != 8 {
			t.Errorf("expected 8 lines of ASCII art, got %d", len(art))
		}
		for _, line := range art {
			if line != " ABCDEFGH" {
				t.Errorf("expected ASCII art line ' ABCDEFGH', got '%s'", line)
			}
		}
	}
}

func TestRenderString(t *testing.T) {
	// Create a test banner map
	bannerMap := map[rune][]string{
		'A': {"AAAAAAA", "A     A", "AAAAAAA", "A     A", "A     A", "A     A", "A     A", "AAAAAAA"},
		'B': {"BBBBBBB", "B     B", "BBBBBBB", "B     B", "B     B", "B     B", "B     B", "BBBBBBB"},
	}

	// Render a string
	input := "AB"
	output, err := renderString(input, bannerMap)
	if err != nil {
		t.Errorf("renderString returned error: %v", err)
	}

	// Check the output
	if len(output) != 8 {
		t.Errorf("expected 8 lines of output, got %d", len(output))
	}
	for i, line := range output {
		if i < 8 {
			if line != "AAAAAAA BBBBBBB" {
				t.Errorf("expected output line '%s', got '%s'", "AAAAAAA BBBBBBB", line)
			}
		} else {
			t.Errorf("expected only 8 lines of output, got %d", len(output))
		}
	}
}

func TestGetAvailableBanners(t *testing.T) {
	// Create a test banners directory
	bannersDir := "test_banners"
	err := os.Mkdir(bannersDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(bannersDir)

	// Create some test banner files
	bannerFiles := []string{"banner1.txt", "banner2.txt", "banner3.txt"}
	for _, file := range bannerFiles {
		err := ioutil.WriteFile(filepath.Join(bannersDir, file), []byte(""), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get the available banners
	banners, err := getAvailableBanners(bannersDir)
	if err != nil {
		t.Errorf("getAvailableBanners returned error: %v", err)
	}

	// Check the banners
	if len(banners) != len(bannerFiles) {
		t.Errorf("expected %d banners, got %d", len(bannerFiles), len(banners))
	}
	for i, banner := range banners {
		if banner != strings.TrimSuffix(bannerFiles[i], ".txt") {
			t.Errorf("expected banner '%s', got '%s'", strings.TrimSuffix(bannerFiles[i], ".txt"), banner)
		}
	}
}

func TestMain(t *testing.T) {
	// Create a test banner file
	bannerFile := "test_banner.txt"
	bannerContent := " ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n ABCDEFGH\n\n"
	err := ioutil.WriteFile(bannerFile, []byte(bannerContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(bannerFile)

	// Create a test input string
	inputString := "A"

	// Run the main function with the test input
	oldArgs := os.Args
	os.Args = []string{"ascii-art", inputString, bannerFile}
	defer func() { os.Args = oldArgs }()

	// Capture the output of the main function
	output := bytes.NewBufferString("")
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	main()

	// Check the output
	expectedOutput := "AAAAAAA \nA     A \nAAAAAAA \nA     A \nA     A \nA     A \nA     A \nAAAAAAA \n"
	if output.String() != expectedOutput {
		t.Errorf("expected output '%s', got '%s'", expectedOutput, output.String())
	}
}
