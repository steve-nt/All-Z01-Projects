package backend

// Defines a struct named `Artist`. A struct is a composite data type that groups together variables (fields) under a single name.
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

// Defines a type `Artists` as a slice (dynamic array) of `Artist` structs.
type Artists []Artist

// Defines a struct named `Locations` to store location-related data.
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

// Defines a struct named `Dates` to store concert date-related data.
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Defines a struct named `Relation` to store relationship data (e.g., concert dates and locations).
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Defines a struct named `Filters` to store filter options for locations and dates.
type Filters struct {
	OptionsLocations []string
	OptionsDates     []string
}

// Defines a struct named `SearchResult` to store search result data.
type SearchResult struct {
	Type     string `json:"type"` // e.g., "artist", "member", "location"
	Display  string `json:"display"`
	ArtistID int    `json:"artistId"` // ID of the associated artist
	Location string `json:"location"` // Location for "location-all" type
}
