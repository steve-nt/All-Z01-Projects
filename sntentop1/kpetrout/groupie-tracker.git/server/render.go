package server

import (
	"encoding/json"
	"fmt"
	"groupie/fetch"
	"groupie/filters"
	search "groupie/searchBar"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	dataMutex sync.RWMutex
	// TopTracks  []fetch.Track
)

type Developer struct {
	ID    int
	Name  string
	Image string
}

type PageData struct {
	Artists           []fetch.Artist
	ArtistsJSON       template.JS
	Query             string
	Dev               []Developer
	YearNow           int
	DatesAndLocations DatesAndLocations
	Src               string

	// Tracks  []fetch.Track
}

func HomeRender(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorPage(w, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("static/welcomepage.html", "static/footer.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	PageData := PageData{
		YearNow: time.Now().Year(),
	}

	tmpl.Execute(w, PageData)
}

func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	suggestions := search.GetSuggestions(query, artists)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func DataRender(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	if r.URL.Path != "/homepage" {
		ErrorPage(w, http.StatusNotFound)
		return
	}

	query := r.URL.Query().Get("query")
	searchedArtists := search.SearchArtists(query, artists)
	filteredArtists := filters.FilterArtists(searchedArtists, r.URL.Query())

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		tmpl, err := template.ParseFiles("static/overlay.html")
		if err != nil {
			ErrorPage(w, http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			Artists: filteredArtists,
			Query:   query,

			// Tracks:  TopTracks,
		}

		tmpl.ExecuteTemplate(w, "overlay", pageData)
		return
	}

	tmpl, err := template.ParseFiles("static/homepage.html", "static/overlay.html", "static/header.html", "static/footer.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Artists:           filteredArtists,
		Query:             query,
		YearNow:           time.Now().Year(),
		DatesAndLocations: getDatesAndLocations(artists),

		// Tracks:  TopTracks,
	}

	tmpl.Execute(w, pageData)
}

func DetailsRender(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/details" {
		ErrorPage(w, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("static/info.html", "static/header.html", "static/footer.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	id--
	if id < 0 || id >= len(artists) {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	var locPtrs []*fetch.Location
	for i := range artists[id].Locations {
		if artists[id].Locations[i].Lat == 0 && artists[id].Locations[i].Lng == 0 {
			locPtrs = append(locPtrs, artists[id].Locations[i])
		} else {
			// fmt.Println("Using cached data from file for", artists[id].Locations[i].Name)
		}
	}

	if len(locPtrs) > 0 {
		fmt.Println()
		fetch.FetchCoords(locPtrs)
		saveData("data.json", artists)
		fmt.Println("Cache coordinates for locations:")
		for i := range locPtrs {
			fmt.Println(locPtrs[i].Name, locPtrs[i].Lat, locPtrs[i].Lng)
		}
		fmt.Println()
	}

	pageData := PageData{
		Artists:           []fetch.Artist{artists[id]},
		YearNow:           time.Now().Year(),
		DatesAndLocations: getDatesAndLocations(artists),
		Src:               os.Getenv("GEO_API_SRC"),
		// Tracks:  TopTracks,
	}

	//JSON marshaling
	artistsJSON, err := json.Marshal(pageData.Artists)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	pageData.ArtistsJSON = template.JS(artistsJSON)

	tmpl.Execute(w, pageData)
}

func ContactRender(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		ErrorPage(w, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("static/contact.html", "static/header.html", "static/footer.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	PageData := PageData{
		YearNow:           time.Now().Year(),
		DatesAndLocations: getDatesAndLocations(artists),
		Dev: []Developer{
			{
				ID:    1,
				Name:  "CGkaldan",
				Image: "https://images.steamusercontent.com/ugc/1838044955942860859/48472B0E42426C20A2CFF6BE890A6C77072DB7AC/?imw=512&&ima=fit&impolicy=Letterbox&imcolor=%23000000&letterbox=false",
			},
			{
				ID:    2,
				Name:  "CMarkos",
				Image: "https://janegoodall.ca/wp-content/uploads/2021/10/Perrine20210511_003.jpg",
			},
			{
				ID:    3,
				Name:  "KPetrout",
				Image: "https://i.scdn.co/image/ab67616d0000b2735ce0d2cd473df39440c6350e",
			},
		},
	}
	tmpl.Execute(w, PageData)
}

func AboutRender(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		ErrorPage(w, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("static/about.html", "static/header.html", "static/footer.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	PageData := PageData{
		YearNow:           time.Now().Year(),
		DatesAndLocations: getDatesAndLocations(artists),
	}

	tmpl.Execute(w, PageData)
}
