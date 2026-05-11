package models

type PostResponse struct {
	Username string   `json:"username"`
	Category []string `json:"category"`
	Posts    []Post   `json:"posts"`
}

type ProfileResponse struct {
	Username            string         `json:"username"`
	Notifications       []Notification `json:"notifications"`
	Activities          []Activity     `json:"activities"`
	UnseenNotifications int            `json:"unseennotifications"`
	LikedPosts          []Post         `json:"liked"`
	CreatedPosts        []Post         `json:"created"`
}
type PostByIdResponse struct {
	Username string
	Post     Post
}
