package main

import (
	"os"

	"forum/src/models"
	"forum/src/utils"
)

// mockGen generates mock data for the database
func mockGen(dbPath string) error {
	if err := models.InitDB(dbPath); err != nil {
		return err
	}

	// Create some mock categories
	categories := []models.Category{
		{Name: "General", Description: "General discussions"},
		{Name: "Tech", Description: "Technology related discussions"},
		{Name: "Random", Description: "Random topics"},
	}

	for _, cat := range categories {
		if !cat.DoesCategoryExist() {
			if err := models.AddCategory(cat); err != nil {
				utils.LogDebug(err)
			}
		}
	}

	return nil
}

func main() {
	db_path := "./db.db"
	if len(os.Args) == 2 {
		db_path = os.Args[1]
	}
	mockGen(db_path)
}
