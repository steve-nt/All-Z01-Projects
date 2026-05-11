package main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	testFiles := []string{
		"samples/goodexample00.txt",
		"samples/goodexample01.txt",
		"samples/goodexample02.txt",
		"samples/goodexample03.txt",
		"samples/goodexample04.txt",
		"samples/hardexam.txt",
	}

	for _, filename := range testFiles {
		start := time.Now()

		cmd := exec.Command("go", "run", ".", filename)

		// Run the command
		err := cmd.Run()
		if err != nil {
			t.Errorf("Error running command for %s: %v", filename, err)
			continue
		}

		duration := time.Since(start)
		durationInSeconds := duration.Seconds() // Get duration in seconds
		fmt.Printf("Execution time for %s: %.6f\n", filename, durationInSeconds)
	}
}
