package fetcher

import (
	"context"
	"sgougoupractice/fetch"
)

type RawLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Fetcher interface {
	Fetch(ctx context.Context, artistID int) (*RawLocation, error)
}

type FetcherFunc func(ctx context.Context, artistID int) (*RawLocation, error)

func (f FetcherFunc) Fetch(ctx context.Context, artistID int) (*RawLocation, error) {
	loc, err := fetch.FetchLocations(artistID)
	if err != nil {
		return nil, err
	}
	return &RawLocation{
		ID:        loc.ID,
		Locations: loc.Locations,
	}, nil
}
