package core

import "strings"

func executeCommands(tokens []string) []string {
	var ret []string
	for i := 0; i < len(tokens); i++ {
		w := tokens[i]
		target := strings.TrimSuffix(w, ")")
		f := cmdMap[target]
		lowW := strings.ToLower(w)
		switch lowW {
		case "a", "an":
			next := false
			j := 1
			for !next && i+j < len(tokens) {
				r := rune(strings.ToLower(tokens[i+j])[0])
				if r >= 'a' && r <= 'z' {
					switch r {
					case 'a', 'e', 'i', 'o', 'u', 'h':
						if lowW == "a" {
							ret = append(ret, w+"n")
						} else {
							ret = append(ret, w)
						}
					default:
						if lowW == "a" {
							ret = append(ret, w)
						} else {
							ret = append(ret, w[:len(w)-1])
						}
					}
					next = true
				} else {
					j++
				}
			}
			if !next {
				ret = append(ret, w)
			}
		case "(up)", "(low)", "(cap)", "(bin)", "(hex)":
			count := 0
			j := 1
			for count < 1 && len(ret)-j >= 0 {
				if !re.MatchString(ret[len(ret)-j]) {
					ret[len(ret)-j] = f(ret[len(ret)-j])
					count++
				}
				j++
			}
		case "(up", "(cap", "(low":
			if i+2 < len(tokens) {
				n := getN(tokens[i+2])
				count := 0
				j := 1
				for count < n && len(ret)-j >= 0 {
					if !re.MatchString(ret[len(ret)-j]) {
						ret[len(ret)-j] = f(ret[len(ret)-j])
						count++
					}
					j++
				}
				i += 2
			}
		default:
			ret = append(ret, w)
		}
	}
	return ret
}
