package converts

import (
	"fmt"
	"strconv"
)

func processHexCommand(result []string) []string {
	if len(result) > 0 {
		decimalValue, err := strconv.ParseInt(result[len(result)-1], 16, 64)
		if err != nil {
			fmt.Println("Error Parsing Hexadecimal to Decimal:", err)
		} else {
			result[len(result)-1] = strconv.Itoa(int(decimalValue))
		}
	}
	return result
}
func processBinCommand(result []string) []string {
	if len(result) > 0 {
		decimalValue, err := strconv.ParseInt(result[len(result)-1], 2, 64)
		if err != nil {
			fmt.Println("Error Parsing Binary to Decimal:", err)
		} else {
			result[len(result)-1] = strconv.Itoa(int(decimalValue))
		}
	}
	return result
}
