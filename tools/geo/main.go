package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"vlg/tools/model"

	"github.com/gocarina/gocsv"
	"github.com/timshannon/badgerhold/v4"
)

type USCensusResult struct {
	Result struct {
		Input struct {
			Address struct {
				Address string `json:"address"`
			} `json:"address"`
			Benchmark struct {
				IsDefault            bool   `json:"isDefault"`
				BenchmarkDescription string `json:"benchmarkDescription"`
				ID                   string `json:"id"`
				BenchmarkName        string `json:"benchmarkName"`
			} `json:"benchmark"`
		} `json:"input"`
		// will be empty if no match
		AddressMatches []struct {
			TigerLine struct {
				Side        string `json:"side"`
				TigerLineID string `json:"tigerLineId"`
			} `json:"tigerLine"`
			Coordinates struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"coordinates"`
			AddressComponents struct {
				Zip             string `json:"zip"`
				StreetName      string `json:"streetName"`
				PreType         string `json:"preType"`
				City            string `json:"city"`
				PreDirection    string `json:"preDirection"`
				SuffixDirection string `json:"suffixDirection"`
				FromAddress     string `json:"fromAddress"`
				State           string `json:"state"`
				SuffixType      string `json:"suffixType"`
				ToAddress       string `json:"toAddress"`
				SuffixQualifier string `json:"suffixQualifier"`
				PreQualifier    string `json:"preQualifier"`
			} `json:"addressComponents"`
			MatchedAddress string `json:"matchedAddress"`
		} `json:"addressMatches"`
	} `json:"result"`
}

type LibPostalResults []struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

const libPostalURL = `http://localhost:8899/parser`

func (l *LibPostalResults) getLastLabel(label string) string {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i].Label == label {
			return (*l)[i].Value
		}
	}
	return ""
}

func queryLibPostal(address string) (string, error) {
	jsonBody := []byte(fmt.Sprintf(`{"query": "%s"}`, address))
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(libPostalURL, "application/json", bodyReader)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var results LibPostalResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return "", err
	}
	line := ""
	state := results.getLastLabel("state")
	if state == "" {
		fmt.Println(results)
		return "", fmt.Errorf("no state found for %s", address)
	}
	line = state
	city := results.getLastLabel("city")
	if city == "" {
		return "", fmt.Errorf("no city found for %s", address)
	}
	line = city + ", " + line
	street := results.getLastLabel("road")
	if street == "" {
		return "", fmt.Errorf("no street found for %s", address)
	}
	line = street + ", " + line
	houseNumber := results.getLastLabel("house_number")
	if houseNumber == "" {
		return "", fmt.Errorf("no house number found for %s", address)
	}
	line = houseNumber + " " + line
	return line, nil
}

const censusURL = `https://geocoding.geo.census.gov/geocoder/locations/onelineaddress?address=%s&benchmark=2020&format=json`

// Geo finds geo-encoded coordinates from address lines in the database
func main() {

	store, err := model.GetStore(false)
	defer func() {
		_ = store.Close()
	}()
	if err != nil {
		panic(err)
	}

	// retrieve addresses
	var addresses []*model.Address
	err = store.Find(&addresses, badgerhold.Where("CountryCodes").Contains("USA"))
	if err != nil {
		panic(err)
	}

	var csvAddresses []*model.USCensusAddress
	misses := 0
	for i := range addresses {
		address := addresses[i]
		address.Normalize()
		if strings.HasPrefix(address.Name, "Unknown Address") {
			continue
		}

		// perform some cleanup
		value := address.Name
		value = strings.ReplaceAll(value, "U.S.A.", "")
		value = strings.ReplaceAll(value, "USA", "")
		value = strings.ReplaceAll(value, "c/o", "")
		value = strings.ReplaceAll(value, "C/o", "")
		value = strings.ReplaceAll(value, "C/O", "")

		value, err = queryLibPostal(value)
		if err != nil {
			log.Printf("No lib postal processing for %s\n", address.Name)
			misses++
			continue
		}

		// query census API
		var result USCensusResult
		resp, err := http.Get(fmt.Sprintf(censusURL, url.QueryEscape(value)))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			panic(err)
		}

		// skip if no match
		if len(result.Result.AddressMatches) == 0 {
			log.Printf("No census match for %s, libpostal term %s\n", address.Name, value)
			misses++
			continue
		}

		// save coordinates
		/*
			address.Location = &model.Location{
				Lat:  result.Result.AddressMatches[0].Coordinates.Y,
				Long: result.Result.AddressMatches[0].Coordinates.X,
			}
			address.GeoSource = "US Census"
			err = store.Update(address.NodeID, address)
			if err != nil {
				panic(err)
			}
		*/
		csvAddresses = append(csvAddresses, &model.USCensusAddress{
			NodeID: address.NodeID,
			Name:   address.Name,
			Lat:    result.Result.AddressMatches[0].Coordinates.Y,
			Long:   result.Result.AddressMatches[0].Coordinates.X,
		})

		time.Sleep(50 * time.Millisecond)

		if i%1000 == 0 {
			log.Printf("Processed %d addresses, misses so far %d\n", i, misses)
		}
	}

	addressesFile, err := os.OpenFile("addresses-geoencoded.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer addressesFile.Close()
	err = gocsv.MarshalFile(&csvAddresses, addressesFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Wrote %d addresses to file, %d misses\n", len(csvAddresses), len(addresses)-len(csvAddresses))
}
