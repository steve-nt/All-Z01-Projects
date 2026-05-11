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

// type location []Locations

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// type date []Dates

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
