package geocode

import (
	"context"
	"fmt"
	"log"
	"sync"

	nominatim "sgougoupractice/geocoding/client"
	"sgougoupractice/helpers"
)

type LocationCoord struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type locationKey struct {
	City    string
	Country string
}

type GeocodingService struct {
	cl    nominatim.GeocodingClient
	cache map[locationKey]nominatim.Coords
	mu    sync.RWMutex

	flight map[locationKey]chan struct{}
	fmu    sync.RWMutex
}

func NewService(cl nominatim.GeocodingClient) *GeocodingService {
	return &GeocodingService{
		cl:     cl,
		cache:  make(map[locationKey]nominatim.Coords),
		flight: make(map[locationKey]chan struct{}),
	}
}

func (s *GeocodingService) Geocode(ctx context.Context, city, country string) (nominatim.Coords, error) {
	log.Print("start")
	key := locationKey{City: city, Country: country}

	s.mu.RLock()
	if c, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return c, nil
	}
	s.mu.RUnlock()

	wait := s.reserveFlight(key)
	if wait != nil {
		<-wait
		s.mu.RLock()
		c := s.cache[key]
		s.mu.RUnlock()
		if c == (nominatim.Coords{}) {
			return nominatim.Coords{}, fmt.Errorf("coords not found for %v", key)
		}
		return c, nil
	}
	defer s.releaseFlight(key)

	c, err := s.cl.Geocode(ctx, city, country)
	if err != nil {
		return nominatim.Coords{}, err
	}
	s.mu.Lock()
	s.cache[key] = c
	s.mu.Unlock()

	return c, nil
}

func (s *GeocodingService) BatchGeocode(ctx context.Context, slugs []string) ([]LocationCoord, error) {
	type result struct {
		coord LocationCoord
		err   error
	}

	// inCh  := make(chan string, len(slugs))
	out := make(chan result, len(slugs))
	var wg sync.WaitGroup

	for _, slug := range slugs {
		wg.Add(1)
		go func(sl string) {
			defer wg.Done()
			city, country := helpers.ParseKey(sl)
			c, err := s.Geocode(ctx, city, country)
			if err != nil {
				out <- result{err: err}
				return
			}
			out <- result{coord: LocationCoord{
				Name: sl,
				Lat:  c.Lat,
				Lon:  c.Lon,
			}}
		}(slug)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	var coords []LocationCoord
	for res := range out {
		if res.err != nil {
			continue

		}
		coords = append(coords, res.coord)
	}
	return coords, nil
}

func (s *GeocodingService) reserveFlight(k locationKey) (wait chan struct{}) {
	s.fmu.Lock()
	if ch, busy := s.flight[k]; busy {
		wait = ch
	} else {
		wait = nil
		s.flight[k] = make(chan struct{})
	}
	s.fmu.Unlock()
	return
}
func (s *GeocodingService) releaseFlight(k locationKey) {
	s.fmu.Lock()
	if ch, ok := s.flight[k]; ok {
		close(ch)
		delete(s.flight, k)
	}
	s.fmu.Unlock()
}
