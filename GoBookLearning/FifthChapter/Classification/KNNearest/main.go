package main

import (
	"fmt"
	"log"
	"math"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

func ReadingIris(pathFile string) {

	irisFile, err := base.ParseCSVToInstances(pathFile, true)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing our KNN model and performing
	// the cross-validation is quick and simple:
	knn := knn.NewKnnClassifier("euclidean", "linear", 2)

	// Use cross-fold validation to successively train and
	// evaluate the model on 5 folds of the data set
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(irisFile, knn, 5)
	if err != nil {
		log.Fatal(err)
	}

	// Get the mean, variance and standard deviation of the
	// accuracy for the cross validation.
	mean, variance := evaluation.GetCrossValidatedMetric(cv, evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	// Output the cross metics to standard out.
	fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)

}

func main() {

	pathFile := "../../datasets/iris_labeled.csv"
	ReadingIris(pathFile)
}
