package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Coords struct {
	Lat float64
	Lon float64
}

/*
	type GeocodeRequest struct {
		City    string
		Country string
	}
*/
const (
	mqEndpoint = "https://www.mapquestapi.com/geocoding/v1/address"
)

type GeocodingClient interface {
	//Geocode(ctx context.Context, address string) (Coords, error)
	Geocode(ctx context.Context, city, country string) (Coords, error)
}

type MQClient struct {
	http *http.Client
	//baseURL   *url.URL
	userAgent string
	key       string
}

func NewMQ(apiKey, userAgent string, hc *http.Client) (*MQClient, error) {
	//u, err := url.Parse("https://nominatim.openstreetmap.org/search")
	if hc == nil {
		hc = &http.Client{Timeout: 5 * time.Second}
	}

	return &MQClient{
		key:       apiKey,
		http:      hc,
		userAgent: userAgent,
	}, nil
}

type mqResp struct {
	Results []struct {
		Locations []struct {
			LatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"locations"`
	} `json:"results"`
}

func (mq *MQClient) Geocode(ctx context.Context, city, country string) (Coords, error) {
	loc := strings.Trim(strings.Join([]string{city, country}, ","), " ,")
	if loc == "" {
		return Coords{}, fmt.Errorf("mapquest: empty location")
	}
	u, _ := url.Parse(mqEndpoint)
	q := u.Query()
	q.Set("key", mq.key)
	q.Set("location", loc)
	q.Set("maxResults", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Coords{}, fmt.Errorf("creating geocode request: %w", err)
	}
	req.Header.Set("User-Agent", mq.userAgent)
	resp, err := mq.http.Do(req)
	if err != nil {
		log.Printf("%v", err)
		return Coords{}, fmt.Errorf("performing geocode request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return Coords{}, fmt.Errorf("geocode API returned %d: %s", resp.StatusCode, string(body))
	}
	var out mqResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Printf("%v", err)
		return Coords{}, fmt.Errorf("decoding geocode response: %w", err)
	}
	if len(out.Results) == 0 || len(out.Results[0].Locations) == 0 {
		return Coords{}, fmt.Errorf("no geocode result for %q", city)
	}

	ll := out.Results[0].Locations[0].LatLng

	if ll.Lat == 0 && ll.Lng == 0 {
		return Coords{}, fmt.Errorf("mapquest: zero coords for %q", loc)
	}
	return Coords{Lat: ll.Lat, Lon: ll.Lng}, nil

}
