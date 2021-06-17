package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Categorical metrics
// let's say that we have a model that is supposed to predict
// some discrete value, such as fraud/not fraud, standing/sitting/walking,
// approved/not approved, and so on. Our data might look something like
// the following

/*
	observed, predicted
		0,		0
		0,		1
		2,		2
		1,		1
		1,		1
		0,		0
		2,		0
		0,		0
		..		.


	Undertanding these metrics and determine which is appropriate for
	our use case, we need to realiza that there a number of different scenarios
	that could occur when we are making discrete predictions




	True Positive(TP): We predicted a certain category, and the
						observation was actually that category
						(for example, we predicted fraud and the
						observation was fraud)


	False Positive(FP): We predicted a certain category, but te
						observation was actually another category
						(for example, we predicted fraud but the
						observation was not fraud)

	True Negative(TN): We predicted that the observation wasn't a
						certain category and the observation was
						not that category (for example, we predicted
						not fraud and the observation was not fraud)


	False Negative(FN): We predicted that the observation wasn't
						a certain category, but the observation was
						actually that category(for example, we predicted
						not fraud but the observation was fraud)



	You can see that there are a number of ways we can combine,
	aggregate, and measure these sceneraios. In fact, we could
	even aggregate/measeure theme in som sor of unique way related
	to our specifi ploblem. HJowever, there are some pretty standard
	ways of aggregating and measuring, these scenerarios that result
	in the fllowinf common metrics:


	Accuracy: The percentage of predictions that were right, or
				(TP + TN)(TP + TN + FP + FN)

	Precision: The percentage of postive predcitions that were actually
				positive, or TP/(TP+FP)

	Recall: The percentage of postive predictions that were indfified as
			positive or TP/(TP+FN)


*/

func CategoricalMetrics(pathFile string) {

	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Let't to create a csv reader reading from the openend file

	reader := csv.NewReader(file)

	// Observed and predicted willl hold the parsed
	// observed and predicted values from the
	// labeled data file

	var observed []int
	var predicted []int

	// line will track wor numbers for logging

	line := 1

	for {

		// read in a row. Check if we are at the end of the file
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		// skip the header

		if line == 1 {
			line++
			continue
		}

		// Read in the observed and predicted values.
		observedVal, err := strconv.Atoi(record[0])

		if err != nil {
			fmt.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		predictedVal, err := strconv.Atoi(record[1])

		if err != nil {
			fmt.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		// Now append the record to our slice,
		// if it has the expected type
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}

	// this variable will hold our count of true positive
	// and true negative values

	var truePosNeg int

	// Accumulate the tru positive/negative count.

	for idx, oVal := range observed {
		if oVal == predicted[idx] {
			truePosNeg++
		}
	}

	// calculate the accuracy (subset accuracy).

	accuracy := float64(truePosNeg) / float64(len(observed))

	// output the accuracy to standard out.

	fmt.Printf("\nAccuracy =  %0.2f\n\n", accuracy)

	// classes contains the three possible classes in the labeled data.

	classes := []int{0, 1, 2}

	// loop over each class

	for _, class := range classes {

		// these variables will hold our count of true positives and
		// our count of false positives
		var truePos int
		var falsePos int
		var falseNeg int

		// Now Accummulate the true positive and false positive counts

		for idx, oVal := range observed {

			switch oVal {

			// if thje observed value is the relevant class,
			// we should check to see if the predicted that class

			case class:
				if predicted[idx] == class {
					truePos++
					continue
				}

				falseNeg++
			default:
				if predicted[idx] == class {
					falsePos++
				}
			}
		}
		// Calculate the precision
		precision := float64(truePos) / float64(truePos+falsePos)

		// Calculate the recall
		recall := float64(truePos) / float64(truePos+falseNeg)

		// Output the precision value to standard out

		fmt.Printf("\nPrecision (class %d) = %0.2f\n\n", class, precision)
		fmt.Printf("\nRecall (class %d = %0.2f\n\n", class, recall)
	}

}

func main() {

	pathFile := "../labeled.csv"
	CategoricalMetrics(pathFile)
}
