package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gonum/matrix/mat64"
)

// logistic implements the logistic function, which
// is used in logistic regression.
func logistic(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// LogisticRegression fits a logistic regression model
// for the given data

/*func logisticRegression(features *mat64.Dense,
	labels []float64, numSteps int,
	learningRate float64) []float64 {

	// Initialize random weights
	_, numWeights := features.Dims()
	weights := make([]float64, numWeights)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for idx := range weights {
		weights[idx] = r.Float64()
	}

	// Iteratively optimize the weights
	for i := 0; i < numSteps; i++ {

		// Initialize a variable to accumulate error for this
		// iteration
		var sumError float64

		// Make predictions for each label and accumulate error.

		for idx, label := range labels {

			// Get the features coresponding to this label
			featureRow := mat64.Row(nil, idx, features)

			// Calculate the error for this iteration's weights
			pred := logistic(featureRow[0]*weights[0] + featureRow[1]*weights[1])
			predError := label - pred
			sumError += math.Pow(predError, 2)

			// update the features wights
			for j := 0; j < len(featureRow); j++ {
				weights[j] += learningRate * predError * pred * (1 - pred) * featureRow[j]
			}
		}
	}

	return weights

}*/
// logisticRegression fits a logistic regression model
// for the given data.
func logisticRegression(features *mat64.Dense, labels []float64, numSteps int, learningRate float64) []float64 {

	// Initialize random weights.
	_, numWeights := features.Dims()
	weights := make([]float64, numWeights)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for idx := range weights {
		weights[idx] = r.Float64()
	}

	// Iteratively optimize the weights.
	for i := 0; i < numSteps; i++ {

		// Initialize a variable to accumulate error for this iteration.
		var sumError float64

		// Make predictions for each label and accumlate error.
		for idx, label := range labels {

			// Get the features corresponding to this label.
			featureRow := mat64.Row(nil, idx, features)

			// Calculate the error for this iteration's weights.
			pred := logistic(featureRow[0]*weights[0] + featureRow[1]*weights[1])
			predError := label - pred
			sumError += math.Pow(predError, 2)

			// Update the feature weights.
			for j := 0; j < len(featureRow); j++ {
				weights[j] += learningRate * predError * pred * (1 - pred) * featureRow[j]
			}
		}
	}

	return weights
}

func main() {

	file, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSVreader reading from the opened file
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2

	rawData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Feature data and labels will hold all the float values
	// That will eventually be used in our training
	featureData := make([]float64, 2*(len(rawData)-1))
	labels := make([]float64, len(rawData)-1)

	// FeatureIndex will track the current index of the features
	// matrix values
	var FeatureIndex int

	// Sequentually move the rows into the slices of floats.
	for idx, record := range rawData {
		if idx == 0 {
			continue
		}

		// Add the FICO score feature.
		featureVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		featureData[FeatureIndex] = featureVal

		// Add an intercept
		featureData[FeatureIndex+1] = 1.0

		// Increment ouw feature row.
		FeatureIndex += 2

		// Add the class label
		labelVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		labels[idx-1] = labelVal
	}

	// Form a matrix from the features
	//fmt.Println("\nfeature data", featureData)
	features := mat64.NewDense(len(rawData)-1, 2, featureData)
	//fmt.Println("\nFeatures", features)

	// Train the logistic regression model.
	weights := logisticRegression(features, labels, 100, 0.3)

	// output the logistic Regression model formula to stdout.
	formula := "p = 1 / (1 + exp(- m1*FICO.score - m2))"
	fmt.Printf("\n%s\n\nm1 = %0.2f\nm2 = %0.2f\n\n", formula, weights[0],
		weights[1])
}
