package models

import (
	"errors"
	"forum/src/utils"
)

type Categories []Category

func GetAllCategories() (Categories, error) {
	rows, err := db.Query(`SELECT id, name, description FROM categories`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return []Category{}, err
	}
	defer rows.Close()
	var categories Categories
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *Categories) IsEmpty() bool {
	if c == nil {
		return true
	}
	if c == (&Categories{}) {
		return true
	}
	return false
}

func GetCategoriesByPostId(post_id int64) (Categories, error) {
	var categories Categories
	rows, err := db.Query(`
	SELECT c.id, c.name, c.description
	FROM categories c
	JOIN posts_categories pc ON c.id = pc.category_id
	WHERE pc.post_id = ?
	`, post_id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Categories{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func AddCategory(category Category) error {
	err := category.ValidateCategory()
	if err != nil {
		return err
	}
	if (&category).DoesCategoryExist() {
		return ErrorCategoryAlreadyExists
	}
	stmt, err := db.Prepare("INSERT INTO categories (name, description) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category.Name, category.Description)
	return err
}
