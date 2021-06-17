package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/gonum/gonum/stat"
	"github.com/montanaflynn/stats"
)

// The goal of an Machine learning application
// is going down to the quality of your data;
// undertanding the data, and the evaluation
// and validation of the result.

// All three of these things require us to have
// an undersantung of statistics

func Distributions(pathFile string) {

	// Distributions
	/*
		a Distribution is a representation of how often values
		appear within a dataset, Let's say, for instance, that on thing
		you are tracjing as a data scientis is the daily sales of a
		certain product or service, and you have a long list(which you
		could represent as a vector or a part of matrix)
	*/

	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	irisDF := dataframe.ReadCSV(file)

	// Get the float values from the "sepal_length" column as
	// we will be looking at the measures for this variable.

	//for _, element := range irisDF.Records() {
	//	fmt.Println(element)
	//}
	//fmt.Println()

	petalLength := irisDF.Col("petal_length").Float()

	log.Println(petalLength)

	// Calculate theMena of the variable
	meanVal := stat.Mean(petalLength, nil)

	fmt.Printf("The meaning of value %0.1f\n\n", meanVal)

	// Calculate the Mode of the variable
	modeVal, modeCount := stat.Mode(petalLength, nil)

	fmt.Printf("the mode val: %0.1f and mode count are %0.1f\n\n", modeVal, modeCount)

	medianVal, err := stats.Median(petalLength)

	if err != nil {
		log.Fatal(err)
	}

	// For all our prupose, we let the decimal length in 2
	// By default

	fmt.Printf("The median val is %0.2f\n\n", medianVal)

	// Taking the mean and mediam from petal length, we can see
	// are not close it, that's mean, some values, drawing up or
	// down

}

func main() {

	pathFile := "../iris_labeled.csv"
	Distributions(pathFile)
}
