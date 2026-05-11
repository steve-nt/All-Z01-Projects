package core

import "strings"

func tokenizeInput(str string) []string {
	fields := strings.Fields(str)
	var ret []string
	for _, f := range fields {
		locs := re.FindAllStringIndex(f, -1)
		if locs != nil {
			idxs := make(map[int]bool)
			for _, loc := range locs {
				for _, v := range loc {
					idxs[v] = true
				}
			}
			var sb strings.Builder
			for i, r := range f {
				if idxs[i] {
					if sb.Len() > 0 {
						ret = append(ret, sb.String())
					}
					sb.Reset()
				}
				sb.WriteRune(r)
			}
			ret = append(ret, sb.String())
		} else {
			ret = append(ret, f)
		}
	}
	return ret
}
