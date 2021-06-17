package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type CSVRecords struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
	ParseError  error
}

func readCSvFIle(pathFile string) {

	file, err := os.Open(pathFile)

	if err != nil {

		log.Fatal(err)
	}
	defer file.Close()
	readers := csv.NewReader(file)

	var csvData []CSVRecords

	for {
		records, err := readers.Read()

		if err == io.EOF {
			break
		}

		var csvRecord CSVRecords

		for idx, value := range records {

			// validate that value in the
			// 4th place is an empty string
			if idx == 4 {
				if value == "" {
					log.Printf("Unexpected type in column %d\n", idx)
					csvRecord.ParseError = fmt.Errorf("empty string value")
					break
				}
				csvRecord.Species = value
				continue
			}

			var floatValue float64
			// if the value cannot be parsed from string into float
			// log an error

			if floatValue, err = strconv.ParseFloat(value, 64); err != nil {
				log.Printf("Unexpected type in column %d\n", idx)
				csvRecord.ParseError = fmt.Errorf("couldn't not parsed float")
				break
			}

			// Otherwise, add the the value in respective place

			switch idx {
			case 0:
				csvRecord.SepalLength = floatValue
			case 1:
				csvRecord.SepalWidth = floatValue
			case 2:
				csvRecord.PetalLength = floatValue
			case 3:
				csvRecord.PetalWidth = floatValue
			}
		}
		if csvRecord.ParseError == nil {
			csvData = append(csvData, csvRecord)
		}
	}

	fmt.Printf("Successfully parsed %d lines\n", len(csvData))

	for _, element := range csvData {
		fmt.Println(element)
	}

}

func main() {

	pathFile := "datasets/iris.csv"

	readCSvFIle(pathFile)

}
