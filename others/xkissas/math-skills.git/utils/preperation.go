package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ReadConvStoreLines(file *os.File) []float64 {
	lineScan := bufio.NewScanner(file)
	var allIntLines []float64

	for lineScan.Scan() {
		strLine := lineScan.Text()
		intLine := convLinesToInt(strLine)
		if intLine == 0 {
			continue
		}
		allIntLines = append(allIntLines, intLine)
	}

	if err := lineScan.Err(); err != nil {
		log.Printf("Error scanning lines. Error message : %v", err)
	}
	file.Close()
	return allIntLines
}

func convLinesToInt(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)

	if err != nil {
		log.Printf("Cannot convert '%s' to intiger, this line will be skipped and it will not be calulated in the population of the elements.\n ---Error message : %v", s, err)
		return 0
	}
	return num
}

func DataValidation(average float64, numLines []float64) {
	if average == 0.0 && len(numLines) == 0 {
		panic("There were no elements with the correct format in the file that was provided.")
	}
}
