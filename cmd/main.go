package main

import (
	"github.com/oedudev/geobed"
	"log"
)

func main() {

	g, err := geobed.NewGeobed()
	if err != nil {
		log.Fatal(err.Error())
	}

	g.ReverseGeocode(0, 0)

}
