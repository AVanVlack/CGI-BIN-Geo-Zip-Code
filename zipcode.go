package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	country_code := "US"

	// If not gob, make gob
	if _, err := os.Stat(country_code + ".gob"); os.IsNotExist(err) {
		zipcodeMap, err := LoadDataset(country_code)
		if err != nil {
			fmt.Printf("Status: 500 Internal Server Error\n\n")
			fmt.Print("Could not find zipcode dataset")
			return
		}
		err = MakeGob(country_code, zipcodeMap)
		if err != nil {
			fmt.Printf("Status: 500 Internal Server Error\n\n")
			fmt.Print("Could not make inital dataset")
			return
		}
	}

	zipcodeMap, err := LoadGob(country_code)
	if err != nil {
		fmt.Printf("Status: 500 Internal Server Error\n\n")
		fmt.Print("Could not find zipcode dataset")
		return
	}

	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	var zip Zipcode
	if err := json.Unmarshal([]byte(line), &zip); err != nil {
		fmt.Printf("Status: 400 Bad Request\n\n")
		fmt.Printf("Could not parse json string %s \n", line)
		return
	}

	foundedZipcode := zipcodeMap.DatasetList[zip.Zipcode]
	if foundedZipcode.ZipCode != zip.Zipcode {
		fmt.Printf("Status: 404 Not Found\n\n")
		fmt.Printf("Could not find the zipcode %s.\n", zip.Zipcode)
		return
	}

	json, err := json.Marshal(foundedZipcode)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Content-type: application/json\n\n")
	fmt.Printf(string(json))
}

func LoadGob(country_code string) (Zipcodes, error) {
	var data Zipcodes

	// open data file
	dataFile, err := os.Open(country_code + ".gob")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataFile.Close()

	return data, nil
}

func MakeGob(country_code string, dataset Zipcodes) error {
	file, err := os.Create(country_code + ".gob")
	if err != nil {
		return fmt.Errorf("zipcodes: error while creating file %v", err)
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(dataset)
	file.Close()
	return nil
}

type Zipcode struct {
	CountryCode string `json:"countryCode"`
	Zipcode     string `json:"zipcode"`
}

type ZipCodeLocation struct {
	ZipCode   string  `json:"zipCode"`
	PlaceName string  `json:"placeName"`
	AdminName string  `json:"adminName"`
	Lat       float64 `json:"latitude"`
	Lon       float64 `json:"longitude"`
}

type Zipcodes struct {
	DatasetList map[string]ZipCodeLocation
}

func LoadDataset(country_code string) (Zipcodes, error) {
	file, err := os.Open(country_code + ".txt")
	if err != nil {
		return Zipcodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	zipcodeMap := Zipcodes{DatasetList: make(map[string]ZipCodeLocation)}
	for scanner.Scan() {
		splittedLine := strings.Split(scanner.Text(), "\t")
		if len(splittedLine) != 12 {
			return Zipcodes{}, fmt.Errorf("zipcodes: file line does not have 12 fields")
		}
		lat, errLat := strconv.ParseFloat(splittedLine[9], 64)
		if errLat != nil {
			return Zipcodes{}, fmt.Errorf("zipcodes: error while converting %s to Latitude", splittedLine[9])
		}
		lon, errLon := strconv.ParseFloat(splittedLine[10], 64)
		if errLon != nil {
			return Zipcodes{}, fmt.Errorf("zipcodes: error while converting %s to Longitude", splittedLine[10])
		}

		zipcodeMap.DatasetList[splittedLine[1]] = ZipCodeLocation{
			ZipCode:   splittedLine[1],
			PlaceName: splittedLine[2],
			AdminName: splittedLine[3],
			Lat:       lat,
			Lon:       lon,
		}
	}

	if err := scanner.Err(); err != nil {
		return Zipcodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	return zipcodeMap, nil
}
