package coordinate

import "github.com/golang/geo/s2"

// Coordinate represents a lat and lon in degrees.
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Coordinates []Coordinate

func (c Coordinate) ToS2LatLng() s2.LatLng {
	return s2.LatLngFromDegrees(c.Lat, c.Lon)
}
