package args

import "fmt"

func ValidateArgs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(`no arguments provided. Please provide at least one
Usage: go run . [STRING] [BANNER]
EX: go run . something standard`)
	}

	if len(args) > 2 {
		return fmt.Errorf(`more than 2 arguments provided. Please provide two at most.
Usage: go run . [STRING] [BANNER]
EX: go run . something standard`)
	}

	return nil
}

func ValidateBanner(banner string) error {
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		return fmt.Errorf("invalid banner. Supported banners are 'standard', 'shadow', and 'thinkertoy'")
	}

	return nil
}
