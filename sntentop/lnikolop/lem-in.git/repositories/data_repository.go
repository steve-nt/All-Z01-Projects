package repositories

import (
	"lem-in/models"
)

type DataRepository interface {
	ReadFile(filename string) ([]string, error)
	FetchData(lines []string) (int, []*models.Room, error)
}
