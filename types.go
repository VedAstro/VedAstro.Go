package vedastro

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// GeoLocation contains the API's location name and coordinates, in degrees.
type GeoLocation struct {
	Name      string  `json:"Name"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
}

// NewGeoLocation builds a location. Coordinates are validated when it is serialized.
func NewGeoLocation(name string, longitude, latitude float64) GeoLocation {
	return GeoLocation{Name: name, Longitude: longitude, Latitude: latitude}
}
func (g GeoLocation) validate() error {
	if strings.TrimSpace(g.Name) == "" || math.IsNaN(g.Longitude) || math.IsNaN(g.Latitude) ||
		math.IsInf(g.Longitude, 0) || math.IsInf(g.Latitude, 0) ||
		g.Longitude < -180 || g.Longitude > 180 || g.Latitude < -90 || g.Latitude > 90 {
		return errors.New("vedastro: a location name and valid longitude/latitude are required")
	}
	return nil
}
func (g GeoLocation) String() string {
	return fmt.Sprintf("%s (%g, %g)", g.Name, g.Longitude, g.Latitude)
}

// MarshalJSON serializes the same Name, Longitude, and Latitude fields as Python.
func (g GeoLocation) MarshalJSON() ([]byte, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}
	type wire GeoLocation
	return json.Marshal(wire(g))
}

// TimeLayout is the API time format, including the explicit UTC offset.
const TimeLayout = "15:04 02/01/2006 -07:00"

// Time contains the API's standard time string and geographic location.
type Time struct {
	StdTime  string      `json:"StdTime"`
	Location GeoLocation `json:"Location"`
}

// NewTime validates an API time string such as "14:30 25/10/1992 +05:30".
func NewTime(stdTime string, location GeoLocation) (Time, error) {
	value := Time{StdTime: stdTime, Location: location}
	if err := value.validate(); err != nil {
		return Time{}, err
	}
	return value, nil
}

// NewTimeFromDate uses a Go time's calendar fields and UTC offset.
func NewTimeFromDate(value time.Time, location GeoLocation) (Time, error) {
	return NewTime(value.Format(TimeLayout), location)
}
func (t Time) validate() error {
	parsed, err := time.Parse(TimeLayout, t.StdTime)
	if err != nil || parsed.Format(TimeLayout) != t.StdTime {
		return errors.New("vedastro: time must use HH:MM DD/MM/YYYY +/-HH:MM")
	}
	return t.Location.validate()
}

// MarshalJSON serializes the same StdTime and Location fields as Python.
func (t Time) MarshalJSON() ([]byte, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}
	type wire Time
	return json.Marshal(wire(t))
}
func (t Time) String() string { return t.StdTime + " at " + t.Location.String() }
