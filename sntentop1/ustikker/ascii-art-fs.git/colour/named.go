package asciiart

import "fmt"

// Function to convert named color to ANSI
func NamedColorToAnsi(name string) (string, error) {
	colors := map[string][3]int{
		"red":     {255, 0, 0},
		"green":   {0, 255, 0},
		"blue":    {0, 0, 255},
		"yellow":  {255, 255, 0},
		"cyan":    {0, 255, 255},
		"magenta": {255, 0, 255},
		"black":   {0, 0, 0},
		"white":   {255, 255, 255},
	}

	if rgb, found := colors[name]; found {
		return RgbToAnsi(rgb[0], rgb[1], rgb[2]), nil
	}
	return "", fmt.Errorf("unknown color name: %s", name)
}
