package geocoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"groupie/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// This joints the path using platform correct symbols, / for linux \ for windows etc.
var filePath = filepath.Join("..", "geodata", "geodata.txt")

type Marker struct {
	Longitude   string `json:"lon"`
	Latitude    string `json:"lat"`
	Class       string `json:"class"`
	Addresstype string `json:"addresstype"`
	Location    string
}

// holds (or will hold) marker information for each location
type GeocodingCacheT struct {
	cache map[string]Marker
	mutex sync.Mutex
}

func (GC *GeocodingCacheT) get(location string) (Marker, bool) {
	GC.mutex.Lock()
	marker, ok := GC.cache[location]
	GC.mutex.Unlock()
	return marker, ok
}

func (GC *GeocodingCacheT) set(location string, marker Marker) {
	GC.mutex.Lock()
	GC.cache[location] = marker
	GC.mutex.Unlock()
}

// creates and initializes a geocoding cache instance
func makeGC() *GeocodingCacheT {
	GC := GeocodingCacheT{}
	GC.cache = make(map[string]Marker, 500)
	return &GC //returning it as a pointer so that there isn't a copied mutex
}

var GeocodingCache = makeGC() //map with geocoding data

// queue of locations that downloader will empty by downloading the coordinates, handlers add to this queue when their location wasn't found in the geocoding cache
type GeocodingQueueT struct {
	queue []string
	mutex sync.Mutex
}

// returns and removes the first element of the queue
func (GCQ *GeocodingQueueT) pop() string {
	GCQ.mutex.Lock()
	location := GCQ.queue[0]
	GCQ.queue = GCQ.queue[1:]
	GCQ.mutex.Unlock()
	return location
}

func (GCQ *GeocodingQueueT) add(location string) {
	GCQ.mutex.Lock()
	GCQ.queue = append(GCQ.queue, location)
	GCQ.mutex.Unlock()
}

var GeocodingQueue = GeocodingQueueT{} //queue with geocoding requests that need to be loaded

// Returns a query's coordinates, eventually. If it's not found in the cache, and order will be placed and it will keep checking the cache until it's found
func FetchCoordinates(location string) (Marker, error) {
	marker, ok := GeocodingCache.get(location)
	if ok {
		// marker found in cache
		marker.Location = utils.FixKey(location)
		return marker, nil
	}

	// marker wasn't found in cache, so we're adding it to the queue to be downloaded
	GeocodingQueue.add(location)

	retryLimit := 200 //retry limit so the goroutine isn't stuck here forever waiting for result
	for {
		time.Sleep(time.Millisecond * 250)

		marker, ok := GeocodingCache.get(location)
		if ok {
			//marker found (after downloader downloaded it)
			marker.Location = utils.FixKey(location)
			return marker, nil
		}

		retryLimit--
		if retryLimit == 0 {
			return Marker{}, errors.New("can't fetch marker, retry timeout reached")
		}
	}

}

// Loads geocode data from file
func LoadGeocodeData() error {

	//check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	//open file
	file, err := os.ReadFile(filePath)
	if nil != err {
		return err
	}

	file = []byte(strings.ReplaceAll(string(file), "\r", ""))

	//turn file into markers and add em to cache
	for _, line := range strings.Split(string(file), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ", ")
		GeocodingCache.cache[parts[0]] = Marker{Longitude: parts[1], Latitude: parts[2]}
	}

	return nil
}

// Saves geocode data into a file
func SaveGeocodeData() error {
	//make new file or truncate (annihilate) existing file
	file, err := os.Create(filePath)
	if nil != err {
		return err
	}
	defer file.Close()

	//turn cache into text
	var builder strings.Builder
	for location, marker := range GeocodingCache.cache {
		builder.WriteString(fmt.Sprintf("%s, %s, %s\n", location, marker.Longitude, marker.Latitude))
	}

	//save text to file
	_, err = file.WriteString(builder.String())
	if err != nil {
		fmt.Println("problem saving cache: ", err)
	}
	return nil
}

// Slowly downloads geocode data as found in the queue
func GeocodeDownloader() {
	for {
		time.Sleep(time.Millisecond * 100)

		//check if there's something in the geocoding queue
		if len(GeocodingQueue.queue) > 0 {
			locationOG := GeocodingQueue.pop()

			//verify that locatin is not actually in the cache already, we don't want to accidentally download something a second time
			_, ok := GeocodingCache.get(locationOG)
			if ok {
				continue
			}

			//format location for api query "osaka-japan" -> "japan+osaka"
			splitLocation := strings.Split(locationOG, "-")
			Aparts := utils.SplitByWords(splitLocation[0])
			Bparts := utils.SplitByWords(splitLocation[1])
			A := strings.Join(Aparts, "+")
			B := strings.Join(Bparts, "+")
			formattedLocation := B + ",+" + A
			formattedLocation = ManualLocationFixs(formattedLocation)

			fmt.Println(formattedLocation)

			link := "https://nominatim.openstreetmap.org/search?q="

			req, err := http.NewRequest("GET", link+formattedLocation+"&format=json", nil)
			if err != nil {
				fmt.Printf("Failed to create request: %v\n", err)
			}

			// Set the User-Agent header to a unique identifier for your app, required by the api to work
			req.Header.Set("User-Agent", "zone01Athens-groupie-tracker-v0.2 (aleksis.gioldaseas@outlook.com)")

			fmt.Println("Downloading marker for: ", locationOG)
			startTime := time.Now()

			// Send the request
			response, err := http.DefaultClient.Do(req)
			fmt.Println("Download complete.(?) It took: ", time.Since(startTime).Milliseconds(), "ms")

			if err != nil {
				fmt.Printf("Failed to fetch data: %v\n", err)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("Failed to read response: %v", err)
				response.Body.Close()
				time.Sleep(time.Second)
				continue
			}

			markers := []Marker{}

			json.Unmarshal(body, &markers)
			var realMarker Marker

			found := false
		loop1:
			for _, potentialMarker := range markers {
				switch potentialMarker.Class {
				case "town", "village", "county", "municipality", "district", "city":
					realMarker = potentialMarker
					found = true
					break loop1
				}
				switch potentialMarker.Addresstype {
				case "town", "village", "county", "municipality", "district", "city":
					realMarker = potentialMarker
					found = true
					break loop1
				}

			}

			if !found {
			loop2:
				for _, potentialMarker := range markers {
					fmt.Println("potential: ", potentialMarker)
					switch potentialMarker.Class {
					case "state", "province", "region", "boundary":
						realMarker = potentialMarker
						found = true
						break loop2
					}
					switch potentialMarker.Addresstype {
					case "state", "province", "region", "boundary":
						realMarker = potentialMarker
						found = true
						break loop2
					}

				}

			}

			if !found && len(markers) > 0 {
				realMarker = markers[0]
			}

			if realMarker.Latitude != "" && realMarker.Longitude != "" {
				GeocodingCache.set(locationOG, realMarker)
			}

			response.Body.Close()

			//sleeping to respect the API's guidelines, of only requesting once per second
			time.Sleep(time.Second)
		}
	}
}

// Saves geocode data permanently, at an interval
func GeocodeLogger() {
	entriesCount := len(GeocodingCache.cache)
	for {
		time.Sleep(time.Second) //every second check of something new was added to cache
		if len(GeocodingCache.cache) > entriesCount {
			entriesCount = len(GeocodingCache.cache)
			err := SaveGeocodeData()
			if err != nil {
				fmt.Println("ERROR saving file: ", err)
			}

			//after successful saving rest for 10 seconds as to not overuse the harddrive
			time.Sleep(time.Second * 10)
		}
	}
}

// manual fixes for random issues with geolocation api
func ManualLocationFixs(location string) string {
	replacer := strings.NewReplacer("los+angeles", "la")
	return replacer.Replace(location)
}
