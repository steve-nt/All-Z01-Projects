package utils

import "forum-image-upload/internal/backend/models"

func SortActivities(acties []models.Activity) (result []models.Activity) {

	for i := 0; i <= len(acties)-2; i++ {
		for j := i + 1; j <= len(acties)-1; j++ {
			if acties[j].CreationDate > acties[i].CreationDate {
				temp := acties[i]
				acties[i] = acties[j]
				acties[j] = temp
			}
		}
	}

	return acties
}
