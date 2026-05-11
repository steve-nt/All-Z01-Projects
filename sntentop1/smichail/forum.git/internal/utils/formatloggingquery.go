package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func FormatLoggingQuery(q string, args ...any) (result string) {
	result = q

	for _, arg := range args {
		switch arg.(type) {
		case string:
			result = strings.Replace(result, "?", fmt.Sprintf("'%v'", arg), 1)
		case uuid.UUID:
			result = strings.Replace(result, "?", fmt.Sprintf("'%v'", arg), 1)
		case time.Time:
			result = strings.Replace(result, "?", fmt.Sprintf("'%v'", arg), 1)
		default:
			result = strings.Replace(result, "?", fmt.Sprintf("%v", arg), 1)
		}
	}
	return result + "\n"
}
