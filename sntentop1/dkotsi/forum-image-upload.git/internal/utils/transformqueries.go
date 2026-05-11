package utils

import (
	"strings"
)

func DesignQueryBasedOnCategories(query string, categories []string) (resultquery string) {

	for i, category := range categories {
		categories[i] = "name==" + "'" + category + "'"
	}
	toaddinquery := strings.Join(categories, " or ")
	query += toaddinquery

	return query
}
