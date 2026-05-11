package main

import (
	"fmt"
	"os"
)

// run executes the appropriate download strategy based on config.
func run(config *Config) error {
	// Mirror mode
	if config.Mirror {
		return mirrorSite(config.URLs[0], config.RejectTypes, config.ExcludeDirs, config.ConvertLinks, config.RateLimit, os.Stdout)
	}

	// Multi-file mode
	if config.InputFile != "" {
		return downloadMultiple(config.InputFile, config.OutputPath, config.RateLimit, os.Stdout)
	}

	// Background mode
	if config.Background {
		if len(config.URLs) == 0 {
			return fmt.Errorf("no URL provided")
		}
		return downloadBackground()
	}

	// Single file mode
	if len(config.URLs) == 0 {
		return fmt.Errorf("no URL provided")
	}
	return downloadFile(config.URLs[0], config.OutputName, config.OutputPath, config.RateLimit, os.Stdout)
}
