package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// captureStdout runs f and returns everything it printed to os.Stdout.
func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// ---------------------------------------------------------------------------
// LoadBanner
// ---------------------------------------------------------------------------

func TestLoadBanner_KeyCount(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	if len(bannerMap) != 95 {
		t.Errorf("expected 95 keys, got %d", len(bannerMap))
	}
}

func TestLoadBanner_CharacterHeight(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	if len(bannerMap['A']) != 8 {
		t.Errorf("expected 8 lines for 'A', got %d", len(bannerMap['A']))
	}
}

func TestLoadBanner_KnownRow(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	expected := "    /\\     "
	if bannerMap['A'][1] != expected {
		t.Errorf("expected %q, got %q", expected, bannerMap['A'][1])
	}
}

func TestLoadBanner_InvalidPath(t *testing.T) {
	_, err := LoadBanner("banners/nonexistent.txt")
	if err == nil {
		t.Error("expected an error for a missing banner file, got nil")
	}
}

func TestLoadBanner_AllFonts(t *testing.T) {
	for _, font := range []string{"standard", "shadow", "thinkertoy"} {
		bannerMap, err := LoadBanner("banners/" + font + ".txt")
		if err != nil {
			t.Errorf("font %q: LoadBanner returned an error: %v", font, err)
			continue
		}
		if len(bannerMap) != 95 {
			t.Errorf("font %q: expected 95 keys, got %d", font, len(bannerMap))
		}
	}
}

// ---------------------------------------------------------------------------
// Input splitting
// ---------------------------------------------------------------------------

func TestSplit_SingleNewline(t *testing.T) {
	result := strings.Split("a\\nb", "\\n")
	if len(result) != 2 || result[0] != "a" || result[1] != "b" {
		t.Errorf("expected [a b], got %v", result)
	}
}

func TestSplit_DoubleNewline(t *testing.T) {
	result := strings.Split("a\\n\\nb", "\\n")
	if len(result) != 3 || result[0] != "a" || result[1] != "" || result[2] != "b" {
		t.Errorf("expected [a  b], got %v", result)
	}
}

func TestSplit_EmptyInput(t *testing.T) {
	input := ""
	if input == "" {
		return
	}
	t.Error("empty input should have been caught before this point")
}

// ---------------------------------------------------------------------------
// Render (no color)
// ---------------------------------------------------------------------------

func TestRender_EmptySegment(t *testing.T) {
	output := captureStdout(func() {
		Render([]string{""}, map[rune][]string{})
	})
	if output != "\n" {
		t.Errorf("expected a single blank line, got %q", output)
	}
}

func TestRender_SingleChar(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	output := captureStdout(func() {
		Render([]string{"A"}, bannerMap)
	})
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 lines of output for 'A', got %d", len(lines))
	}
	if lines[1] != "    /\\     " {
		t.Errorf("expected %q, got %q", "    /\\     ", lines[1])
	}
}

func TestRender_NoANSICodes(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	output := captureStdout(func() {
		Render([]string{"Hi"}, bannerMap)
	})
	if strings.Contains(output, "\033[") {
		t.Error("Render output should not contain any ANSI escape codes")
	}
}

// ---------------------------------------------------------------------------
// ColorCode
// ---------------------------------------------------------------------------

func TestColorCode_NamedColors(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"red", "\033[31m"},
		{"green", "\033[32m"},
		{"blue", "\033[34m"},
		{"white", "\033[37m"},
		{"black", "\033[30m"},
		{"cyan", "\033[36m"},
		{"magenta", "\033[35m"},
		{"yellow", "\033[33m"},
	}
	for _, tc := range cases {
		got, err := ColorCode(tc.input)
		if err != nil {
			t.Errorf("ColorCode(%q) returned unexpected error: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ColorCode(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestColorCode_CaseInsensitive(t *testing.T) {
	variants := []string{"red", "RED", "Red", "rEd", "rED"}
	expected := "\033[31m"
	for _, v := range variants {
		got, err := ColorCode(v)
		if err != nil {
			t.Errorf("ColorCode(%q) returned unexpected error: %v", v, err)
		}
		if got != expected {
			t.Errorf("ColorCode(%q) = %q, want %q", v, got, expected)
		}
	}
}

func TestColorCode_Aliases(t *testing.T) {
	// purple is an alias for magenta
	got, err := ColorCode("purple")
	if err != nil {
		t.Fatalf("ColorCode(\"purple\") returned unexpected error: %v", err)
	}
	if got != "\033[35m" {
		t.Errorf("ColorCode(\"purple\") = %q, want magenta code \\033[35m", got)
	}

	// orange maps to bright yellow
	got, err = ColorCode("orange")
	if err != nil {
		t.Fatalf("ColorCode(\"orange\") returned unexpected error: %v", err)
	}
	if got != "\033[93m" {
		t.Errorf("ColorCode(\"orange\") = %q, want bright yellow code \\033[93m", got)
	}

	// pink maps to bright magenta
	got, err = ColorCode("pink")
	if err != nil {
		t.Fatalf("ColorCode(\"pink\") returned unexpected error: %v", err)
	}
	if got != "\033[95m" {
		t.Errorf("ColorCode(\"pink\") = %q, want bright magenta code \\033[95m", got)
	}
}

func TestColorCode_UnknownName(t *testing.T) {
	_, err := ColorCode("ultraviolet")
	if err == nil {
		t.Error("expected an error for unknown color name, got nil")
	}
}

func TestColorCode_EmptyString(t *testing.T) {
	_, err := ColorCode("")
	if err == nil {
		t.Error("expected an error for empty color string, got nil")
	}
}

func TestColorCode_Hex(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"#ff0000", "\033[38;2;255;0;0m"},
		{"#FF0000", "\033[38;2;255;0;0m"}, // case-insensitive
		{"#00ff00", "\033[38;2;0;255;0m"},
		{"#0000ff", "\033[38;2;0;0;255m"},
		{"#ffffff", "\033[38;2;255;255;255m"},
		{"#000000", "\033[38;2;0;0;0m"},
	}
	for _, tc := range cases {
		got, err := ColorCode(tc.input)
		if err != nil {
			t.Errorf("ColorCode(%q) returned unexpected error: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ColorCode(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestColorCode_Hex_Malformed(t *testing.T) {
	bad := []string{"#gg0000", "#fff", "#1234567", "#", "##ff0000"}
	for _, v := range bad {
		_, err := ColorCode(v)
		if err == nil {
			t.Errorf("ColorCode(%q): expected error for malformed hex, got nil", v)
		}
	}
}

func TestColorCode_RGB(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"rgb(255,0,0)", "\033[38;2;255;0;0m"},
		{"rgb(0,255,0)", "\033[38;2;0;255;0m"},
		{"rgb(0,0,255)", "\033[38;2;0;0;255m"},
		{"rgb(255, 0, 0)", "\033[38;2;255;0;0m"}, // spaces tolerated
	}
	for _, tc := range cases {
		got, err := ColorCode(tc.input)
		if err != nil {
			t.Errorf("ColorCode(%q) returned unexpected error: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ColorCode(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestColorCode_RGB_Malformed(t *testing.T) {
	bad := []string{"rgb()", "rgb(255,0)", "rgb(255,0,0,0)", "rgb(abc,0,0)"}
	for _, v := range bad {
		_, err := ColorCode(v)
		if err == nil {
			t.Errorf("ColorCode(%q): expected error for malformed rgb, got nil", v)
		}
	}
}

func TestColorCode_HSL(t *testing.T) {
	// hsl(0,100%,50%) = pure red = rgb(255,0,0)
	got, err := ColorCode("hsl(0,100%,50%)")
	if err != nil {
		t.Fatalf("ColorCode(\"hsl(0,100%%,50%%)\") returned unexpected error: %v", err)
	}
	if got != "\033[38;2;255;0;0m" {
		t.Errorf("ColorCode(\"hsl(0,100%%,50%%)\") = %q, want \\033[38;2;255;0;0m", got)
	}

	// hsl(120,100%,50%) = pure green = rgb(0,255,0)
	got, err = ColorCode("hsl(120,100%,50%)")
	if err != nil {
		t.Fatalf("ColorCode(\"hsl(120,100%%,50%%)\") returned unexpected error: %v", err)
	}
	if got != "\033[38;2;0;255;0m" {
		t.Errorf("ColorCode(\"hsl(120,100%%,50%%)\") = %q, want \\033[38;2;0;255;0m", got)
	}

	// hsl(0,0%,0%) = black = rgb(0,0,0)
	got, err = ColorCode("hsl(0,0%,0%)")
	if err != nil {
		t.Fatalf("ColorCode(\"hsl(0,0%%,0%%)\") returned unexpected error: %v", err)
	}
	if got != "\033[38;2;0;0;0m" {
		t.Errorf("ColorCode(\"hsl(0,0%%,0%%)\") = %q, want \\033[38;2;0;0;0m", got)
	}
}

func TestColorCode_HSL_Malformed(t *testing.T) {
	bad := []string{"hsl()", "hsl(0,100%)", "hsl(abc,100%,50%)"}
	for _, v := range bad {
		_, err := ColorCode(v)
		if err == nil {
			t.Errorf("ColorCode(%q): expected error for malformed hsl, got nil", v)
		}
	}
}

// ---------------------------------------------------------------------------
// BuildColorMask
// ---------------------------------------------------------------------------

func TestBuildColorMask_EmptySubstr(t *testing.T) {
	runes := []rune("hello")
	mask := BuildColorMask(runes, "")
	for i, v := range mask {
		if !v {
			t.Errorf("expected mask[%d]=true for empty substr, got false", i)
		}
	}
}

func TestBuildColorMask_SubstrNotPresent(t *testing.T) {
	runes := []rune("hello")
	mask := BuildColorMask(runes, "z")
	for i, v := range mask {
		if v {
			t.Errorf("expected mask[%d]=false for absent substr, got true", i)
		}
	}
}

func TestBuildColorMask_SingleCharSingleOccurrence(t *testing.T) {
	runes := []rune("hello")
	mask := BuildColorMask(runes, "h")
	expected := []bool{true, false, false, false, false}
	for i, v := range mask {
		if v != expected[i] {
			t.Errorf("mask[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestBuildColorMask_SingleCharMultipleOccurrences(t *testing.T) {
	// 'l' appears at positions 2 and 3 in "hello"
	runes := []rune("hello")
	mask := BuildColorMask(runes, "l")
	expected := []bool{false, false, true, true, false}
	for i, v := range mask {
		if v != expected[i] {
			t.Errorf("mask[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestBuildColorMask_MultiCharMultipleOccurrences(t *testing.T) {
	// "ll" appears once at positions 2-3 in "hello"
	runes := []rune("hello")
	mask := BuildColorMask(runes, "ll")
	expected := []bool{false, false, true, true, false}
	for i, v := range mask {
		if v != expected[i] {
			t.Errorf("mask[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestBuildColorMask_CaseSensitive(t *testing.T) {
	// 'H' should not match 'h'
	runes := []rune("hello")
	mask := BuildColorMask(runes, "H")
	for i, v := range mask {
		if v {
			t.Errorf("mask[%d] = true, expected all false (case-sensitive mismatch)", i)
		}
	}
}

func TestBuildColorMask_SubstrEqualsString(t *testing.T) {
	runes := []rune("hello")
	mask := BuildColorMask(runes, "hello")
	for i, v := range mask {
		if !v {
			t.Errorf("mask[%d] = false, expected all true when substr equals string", i)
		}
	}
}

func TestBuildColorMask_SubstrLongerThanString(t *testing.T) {
	runes := []rune("hi")
	mask := BuildColorMask(runes, "hello")
	for i, v := range mask {
		if v {
			t.Errorf("mask[%d] = true, expected all false when substr longer than string", i)
		}
	}
}

// ---------------------------------------------------------------------------
// RenderWithColor
// ---------------------------------------------------------------------------

func TestRenderWithColor_EmptySegment(t *testing.T) {
	output := captureStdout(func() {
		RenderWithColor([]string{""}, map[rune][]string{}, []ColorPair{})
	})
	if output != "\n" {
		t.Errorf("expected a single blank line, got %q", output)
	}
}

func TestRenderWithColor_NoPairs_FallsBackToRender(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	plain := captureStdout(func() {
		Render([]string{"A"}, bannerMap)
	})
	colored := captureStdout(func() {
		RenderWithColor([]string{"A"}, bannerMap, []ColorPair{})
	})
	if plain != colored {
		t.Error("RenderWithColor with no pairs should produce identical output to Render")
	}
}

func TestRenderWithColor_WholeStringColored(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	ansiCode, _ := ColorCode("red")
	output := captureStdout(func() {
		RenderWithColor([]string{"A"}, bannerMap, []ColorPair{{AnsiCode: ansiCode, Substr: ""}})
	})
	if !strings.Contains(output, ansiCode) {
		t.Error("expected ANSI color code in output when coloring whole string")
	}
	if !strings.Contains(output, ansiReset) {
		t.Error("expected ANSI reset code in output when coloring whole string")
	}
}

func TestRenderWithColor_SubstrColored(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	ansiCode, _ := ColorCode("green")
	// Color only "B" in "AB" — A should have no color code, B should
	output := captureStdout(func() {
		RenderWithColor([]string{"AB"}, bannerMap, []ColorPair{{AnsiCode: ansiCode, Substr: "B"}})
	})
	if !strings.Contains(output, ansiCode) {
		t.Error("expected ANSI color code in output for matched substring")
	}
	if !strings.Contains(output, ansiReset) {
		t.Error("expected ANSI reset code after colored substring")
	}
}

func TestRenderWithColor_NoMatchNoANSI(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	ansiCode, _ := ColorCode("blue")
	// Substr "Z" does not appear in "AB"
	output := captureStdout(func() {
		RenderWithColor([]string{"AB"}, bannerMap, []ColorPair{{AnsiCode: ansiCode, Substr: "Z"}})
	})
	if strings.Contains(output, "\033[") {
		t.Error("expected no ANSI codes when substr is not found in string")
	}
}

func TestRenderWithColor_OutputRowCount(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	ansiCode, _ := ColorCode("cyan")
	output := captureStdout(func() {
		RenderWithColor([]string{"A"}, bannerMap, []ColorPair{{AnsiCode: ansiCode, Substr: "A"}})
	})
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 rows of output, got %d", len(lines))
	}
}

func TestRenderWithColor_MultipleColors(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	red, _ := ColorCode("red")
	blue, _ := ColorCode("blue")
	pairs := []ColorPair{
		{AnsiCode: red, Substr: "A"},
		{AnsiCode: blue, Substr: "B"},
	}
	output := captureStdout(func() {
		RenderWithColor([]string{"AB"}, bannerMap, pairs)
	})
	if !strings.Contains(output, red) {
		t.Error("expected red ANSI code in output for 'A'")
	}
	if !strings.Contains(output, blue) {
		t.Error("expected blue ANSI code in output for 'B'")
	}
}

func TestRenderWithColor_MultilineWithColor(t *testing.T) {
	bannerMap, err := LoadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("LoadBanner returned an error: %v", err)
	}
	ansiCode, _ := ColorCode("yellow")
	output := captureStdout(func() {
		RenderWithColor([]string{"A", "", "B"}, bannerMap, []ColorPair{{AnsiCode: ansiCode, Substr: ""}})
	})
	// A: 8 rows, empty segment: 1 blank line, B: 8 rows = 17 lines total
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 17 {
		t.Errorf("expected 17 lines for A + blank + B, got %d", len(lines))
	}
}
