package assets

import (
	"archive/zip"
	"bytes"
	"embed"
	"errors"
	"io"
	"io/fs"
)

const (
	ASSETS_COUNTRY_DATA  = "g.country_data.dmp"
	ASSETS_CITY_DATA     = "g.city_data.dmp.zip"
	ASSETS_CITY_NAME_IDX = "g.city_name_idx.dmp"
)

//go:embed g.country_data.dmp
var countryData embed.FS

//go:embed g.city_data.dmp.zip
var cityDataZip embed.FS

//go:embed g.city_name_idx.dmp
var cityNameIdx embed.FS

func LoadCountryData() (fs.File, error) {
	return countryData.Open(ASSETS_COUNTRY_DATA)
}

func LoadCityData() (io.ReadCloser, error) {
	zipData, err := cityDataZip.ReadFile(ASSETS_CITY_DATA)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(zipData)
	zipReader, err := zip.NewReader(readerAt, int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	if len(zipReader.File) != 1 {
		return nil, errors.New("more than one file on g.city_data.dmp.zip")
	}

	zipFile := zipReader.File[0]
	return zipFile.Open()
}

func LoadCityNameIdx() (fs.File, error) {
	return cityNameIdx.Open(ASSETS_CITY_NAME_IDX)
}
