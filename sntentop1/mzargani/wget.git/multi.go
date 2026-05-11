package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// downloadMultiple reads URLs from a file and downloads them all concurrently.
func downloadMultiple(inputFile string, outputPath string, rateLimit int64, w io.Writer) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found in %s", inputFile)
	}

	// Get content sizes
	var sizes []int64
	for _, url := range urls {
		resp, err := http.Head(url)
		if err == nil {
			sizes = append(sizes, resp.ContentLength)
			resp.Body.Close()
		} else {
			sizes = append(sizes, -1)
		}
	}

	// Print content sizes
	sizeStrs := make([]string, len(sizes))
	for i, s := range sizes {
		if s > 0 {
			sizeStrs[i] = fmt.Sprintf("%d", s)
		} else {
			sizeStrs[i] = "unknown"
		}
	}
	fmt.Fprintf(w, "content size: [%s]\n", strings.Join(sizeStrs, ", "))

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			name := extractFileName(u)
			savePath := name
			if outputPath != "" {
				if strings.HasPrefix(outputPath, "~/") {
					home, err := os.UserHomeDir()
					if err == nil {
						outputPath = filepath.Join(home, outputPath[2:])
					}
				}
				os.MkdirAll(outputPath, 0755)
				savePath = filepath.Join(outputPath, name)
			}

			err := downloadSimple(u, savePath, rateLimit)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fmt.Fprintf(w, "failed %s: %v\n", name, err)
			} else {
				fmt.Fprintf(w, "finished %s\n", name)
			}
		}(url)
	}

	wg.Wait()

	fmt.Fprintf(w, "\nDownload finished: [%s]\n", strings.Join(urls, " "))
	return nil
}

// downloadSimple downloads a URL to a file without progress output (used for concurrent downloads).
func downloadSimple(url, savePath string, rateLimit int64) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var reader io.Reader = resp.Body
	if rateLimit > 0 {
		reader = &rateLimitedReader{r: resp.Body, rateLimit: rateLimit}
	}

	_, err = io.Copy(outFile, reader)
	return err
}

// downloadBackground spawns a new process without -B and redirects its output to wget-log.
func downloadBackground() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable: %v", err)
	}

	var newArgs []string
	for _, arg := range os.Args[1:] {
		if arg != "-B" {
			newArgs = append(newArgs, arg)
		}
	}

	logFile, err := os.Create("wget-log")
	if err != nil {
		return fmt.Errorf("failed to create wget-log: %v", err)
	}

	cmd := exec.Command(exe, newArgs...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start background download: %v", err)
	}

	logFile.Close()
	fmt.Println(`Output will be written to "wget-log".`)
	return nil
}
