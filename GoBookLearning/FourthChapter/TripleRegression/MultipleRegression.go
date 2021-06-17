package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

func MakingTripleModel(trainingFile, testingFile string) {

	fileTraining, err := os.Open(trainingFile)
	if err != nil {
		log.Fatal(err)
	}

	defer fileTraining.Close()

	reader := csv.NewReader(fileTraining)
	reader.FieldsPerRecord = 4

	dataReader, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Setting three independents variables

	var r regression.Regression

	r.SetObserved("Sales")
	r.SetVar(0, "TV")
	r.SetVar(1, "Radio")
	r.SetVar(2, "Newspaper")

	for i, record := range dataReader {

		if i == 0 {
			continue
		}

		// Parsing yvalue

		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		raVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		newsVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// added each value into the statement
		r.Train(regression.DataPoint(yVal, []float64{tvVal, raVal, newsVal}))
	}
	r.Run()

	fmt.Printf("\nRegression Formula\n%v\n\n", r.Formula)

	testFile, err := os.Open(testingFile)
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()

	testReader := csv.NewReader(testFile)
	testReader.FieldsPerRecord = 4

	testData, err := testReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var mAE float64

	for i, record := range testData {

		if i == 0 {
			continue
		}

		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		raVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		newsVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		yPredicted, err := r.Predict([]float64{tvVal, raVal, newsVal})

		if err != nil {
			log.Fatal(err)
		}
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))

	}

	fmt.Printf("\nMAE = \n%v\n\n", mAE)

}

func main() {
	pathFileTraining := "training.csv"
	pathFileTesting := "test.csv"

	MakingTripleModel(pathFileTraining, pathFileTesting)
}
