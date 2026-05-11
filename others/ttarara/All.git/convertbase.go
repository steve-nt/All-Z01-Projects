package piscine

func ConvertBase(nbr, baseFrom, baseTo string) string {
	numBase10 := fromBaseToDecimal(nbr, baseFrom)
	return fromDecimalToBase(numBase10, baseTo)
}

func fromBaseToDecimal(nbr string, base string) int {
	baseLen, num := len(base), 0
	for i := 0; i < len(nbr); i++ {
		value := indexOf(base, nbr[i])
		num = num*baseLen + value
	}
	return num
}

func indexOf(str string, char byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == char {
			return i
		}
	}
	return -1
}

func fromDecimalToBase(num int, base string) string {
	if num == 0 {
		return string(base[0])
	}
	baseLen, result := len(base), ""
	for num > 0 {
		result = string(base[num%baseLen]) + result
		num /= baseLen
	}
	return result
}
