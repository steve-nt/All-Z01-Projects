package models

type SocialUser struct {
	ID             int
	UUID           string
	Provider       string
	ProviderUserID string
}
