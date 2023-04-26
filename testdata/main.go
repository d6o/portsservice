package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const (
	codeLen   = 9
	codeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Location struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

const (
	maxRandNum       = 1000
	maxCoordinateNum = 90
	filePermissions  = 0o644
)

func main() {
	jsonData := make(map[string]Location)

	for i := 0; i < 1000000; i++ {
		location := Location{
			Name:    "Location " + strconv.Itoa(i+1),
			City:    "City " + strconv.Itoa(rand.Intn(maxRandNum)),
			Country: "Country " + strconv.Itoa(rand.Intn(maxRandNum)),
			Alias:   []string{},
			Regions: []string{},
			Coordinates: []float64{
				float64(rand.Intn(maxCoordinateNum)),
				float64(rand.Intn(maxCoordinateNum)),
			},
			Province: "Province " + strconv.Itoa(rand.Intn(maxRandNum)),
			Timezone: "Timezone " + strconv.Itoa(rand.Intn(maxRandNum)),
			Unlocs:   []string{},
			Code:     generateCode(i),
		}

		jsonData[location.Code] = location
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON data: %v", err)
	}

	err = os.WriteFile("data.json", jsonBytes, filePermissions)
	if err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	fmt.Println("JSON data successfully generated and saved to data.json")
}

func generateCode(num int) string {
	code := ""
	for i := 0; i < codeLen; i++ {
		code += string(codeChars[num%len(codeChars)])
		num /= len(codeChars)
	}
	return code
}
