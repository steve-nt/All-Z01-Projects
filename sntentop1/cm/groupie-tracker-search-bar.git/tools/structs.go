package tools

type Artist struct {
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	Image        string              `json:"image"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Relations    map[string][]string `json:"-"`
	Locations    []string            `json:"-"`
}

type RelationData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
