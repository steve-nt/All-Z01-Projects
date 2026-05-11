package models

type Artist struct {
	Id                int      `json:"id"`
	Image             string   `json:"image"`
	Name              string   `json:"name"`
	Members           []string `json:"members"`
	MembersFormatted  []string
	CreationDate      int                 `json:"creationDate"`
	FirstAlbum        string              `json:"firstAlbum"`
	Concerts          map[string][]string `json:"datesLocations"`
	ConcertsFormatted []string
}
