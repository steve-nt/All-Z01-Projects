package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

// downloadFile downloads a single URL and saves it to the configured path/name.
// Output is written to the provided writer (stdout or log file).
func downloadFile(url, outputName, outputPath string, rateLimit int64, w io.Writer) error {
	fmt.Fprintf(w, "start at %s\n", time.Now().Format(timeFormat))

	fmt.Fprintf(w, "sending request, awaiting response... ")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "\n")
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "status %s\n", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	contentLength := resp.ContentLength
	if contentLength > 0 {
		fmt.Fprintf(w, "content size: %d [~%.2fMB]\n", contentLength, float64(contentLength)/(1024*1024))
	} else {
		fmt.Fprintf(w, "content size: unknown\n")
	}

	// Determine output filename
	if outputName == "" {
		outputName = extractFileName(url)
	}

	// Determine output path
	savePath := outputName
	if outputPath != "" {
		// Expand ~ if present
		if strings.HasPrefix(outputPath, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				outputPath = filepath.Join(home, outputPath[2:])
			}
		}
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", outputPath, err)
		}
		savePath = filepath.Join(outputPath, outputName)
	}

	displayPath := savePath
	if !filepath.IsAbs(savePath) {
		displayPath = "./" + savePath
	}
	fmt.Fprintf(w, "saving file to: %s\n", displayPath)

	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", savePath, err)
	}
	defer outFile.Close()

	// Show progress bar only when writing to stdout (not background log)
	showProgress := isStdout(w)

	var reader io.Reader = resp.Body
	if rateLimit > 0 {
		reader = &rateLimitedReader{r: resp.Body, rateLimit: rateLimit}
	}

	if showProgress {
		err = copyWithProgress(outFile, reader, contentLength, w)
	} else {
		_, err = io.Copy(outFile, reader)
	}
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	fmt.Fprintf(w, "\nDownloaded [%s]\n", url)
	fmt.Fprintf(w, "finished at %s\n", time.Now().Format(timeFormat))

	return nil
}

func extractFileName(url string) string {
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	// Remove query string
	if idx := strings.Index(name, "?"); idx != -1 {
		name = name[:idx]
	}
	if name == "" {
		name = "index.html"
	}
	return name
}

func isStdout(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return f == os.Stdout
	}
	return false
}

// copyWithProgress copies from src to dst while displaying a progress bar.
func copyWithProgress(dst io.Writer, src io.Reader, totalSize int64, out io.Writer) error {
	buf := make([]byte, 32*1024)
	var downloaded int64
	startTime := time.Now()

	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, werr := dst.Write(buf[:n])
			if werr != nil {
				return werr
			}
			downloaded += int64(n)
			printProgress(downloaded, totalSize, startTime, out)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func printProgress(downloaded, total int64, start time.Time, w io.Writer) {
	elapsed := time.Since(start).Seconds()
	speed := float64(downloaded) / elapsed // bytes/sec

	var sizeStr string
	if total > 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f MiB / %.2f MiB", float64(downloaded)/(1024*1024), float64(total)/(1024*1024))
	} else {
		sizeStr = fmt.Sprintf("%.2f KiB / %.2f KiB", float64(downloaded)/1024, float64(total)/1024)
	}

	var pct float64
	var bar string
	var eta string

	if total > 0 {
		pct = float64(downloaded) / float64(total) * 100
		barLen := 50
		filled := int(float64(barLen) * float64(downloaded) / float64(total))
		bar = "[" + strings.Repeat("=", filled) + strings.Repeat(" ", barLen-filled) + "]"

		remaining := float64(total-downloaded) / speed
		if speed > 0 {
			eta = formatDuration(remaining)
		} else {
			eta = "?"
		}
	} else {
		bar = ""
		pct = 0
		eta = "?"
	}

	var speedStr string
	if speed > 1024*1024 {
		speedStr = fmt.Sprintf("%.2f MiB/s", speed/(1024*1024))
	} else {
		speedStr = fmt.Sprintf("%.2f KiB/s", speed/1024)
	}

	fmt.Fprintf(w, "\r %s %s %.2f%% %s %s", sizeStr, bar, pct, speedStr, eta)
}

func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}
	mins := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm%ds", mins, secs)
}
