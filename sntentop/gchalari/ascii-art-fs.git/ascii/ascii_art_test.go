package ascii

import "testing"

func TestRenderOnlyNewline(t *testing.T) {
	banner := map[rune][charHeight]string{
		'a': {"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7"},
	}

	got := Render("\\n", banner)
	want := "\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderSimpleCharacter(t *testing.T) {
	banner := map[rune][charHeight]string{
		'a': {"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7"},
	}

	got := Render("a", banner)
	want := "a0\n" +
		"a1\n" +
		"a2\n" +
		"a3\n" +
		"a4\n" +
		"a5\n" +
		"a6\n" +
		"a7\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderTwoLines(t *testing.T) {
	banner := map[rune][charHeight]string{
		'a': {"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7"},
		'b': {"b0", "b1", "b2", "b3", "b4", "b5", "b6", "b7"},
	}

	got := Render("a\\nb", banner)

	want := "a0\n" +
		"a1\n" +
		"a2\n" +
		"a3\n" +
		"a4\n" +
		"a5\n" +
		"a6\n" +
		"a7\n" +
		"b0\n" +
		"b1\n" +
		"b2\n" +
		"b3\n" +
		"b4\n" +
		"b5\n" +
		"b6\n" +
		"b7\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
