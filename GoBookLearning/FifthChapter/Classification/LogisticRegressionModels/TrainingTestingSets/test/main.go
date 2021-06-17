package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

// predict makes a predictions based
// on our trained logistic regression model.
func predict(score float64) float64 {

	// calculate the predicted probability
	p := 1 / (1 + math.Exp(-13.65*score+4.89))

	// output the corresponding class.
	if p >= 0.5 {
		return 1.0
	}

	return 0.0
}

func main() {

	f, err := os.Open("../testing.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader  reading from the opened file
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	// observed and predicted will hold the parsed observed
	// and predicted values from the labeled data file.
	var observed []float64
	var predicted []float64

	// lne will track row numbers for logging
	line := 1

	// Read in the records looking for unexpected types in the
	// columns
	for {

		// Read in a row, Check uf we are at the end of the file
		record, err := reader.Read()
		if err == io.EOF {
			//log.Fatal(err)
			break
		}

		// skipe the header
		if line == 1 {
			line++
			continue
		}

		// Read in the observed value.
		observedVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Parsing line %d falled, unexpected type\n", line)
			continue
		}

		// Make the corresponding prediction
		score, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Parsing line %d falled, unexpected type\n", line)
			continue
		}

		predictedVal := predict(score)

		// Append the record to our slice
		// if it has the expected type.
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++

	}

	// this variable will hold our count of true positive and
	// true negative values.
	var truePosNeg int

	// Accumulate the true positive/negative count
	for idx, oVal := range observed {
		if oVal == predicted[idx] {
			truePosNeg++
		}
	}

	// Calculate 5the accuracy (subset accuracy)
	accuracy := float64(truePosNeg) / float64(len(observed))

	// output the accuracy value to standard out.

	fmt.Printf("\n Accuracy = %0.2f\n\n", accuracy)
}
