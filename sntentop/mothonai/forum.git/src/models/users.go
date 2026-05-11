package models

import (
	"errors"
	"forum/src/utils"
)

func GetAllUsernames() ([]string, error) {
	rows, err := db.Query(`SELECT username FROM users`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return []string{}, err
	}
	defer rows.Close()
	var usernames []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []string{}, err
		}
		usernames = append(usernames, email)
	}
	return usernames, nil
}

func GetAllUserEmails() ([]string, error) {
	rows, err := db.Query(`SELECT email FROM users`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return []string{}, err
	}
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return []string{}, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}
