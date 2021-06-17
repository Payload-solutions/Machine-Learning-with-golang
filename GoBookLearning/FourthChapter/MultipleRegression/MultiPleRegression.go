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

func MultipleLinear(pathFileTraining, pathFileTesting string) {

	file, err := os.Open(pathFileTraining)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// creating another csv files form the file opened

	reader := csv.NewReader(file)

	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	// In this case we are going to try and model our sales
	// by the Tv and Radio features plus an intercept

	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")
	r.SetVar(1, "Radio")

	// Loop over the csv records adding the training data.
	for i, record := range trainingData {

		if i == 0 {
			continue
		}

		// parsing the sales

		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// parsing tv and radio values

		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		raVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		// add these pounts to the regression value
		r.Train(regression.DataPoint(yVal, []float64{tvVal, raVal}))
	}

	// Train/fit the regression model
	r.Run()

	// output the trained model parameters.
	fmt.Printf("\nRegression Formula:\n%v\n\n", r.Formula)

	fileTest, err := os.Open(pathFileTesting)
	if err != nil {
		log.Fatal(err)
	}
	defer fileTest.Close()

	testReader := csv.NewReader(fileTest)

	testReader.FieldsPerRecord = 4

	testData, err := testReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the test data predicting y and evaluatin the
	// prediction with the mean absolute error.

	var mAE float64
	for i, record := range testData {

		// skip the header
		if i == 0 {
			continue
		}

		// Parsing the values
		yObserverd, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// parsing tv val and radio val

		tvVal, err := strconv.ParseFloat(record[0], 64)

		if err != nil {
			log.Fatal(err)
		}

		raVal, err := strconv.ParseFloat(record[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		// predict y with our trained model.
		yPredicted, err := r.Predict([]float64{tvVal, raVal})

		if err != nil {
			log.Fatal(err)
		}

		// Adding to the mean absolute error.
		mAE += math.Abs(yObserverd-yPredicted) / float64(len(testData))
	}

	fmt.Printf("\nMAE = %0.2f\n\n", mAE)

}

func main() {

	pathFileTraining := "training.csv"
	pathFileTesting := "test.csv"
	MultipleLinear(pathFileTraining, pathFileTesting)
}
