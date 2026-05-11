package main

import (
	"fmt"
	"strings"
)

type Config struct {
	URLs        []string
	OutputName  string
	OutputPath  string
	RateLimit   int64 // bytes per second, 0 = unlimited
	Background  bool
	InputFile   string
	Mirror      bool
	RejectTypes []string
	ExcludeDirs []string
	ConvertLinks bool
}

func parseArgs(args []string) (*Config, error) {
	config := &Config{}

	var urls []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-B":
			config.Background = true

		case arg == "--mirror":
			config.Mirror = true

		case arg == "--convert-links":
			config.ConvertLinks = true

		case strings.HasPrefix(arg, "-O="):
			config.OutputName = strings.TrimPrefix(arg, "-O=")

		case arg == "-O":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-O requires a filename argument")
			}
			i++
			config.OutputName = args[i]

		case strings.HasPrefix(arg, "-P="):
			config.OutputPath = strings.TrimPrefix(arg, "-P=")

		case arg == "-P":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-P requires a path argument")
			}
			i++
			config.OutputPath = args[i]

		case strings.HasPrefix(arg, "--rate-limit="):
			val := strings.TrimPrefix(arg, "--rate-limit=")
			rate, err := parseRateLimit(val)
			if err != nil {
				return nil, fmt.Errorf("invalid rate-limit: %v", err)
			}
			config.RateLimit = rate

		case strings.HasPrefix(arg, "-i="):
			config.InputFile = strings.TrimPrefix(arg, "-i=")

		case arg == "-i":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-i requires a filename argument")
			}
			i++
			config.InputFile = args[i]

		case strings.HasPrefix(arg, "-R="):
			val := strings.TrimPrefix(arg, "-R=")
			config.RejectTypes = strings.Split(val, ",")

		case arg == "-R":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-R requires a file extension list")
			}
			i++
			config.RejectTypes = strings.Split(args[i], ",")

		case strings.HasPrefix(arg, "--reject="):
			val := strings.TrimPrefix(arg, "--reject=")
			config.RejectTypes = strings.Split(val, ",")

		case strings.HasPrefix(arg, "-X="):
			val := strings.TrimPrefix(arg, "-X=")
			config.ExcludeDirs = strings.Split(val, ",")

		case arg == "-X":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("-X requires a directory list")
			}
			i++
			config.ExcludeDirs = strings.Split(args[i], ",")

		case strings.HasPrefix(arg, "--exclude="):
			val := strings.TrimPrefix(arg, "--exclude=")
			config.ExcludeDirs = strings.Split(val, ",")

		case strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") || strings.HasPrefix(arg, "ftp://"):
			urls = append(urls, arg)

		default:
			return nil, fmt.Errorf("unknown argument: %s", arg)
		}
	}

	config.URLs = urls

	// Validate
	if config.InputFile == "" && !config.Mirror && len(config.URLs) == 0 {
		return nil, fmt.Errorf("no URL provided")
	}
	if config.Mirror && len(config.URLs) == 0 {
		return nil, fmt.Errorf("--mirror requires a URL")
	}

	return config, nil
}

func parseRateLimit(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0, fmt.Errorf("empty rate limit")
	}

	multiplier := int64(1)
	switch s[len(s)-1] {
	case 'k', 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	var val int64
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return val * multiplier, nil
}
