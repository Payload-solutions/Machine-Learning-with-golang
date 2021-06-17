package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/berkmancenter/ridge"
	"github.com/gonum/matrix/mat64"
)

func main() {

	trainFile, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer trainFile.Close()

	reader := csv.NewReader(trainFile)
	reader.FieldsPerRecord = 4

	rawData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// FeatureData wull hold all the float values that will
	// eventually be used to form our matrix of features

	featureData := make([]float64, 4*len(rawData))
	yData := make([]float64, len(rawData))

	// FeatureIndex and yIndex will track the current
	// index of the matrix values
	var featureIndex int
	var yIndex int

	// Sequentially move the rows into a slice of floats.

	for idx, record := range rawData {
		if idx == 0 {
			continue
		}
		// loop over the float columns.
		for i, val := range record {
			valParsed, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			if i < 3 {
				// Add an intercept to the model.

				if i == 0 {
					featureData[featureIndex] = 1
					featureIndex++
				}

				// Add the float value to the slice of feature
				// floats.

				featureData[featureIndex] = valParsed
				featureIndex++
			}

			if i == 3 {
				// Add the float value to the slice of y floats.
				yData[yIndex] = valParsed
				yIndex++
			}
		}
	}

	// Form the matrices that will be input to our regression.
	features := mat64.NewDense(len(rawData), 4, featureData)
	y := mat64.NewVector(len(rawData), yData)

	// Create a new RifgeRegression value, where 1.0
	// is the penalty value.

	r := ridge.New(features, y, 1.0)

	// Train the regression model
	r.Regress()

	// Print our regression formula.

	c1 := r.Coefficients.At(0, 0)
	c2 := r.Coefficients.At(1, 0)
	c3 := r.Coefficients.At(2, 0)
	c4 := r.Coefficients.At(3, 0)

	fmt.Printf("\nREgression Formula:\n")
	fmt.Printf("y = %0.3f + %0.3fTV + %0.3f Radio + %0.3f NewsPaper\n\n", c1, c2, c3, c4)

	// Testing the another file

	testFile, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer testFile.Close()

	testReader := csv.NewReader(testFile)
	testReader.FieldsPerRecord = 4

	// reading all csv lines

	testData, err := testReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// looping over the houldout data predicting y and
	// evaluationg the prediction with the mean absoljute error

	var mAE float64

	for i, record := range testData {

		if i == 0 {
			continue
		}

		// parsing the sales
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

		newsVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Predict y with our trained model.
		yPredicted := Predict(tvVal, raVal, newsVal)
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))

	}

	fmt.Printf("\nMAE = %0.2f\n\n", mAE)

}

func Predict(tv, radio, newspaper float64) float64 {
	return (3.038 + tv*0.047 + 0.177*radio + 0.001*newspaper)
}
