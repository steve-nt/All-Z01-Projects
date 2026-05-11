package core

import "strings"

func buildOutput(tokens []string) string {
	var sb strings.Builder
	q := false
	for i := range tokens {
		sb.WriteString(tokens[i])
		if tokens[i] == "'" {
			q = !q
			if q {
				continue
			}
		}
		if i < len(tokens)-1 && !re.MatchString(tokens[i+1]) {
			sb.WriteRune(' ')
		} else if i < len(tokens)-1 && tokens[i+1] == "'" && !q {
			sb.WriteRune(' ')
		}
	}
	return sb.String()
}
