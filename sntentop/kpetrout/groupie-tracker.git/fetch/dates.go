package fetch

type Date struct {
	Dates []string `json:"dates"`
}

func DatesParse(relations map[string][]string) []string {
	var dates []string
	for _, relation := range relations {
		dates = append(dates, relation...)
	}
	return dates
}
