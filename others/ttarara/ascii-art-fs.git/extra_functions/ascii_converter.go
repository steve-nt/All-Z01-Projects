package extra_functions

// Return2DASCIIArray converts lines read from the file into a 2D array of ASCII templates.
func Return2DASCIIArray(fileLines []string) [][]string {
	var asciiTemplates [][]string
	counter := 0
	var tempAsciArray []string

	for _, line := range fileLines {
		counter++
		if counter != 1 { // Skip the first line if unnecessary
			tempAsciArray = append(tempAsciArray, line)
		}
		if counter == 9 { // Each character is represented by 9 lines
			asciiTemplates = append(asciiTemplates, tempAsciArray)
			counter = 0
			tempAsciArray = nil
		}
	}
	return asciiTemplates
}

// ReturnAsciiCodeInt converts each character to its ASCII template index.
func ReturnAsciiCodeInt(s string) []int {
	var tempIntArrLetter []int
	for _, v := range s {
		tempIntArrLetter = append(tempIntArrLetter, int(v)-32)
	}
	return tempIntArrLetter
}
