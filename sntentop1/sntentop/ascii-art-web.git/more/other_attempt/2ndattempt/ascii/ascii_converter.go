// ascii/ascii_converter.go
package ascii

func ReturnAsciiCodeInt(s string) []int {
	var asciiArr []int
	for _, v := range s {
		asciiArr = append(asciiArr, int(v)-32)
	}
	return asciiArr
}
