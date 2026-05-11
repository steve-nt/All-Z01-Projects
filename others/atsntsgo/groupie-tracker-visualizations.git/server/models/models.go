package models

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationData struct {
	Index []Location `json:"index"`
}

type DateData struct {
	Index []Date `json:"index"`
}

type RelationData struct {
	Index []Relation `json:"index"`
}

// Indexes holds all the data needed for the template
type Indexes struct {
	Artists   []Artist   // Artist data
	Locations []Location // Location data from the index field
	Dates     []Date     // Date data from the index field
	Relations []Relation // Relation data from the index field
}
