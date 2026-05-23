package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ansiReset = "\033[0m"

// ColorCode resolves a color name or notation to an ANSI foreground escape sequence.
// Supported formats: named colors, #rrggbb, rgb(r,g,b), hsl(h,s%,l%).
// Returns an error if the format is not recognized.
func ColorCode(name string) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(name))

	// --- Named colors ---
	switch lower {
	case "black":
		return "\033[30m", nil
	case "red":
		return "\033[31m", nil
	case "green":
		return "\033[32m", nil
	case "yellow":
		return "\033[33m", nil
	case "blue":
		return "\033[34m", nil
	case "magenta", "purple":
		return "\033[35m", nil
	case "cyan":
		return "\033[36m", nil
	case "white":
		return "\033[37m", nil
	case "orange":
		return "\033[93m", nil // bright yellow — closest standard ANSI to orange
	case "pink":
		return "\033[95m", nil // bright magenta
	case "brightblack", "darkgray", "dark gray":
		return "\033[90m", nil
	case "brightred":
		return "\033[91m", nil
	case "brightgreen":
		return "\033[92m", nil
	case "brightyellow":
		return "\033[93m", nil
	case "brightblue":
		return "\033[94m", nil
	case "brightmagenta":
		return "\033[95m", nil
	case "brightcyan":
		return "\033[96m", nil
	case "brightwhite":
		return "\033[97m", nil
	}

	// --- Hex: #rrggbb ---
	if strings.HasPrefix(lower, "#") {
		hex := lower[1:]
		if len(hex) != 6 {
			return "", fmt.Errorf("invalid hex color: %q (must be #rrggbb)", name)
		}
		r, err1 := strconv.ParseUint(hex[0:2], 16, 8)
		g, err2 := strconv.ParseUint(hex[2:4], 16, 8)
		b, err3 := strconv.ParseUint(hex[4:6], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return "", fmt.Errorf("invalid hex color: %q", name)
		}
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
	}

	// --- RGB: rgb(r,g,b) ---
	if strings.HasPrefix(lower, "rgb(") && strings.HasSuffix(lower, ")") {
		inner := lower[4 : len(lower)-1] // strip "rgb(" and ")"
		parts := strings.Split(inner, ",")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid rgb color: %q", name)
		}
		r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err1 != nil || err2 != nil || err3 != nil {
			return "", fmt.Errorf("invalid rgb color: %q", name)
		}
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
	}

	// --- HSL: hsl(h,s%,l%) ---
	if strings.HasPrefix(lower, "hsl(") && strings.HasSuffix(lower, ")") {
		inner := lower[4 : len(lower)-1] // strip "hsl(" and ")"
		parts := strings.Split(inner, ",")
		if len(parts) != 3 {
			return "", fmt.Errorf("invalid hsl color: %q", name)
		}
		h, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		s, err2 := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"), 64)
		l, err3 := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"), 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return "", fmt.Errorf("invalid hsl color: %q", name)
		}
		r, g, b := hslToRGB(h, s/100, l/100)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
	}

	return "", fmt.Errorf("unknown color: %q", name)
}

// hslToRGB converts HSL values (H: 0–360, S: 0–1, L: 0–1) to RGB (0–255 each).
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		v := uint8(l * 255)
		return v, v, v
	}

	hueToChannel := func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		switch {
		case t < 1.0/6:
			return p + (q-p)*6*t
		case t < 1.0/2:
			return q
		case t < 2.0/3:
			return p + (q-p)*(2.0/3-t)*6
		default:
			return p
		}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hNorm := h / 360

	r := hueToChannel(p, q, hNorm+1.0/3)
	g := hueToChannel(p, q, hNorm)
	b := hueToChannel(p, q, hNorm-1.0/3)

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}
