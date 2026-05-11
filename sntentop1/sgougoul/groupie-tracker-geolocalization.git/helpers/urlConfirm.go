package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

// function to check invalid url query
func CheckUrl(s string) (bool, int) {
	reg := regexp.MustCompile(`^\d+$`)
	urlArray := strings.Split(s, "/")
	if len(urlArray) < 3 || len(urlArray) > 4 {
		return false, 0
	}
	if len(urlArray) >= 3 {
		isValid := reg.MatchString(urlArray[3])
		if !isValid {
			return false, 0
		}
	}

	idStr := urlArray[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return false, 0
	}

	return true, id
}
