package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const citiBikeURL = "https://gbfs.citibikenyc.com/gbfs/en/station_status.json"

type stationData struct {
	LastUpdated int `json:"last_updated"`
	TTL         int `json:"ttl"`
	Data        struct {
		Stations []station `json:"stations"`
	} `json:"data"`
}

type station struct {
	ID                string `json:"station_id"`
	NumBikesAvailable int    `json:"num_bikes_available"`
	NumBikesDisabled  int    `json:"num_bike_disabled"`
	NumDocksAvailable int    `json:"num_docks_available"`
	NumDocksDisabled  int    `json:"num_docks_disabled"`
	IsInstalled       int    `json:"is_installed"`
	IsRenting         int    `json:"is_renting"`
	IsReturning       int    `json:"is_returning"`
	LastReported      int    `json:"last_reported"`
	HasAvailableKeys  bool   `json:"eightd_has_available_keys"`
}

func main() {

	response, err := http.Get(citiBikeURL)

	if err != nil {
		log.Fatal(err)
	}
	// Now Read the body of the repsonse into []byte

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var sd stationData

	if err := json.Unmarshal(body, &sd); err != nil {
		log.Fatal(err)

	}

	//fmt.Println(sd)

	fmt.Printf("%+v\n\n", sd.Data.Stations[0])

	// Now let's sat that we have the Citi Bike station
	// in our stationData struct value and we want to
	// save that data out to a file. We can do this
	// with json.marshal

	/**
	 * Marshal the data
	 */

	outputData, err := json.Marshal(sd)

	if err != nil {
		log.Fatal(err)
	}

	// save the marshalled data to a file

	if err := ioutil.WriteFile("citibike.json", outputData, 0644); err != nil {
		log.Fatal(err)
	}
}
