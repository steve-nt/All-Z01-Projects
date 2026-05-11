package api

import (
	"groupie/utils"
)

// holds relational information for artists
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type relationResponse struct {
	Relations []Relation `json:"index"`
}

// Holds artist information
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// global variables to be loaded once when the server starts
var (
	Artists           []Artist
	RelationData      []Relation
	ArtistRelationMap = make(map[int]Relation)
	ArtistMap         = make(map[int]Artist)
)

// filters that come from the front and are used to filter results
type FilterT struct {
	BandSizeFilter           []int
	BandSizeFilterCheckboxes [10]int
	CreationYearStart        int
	CreationYearEnd          int
	FirstAlbumYearStart      int
	FirstAlbumYearEnd        int
	ConcertFilter            string
	SearchBar                string
}

// makes API call to download the artists and their relation data
func LoadArtistData() ([]Artist, []Relation, error) {
	var artists []Artist
	var relationData relationResponse

	artistChan := make(chan []Artist)
	artistErrChan := make(chan error)
	go func() {
		// Fetch artist data
		res, err := utils.LoadDataFromURL("https://groupietrackers.herokuapp.com/api/artists", &artists, 20)
		artistChan <- *res
		artistErrChan <- err
	}()

	relationChan := make(chan relationResponse)
	relationErrChan := make(chan error)
	go func() {
		// Fetch relations data
		res, err := utils.LoadDataFromURL("https://groupietrackers.herokuapp.com/api/relation", &relationData, 20)
		relationChan <- *res
		relationErrChan <- err
	}()

	artists, relations, err := consumeChannels(artistChan, relationChan, artistErrChan, relationErrChan)

	return artists, relations, err
}

func consumeChannels(artistChan chan []Artist, relationChan chan relationResponse, artistErrChan, relationErrChan chan error) ([]Artist, []Relation, error) {

	var artists []Artist
	var relations []Relation
	for range 4 {
		select {
		case artistsResp := <-artistChan:
			for _, artist := range artistsResp {
				ArtistMap[artist.ID] = artist
			}
			artists = artistsResp
		case relationResp := <-relationChan:
			for _, relation := range relationResp.Relations {
				ArtistRelationMap[relation.ID] = relation
			}
			relations = relationResp.Relations

		case errResp := <-artistErrChan:
			if errResp != nil {
				return nil, nil, errResp
			}

		case errResp := <-relationErrChan:
			if errResp != nil {
				return nil, nil, errResp
			}
		}
	}
	return artists, relations, nil
}
