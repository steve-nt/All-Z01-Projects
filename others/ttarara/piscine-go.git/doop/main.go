package main

import (
	"os"
)

const (
	maxInt = int64(^uint64(0) >> 1)
	minInt = -maxInt - 1
)

func main() {
	if len(os.Args) != 4 {
		return
	}

	arg1 := os.Args[1]
	op := os.Args[2]
	arg2 := os.Args[3]

	firstNbr, err1 := Atoi(arg1)
	secondNbr, err2 := Atoi(arg2)
	if err1 != nil || err2 != nil {
		return
	}

	var result string

	switch op {
	case "+":
		if (secondNbr > 0 && firstNbr > maxInt-secondNbr) || (secondNbr < 0 && firstNbr < minInt-secondNbr) {
			return
		}
		result = Itoa(firstNbr + secondNbr)
	case "-":
		if (secondNbr < 0 && firstNbr > maxInt+secondNbr) || (secondNbr > 0 && firstNbr < minInt+secondNbr) {
			return
		}
		result = Itoa(firstNbr - secondNbr)
	case "*":
		if firstNbr != 0 && secondNbr != 0 {
			if firstNbr > maxInt/secondNbr || firstNbr < minInt/secondNbr {
				return
			}
		}
		result = Itoa(firstNbr * secondNbr)
	case "/":
		if secondNbr == 0 {
			result = "No division by 0"
		} else {
			result = Itoa(firstNbr / secondNbr)
		}
	case "%":
		if secondNbr == 0 {
			result = "No modulo by 0"
		} else {
			result = Itoa(firstNbr % secondNbr)
		}
	default:
		return
	}

	os.Stdout.WriteString(result + "\n")
}

// Atoi converts a string to an integer.
func Atoi(s string) (int64, error) {
	result := int64(0)
	sign := int64(1)
	start := 0

	if len(s) == 0 {
		return 0, os.ErrInvalid
	}

	if s[0] == '-' {
		sign = -1
		start = 1
	} else if s[0] == '+' {
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, os.ErrInvalid
		}
		result = result*10 + int64(s[i]-'0')
		if result < 0 { // overflow detection
			return 0, os.ErrInvalid
		}
	}

	return result * sign, nil
}

// Itoa converts an integer to a string.
func Itoa(n int64) string {
	if n == 0 {
		return "0"
	}

	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}

	digits := [20]byte{}
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte(n%10) + '0'
		n /= 10
	}

	return sign + string(digits[i:])
}
