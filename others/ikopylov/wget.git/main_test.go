package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNoOpBinary_ExitZero_NoOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("binary name differs on windows; not needed for this repo")
	}

	tmp := t.TempDir()
	binPath := filepath.Join(tmp, "wget")

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Env = os.Environ()
	out, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, string(out))
	}

	run := exec.Command(binPath)
	var stdout, stderr bytes.Buffer
	run.Stdout = &stdout
	run.Stderr = &stderr

	if err := run.Run(); err != nil {
		t.Fatalf("binary run failed: %v", err)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout; got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr; got %q", stderr.String())
	}
}

