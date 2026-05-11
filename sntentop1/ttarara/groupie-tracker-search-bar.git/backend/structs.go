package backend

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
	Relation     Relation
	Locations    string `json:"locations"`
	Location     Locations
	Dates        string `json:"concertDates"`
	Date         Dates
}

type Artists []Artist

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Filters struct {
	OptionsLocations []string
	OptionsDates     []string
}

type SearchResult struct {
	Type     string `json:"type"` // e.g., "artist", "member", "location"
	Display  string `json:"display"`
	ArtistID int    `json:"artistId"` // ID of the associated artist
	Location string `json:"location"` // Location for "location-all" type
}
