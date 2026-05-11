package banner

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Add regex patterns for hex, RGB, and HSL formats
var hexPattern = regexp.MustCompile(`^#?([a-fA-F0-9]{6})$`)
var rgbPattern = regexp.MustCompile(`^rgb\((\d{1,3}),\s*(\d{1,3}),\s*(\d{1,3})\)$`)
var hslPattern = regexp.MustCompile(`^hsl\((\d{1,3}),\s*(\d{1,3})%,\s*(\d{1,3})%\)$`)

func SetColor(colorName string) string {
	switch {
	case colorName == "red":
		return "\033[31m"
	case colorName == "green":
		return "\033[32m"
	case colorName == "yellow":
		return "\033[33m"
	case colorName == "blue":
		return "\033[34m"
	case colorName == "magenta":
		return "\033[35m"
	case colorName == "cyan":
		return "\033[36m"
	case colorName == "gray":
		return "\033[37m"
	case colorName == "white":
		return "\033[97m"
	case colorName == "orange":
		return "\033[38;5;208m"

	// Hex color format #RRGGBB
	case hexPattern.MatchString(colorName):
		r, g, b := hexToRGB(colorName)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

	// RGB color format rgb(r, g, b)
	case rgbPattern.MatchString(colorName):
		r, g, b := parseRGB(colorName)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

	// HSL color format hsl(h, s%, l%)
	case hslPattern.MatchString(colorName):
		r, g, b := hslToRGB(colorName)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

	default:
		return "\033[0m"
	}
}

func hexToRGB(hex string) (int, int, int) {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return int(r), int(g), int(b)
}

func parseRGB(rgb string) (int, int, int) {
	matches := rgbPattern.FindStringSubmatch(rgb)
	r, _ := strconv.Atoi(matches[1])
	g, _ := strconv.Atoi(matches[2])
	b, _ := strconv.Atoi(matches[3])
	return r, g, b
}

func hslToRGB(hsl string) (int, int, int) {
	matches := hslPattern.FindStringSubmatch(hsl)
	h, _ := strconv.Atoi(matches[1])
	s, _ := strconv.Atoi(matches[2])
	l, _ := strconv.Atoi(matches[3])
	return hslToRgbConversion(h, s, l)
}

func hslToRgbConversion(h, s, l int) (int, int, int) {
	hue := float64(h) / 360.0
	saturation := float64(s) / 100.0
	lightness := float64(l) / 100.0

	var r, g, b float64
	if saturation == 0 {
		r = lightness
		g = lightness
		b = lightness
	} else {
		var q float64
		if lightness < 0.5 {
			q = lightness * (1 + saturation)
		} else {
			q = lightness + saturation - lightness*saturation
		}
		p := 2*lightness - q
		r = hueToRGB(p, q, hue+1.0/3.0)
		g = hueToRGB(p, q, hue)
		b = hueToRGB(p, q, hue-1.0/3.0)
	}

	return int(r * 255), int(g * 255), int(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func ReadBannerFiles(txt string) map[int][]string {
	DATA := make(map[int][]string)
	file, err := os.Open(txt)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	read := bufio.NewScanner(file)
	ascciStart := 31
	for read.Scan() {
		if read.Text() == "" {
			ascciStart++
		} else {
			DATA[ascciStart] = append(DATA[ascciStart], read.Text())
		}
	}
	return DATA
}

func CheckIfAllCharInFile(words []string) bool {
	Temp := strings.Join(words, "")
	for _, char := range Temp {
		if char < ' ' || char > '~' {
			return false
		}
	}
	return true
}

func PointsToBeColored(word, LettersToBeColored string) [][]int {
	size_Word := len(word)
	Points := make([][]int, 0)
	sizeOfLetters := len(LettersToBeColored)
	for i := 0; i <= size_Word-1; i++ {
		if i <= size_Word-sizeOfLetters && word[i:sizeOfLetters+i] == LettersToBeColored {
			Points = append(Points, []int{i, i + sizeOfLetters - 1})
		}
	}
	return Points
}

func PrintChars(word, LettersColored, color string, banner map[int][]string) {
	Points := PointsToBeColored(word, LettersColored)
	Size_Points := len(Points)
	activeIndex := 0
	wordSize := len(word)
	for i := 0; i < 8; i++ {
		for j := 0; j <= wordSize-1; j++ {
			if Size_Points > activeIndex && j >= Points[activeIndex][0] && j <= Points[activeIndex][1] {
				print(SetColor(color))
				if j == Points[activeIndex][1] {
					activeIndex++
				}
			} else {
				print(SetColor("reset"))
			}
			if LettersColored == "" {
				print(SetColor(color))
			}
			if j == len(word)-1 {
				fmt.Println(banner[int(word[j])][i])
				continue
			}
			fmt.Print(banner[int(word[j])][i])
		}
		activeIndex = 0
	}
}

func Result(words []string, newLineCounter int, banner map[int][]string, color, Letters string) {
	counter := 1
	for _, word := range words {
		if word == "" && counter <= newLineCounter {
			fmt.Println()
			counter++
			continue
		}
		PrintChars(word, Letters, color, banner)
	}
	print(SetColor("reset"))
}
