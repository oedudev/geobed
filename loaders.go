package geobed

import (
	"encoding/gob"
	"github.com/oedudev/geobed/assets"
)

func loadCityData() ([]GeobedCity, error) {
	fh, err := assets.LoadCityData()
	if err != nil {
		return nil, err
	}

	var cityData []GeobedCity
	dec := gob.NewDecoder(fh)
	err = dec.Decode(&cityData)
	if err != nil {
		return nil, err
	}
	return cityData, nil
}

func loadCountryData() ([]CountryInfo, error) {
	fh, err := assets.LoadCountryData()
	if err != nil {
		return nil, err
	}

	var countryInfo []CountryInfo
	dec := gob.NewDecoder(fh)
	err = dec.Decode(&countryInfo)
	if err != nil {
		return nil, err
	}
	return countryInfo, nil
}

func loadCityNameIdx() error {
	fh, err := assets.LoadCityNameIdx()
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(fh)
	cityNameIdx = make(map[string]int)
	err = dec.Decode(&cityNameIdx)
	if err != nil {
		return err
	}
	return nil
}
