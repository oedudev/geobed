package geobed

import (
	"github.com/TomiHiltunen/geohash-golang"
)

// Contains all of the city and country data. Cities are split into buckets by country to increase lookup speed when the country is known.
type GeoBed struct {
	c  Cities
	co []CountryInfo
}

type Cities []GeobedCity

func (c Cities) Len() int {
	return len(c)
}
func (c Cities) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Cities) Less(i, j int) bool {
	return toLower(c[i].City) < toLower(c[j].City)
}

// A combined city struct (the various data sets have different fields, this combines what's available and keeps things smaller).
type GeobedCity struct {
	City    string
	CityAlt string
	// TODO: Think about converting this to a small int to save on memory allocation. Lookup requests can have the strings converted to the same int if there are any matches.
	// This could make lookup more accurate, easier, and faster even. IF the int uses less bytes than the two letter code string.
	Country    string
	Region     string
	Latitude   float64
	Longitude  float64
	Population int32
	Geohash    string
}

// Holds information about the index ranges for city names (1st and 2nd characters) to help narrow down sets of the GeobedCity slice to scan when looking for a match.
var cityNameIdx map[string]int

// Information about each country from Geonames including; ISO codes, FIPS, country capital, area (sq km), population, and more.
// Particularly useful for validating a location string contains a country name which can help the search process.
// Adding to this info, a slice of partial geohashes to help narrow down reverse geocoding lookups (maps to country buckets).
type CountryInfo struct {
	Country            string
	Capital            string
	Area               int32
	Population         int32
	GeonameId          int32
	ISONumeric         int16
	ISO                string
	ISO3               string
	Fips               string
	Continent          string
	Tld                string
	CurrencyCode       string
	CurrencyName       string
	Phone              string
	PostalCodeFormat   string
	PostalCodeRegex    string
	Languages          string
	Neighbours         string
	EquivalentFipsCode string
}

// Options when geocoding. For now just an exact match on city name, but there will be potentially other options that can be set to adjust how searching/matching works.
type GeocodeOptions struct {
	ExactCity bool
}

// Creates a new Geobed instance. You do not need more than one. You do not want more than one. There's a fair bit of data to load into memory.
func NewGeobed() (*GeoBed, error) {
	g := GeoBed{}

	var err error
	g.c, err = loadCityData()
	if err != nil {
		return nil, err
	}

	g.co, err = loadCountryData()
	if err != nil {
		return nil, err
	}

	err = loadCityNameIdx()
	if err != nil {
		return nil, err
	}

	return &g, nil
}

// ReverseGeocode - Provide the location given lat and long
func (g *GeoBed) ReverseGeocode(lat float64, lng float64) GeobedCity {
	c := GeobedCity{}

	gh := geohash.Encode(lat, lng)
	// This is produced with empty lat/lng values - don't look for anything.
	if gh == "7zzzzzzzzzzz" {
		return c
	}

	// Note: All geohashes are going to be 12 characters long. Even if the precision on the lat/lng isn't great. The geohash package will center things.
	// Obviously lat/lng like 37, -122 is a guess. That's no where near the resolution of a city. Though we're going to allow guesses.
	mostMatched := 0
	matched := 0
	for k, v := range g.c {
		// check first two characters to reduce the number of loops
		if v.Geohash[0] == gh[0] && v.Geohash[1] == gh[1] {
			matched = 2
			for i := 2; i <= len(gh); i++ {
				if v.Geohash[0:i] == gh[0:i] {
					matched++
				}
			}
			// tie breakers go to city with larger population (NOTE: There's still a chance that the next pass will uncover a better match)
			if matched == mostMatched && g.c[k].Population > c.Population {
				c = g.c[k]
			}
			if matched > mostMatched {
				c = g.c[k]
				mostMatched = matched
			}
		}
	}

	return c
}

// A slightly faster lowercase function.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
