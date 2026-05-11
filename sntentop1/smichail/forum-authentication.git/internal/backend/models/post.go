package models

type Post struct {
	ID               string
	Title            string
	Content          string
	User_UUID        string
	CreationDate     string
	Categories       []string
	Liked            bool
	Disliked         bool
	NumberOfComments int
	NumberOfLikes    int
	NumberOfDislikes int
	Comments         []Comment
	ImagePath        string
	// BestComment      Comment
	//The Best Comment atribute can be implemented in order to show in the home page under the post the best comment for each post
}

type Comment struct {
	ID               string
	Content          string
	User_UUID        string
	Post_id          string
	CreationDate     string
	NumberOfLikes    int
	NumberOfDislikes int
	Liked            bool
	Disliked         bool
}

type PostInfo struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
	RemoveImg  string   `json:"remove-img"`
	ImagePath  string   `json:"imagepath"`
}
