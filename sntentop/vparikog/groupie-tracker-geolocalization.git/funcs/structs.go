package groupie_tracker_search

type All_API struct {
	Artists_API      string `json:"artists"`
	Locations_API    string `json:"locations"`
	ConcertDates_API string `json:"dates"`
	Relations_API    string `json:"relation"`
}

type Artist struct {
	ID             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	NumMembers     int      `json:"numMembers"`
	NameMembers    string   `json:"NameMembers"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Locations      `json:"-"`
	Dates          `json:"-"`
	Relations      `json:"-"`
	GeoCoordinates map[string]GeoResult `json:"geoCoordinates"`
}

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationsResponse struct {
	Index []Locations `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesResponse struct {
	Index []Dates `json:"index"`
}

type Relations struct {
	ID             int                 `json:"id"`
	Dates_Location map[string][]string `json:"datesLocations"`
}

type RelationsResponse struct {
	Index []Relations `json:"index"`
}

type GeoResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}
