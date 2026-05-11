package utils

import "forum-authentication/internal/backend/models"

func ReversePosts(posts *[]models.Post) {
	s := *posts
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	*posts = s
}
