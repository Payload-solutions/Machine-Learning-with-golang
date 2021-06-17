package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
)

/*
Differents types of results:

	Continuous: result such as total sales, stock price
				and temperature can take any continuis numerica
				value($ 12102.21, 92 degrees, and so on)

	Categorical: results such as fraud/not fraud, activity, and name
					that can take one of a finite number of categories
					(fraud, standing, Frank, and so on)

*/

/*

Continuos metrics:

	Let's say that we have a model that is supposed to predict
	some continuous value, like a stock price. Suppose that we
	have accumulated some predicted values that we can compare
	to actual observed values:

	observations, prediction

	22.1,17.9
	10.4,9.1
	9.3,7.8
	18.5,14.2
	12.9,15.6
	7.2,7.4
	11.8,9.7

	Now how do we measure the perfomance of this model? Well, the
	first step would be taking the difference between the observed and
	predictted values go get an <<error>>

	observation, prediction, error
	22.1,17.9,4.2
	10.4,9.1,1.3
	9.3,7.8,1.5
	18.5,14.2,4.3
	12.9,15.6,-2.7
	7.2,7.4,-0.2
	11.8,9.7,2.1


	The error gives us a general idea of how far off we were
	from the value thata we were supposed to predict. However,
	it's not really feasible por practical to look at all the error
	values individually, especially when there is a lot of data.
	There could be a million or more of these error values. Thus, we
	need a way to understand the errors in aggregate.



	the mean squared error (MSE) and mean absolute error (MAE)
	provide us with a view on errors in aggregate:

		MES or mean squared deviation (MSD) is the average of the
		squares of all the errors.

		MAE is the average fo the  absolute values of all errors.



*/

func ErrorEvaluations(pathFile string) {
	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Create a new CSV reader reading from the opened file

	reader := csv.NewReader(file)

	// Observed  and predicted will hod the parsed observed
	// and predicted values
	// from the continuos data file
	var observed []float64
	var predicted []float64

	// line will track row numbers for logging

	line := 1

	// Read in the records looking for unexpected
	// types in te columns

	for {

		// Read in a row. Check if we are at the end of the file

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		// skip the header.

		if line == 1 {
			line++
			continue
		}

		// Read in the observed and predicted values.
		observedVal, err := strconv.ParseFloat(record[0], 64)

		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		predictedVal, err := strconv.ParseFloat(record[1], 64)

		if err != nil {
			log.Printf("Parsing line %d filed unexpected type\n", line)
		}

		// Append the record to our slice, if it has the expected type.

		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}

	// Calculate the mean absolute errors and mean squared error.

	var mAE float64
	var mSE float64

	for idx, oVal := range observed {
		mAE += math.Abs(oVal-predicted[idx]) / float64(len(observed))
		mSE += math.Pow(oVal-predicted[idx], 2) / float64(len(observed))
	}

	mean := stat.Mean(observed, nil)

	fmt.Printf("\nMAE = %0.2f", mAE)
	fmt.Printf("\nMSE = %0.2f\n\n", mSE)
	fmt.Printf("The mean %0.2f\n", mean)

	// Calculate R² value.
	rSquared := stat.RSquaredFrom(observed, predicted, nil)

	fmt.Printf("\n the R² is %0.2f\n\n", rSquared)

	// 0.37
	// So, is this a good or a bad R-squared? Remember that
	// R-squared is a percentage, and higher percentages are
	// better. Here, we are capturing about 37% of the variance
	// in the variable, that we are trying to predict. Not very good

}

func main() {
	pathFile := "../continuous_data.csv"
	ErrorEvaluations(pathFile)
}
