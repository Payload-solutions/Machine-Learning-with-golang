package main

import (
	"fmt"
	"log"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/naive"
)

func convertToBinary(src base.FixedDataGrid) base.FixedDataGrid {
	b := filters.NewBinaryConvertFilter()
	attrs := base.NonClassAttributes(src)
	for _, a := range attrs {
		b.AddAttribute(a)
	}
	b.Train()
	ret := base.NewLazilyFilteredInstances(src, b)
	return ret
}

func main() {

	trainingData, err := base.ParseCSVToInstances("../../training.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a new Naive Bayes classifier
	nb := naive.NewBernoulliNBClassifier()

	// Fit the native bayes classifier
	nb.Fit(convertToBinary(trainingData))

	// Read in the loan test data set into golearn instasnces
	// this time we wil utilize a template of the previous set
	// of instances to validate the format of the test set.
	testData, err := base.ParseCSVToTemplatedInstances("../../testing.csv", true, trainingData)
	if err != nil {
		log.Fatal(err)
	}

	// Make our predictions
	predictions, err := nb.Predict(convertToBinary(testData))
	if err != nil {
		log.Fatal(err)
	}
	// Generate a confusion matrix
	cm, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		log.Fatal(err)
	}
	// Retrieve the accuracy
	accuracy := evaluation.GetAccuracy(cm)
	fmt.Printf("\nAccuracy: %0.2f\n\n", accuracy)

}
