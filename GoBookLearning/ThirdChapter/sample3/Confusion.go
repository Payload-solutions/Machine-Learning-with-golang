package main

import (
	"fmt"

	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/stat"
)

/*
Confusion matrices, AUC, and ROC

In addition to calculatin individual numerical metrics for our
models, there are a varietu of techiques to combine various metrics
into a form that gives you a more complete representation of model
perfomance. These include, but are certainly not limite to,
consusion matrices and
area under the curve(AUC)/Receiver Operating Characteristic (ROC)curves

Confision matrices allow us to visualize the varous TP, TN, FP and FN
values that we predict in a two-dimensional format. A confusion matrix
has rows corresponding to the categories that you ere supposed to predict,
and columns corresponding to categories that were predicted.
Then, the value of each element is the corresponding count:

								Predicted

						Fraud		Not fraud
Observed	fraud		  TP			FN
			Not fraund	  FP			TN

*/

func CalculateAUCROC() {

	// Define our scores and classes.
	scores := []float64{0.1, 0.35, 0.4, 0.8}
	classes := []bool{true, false, true, false}

	// Calculate the true positive rates (recalls) and
	// false positive rates.
	tpr, fpr, _ := stat.ROC(nil, scores, classes, nil)

	// Compute the Area Under Curve.
	auc := integrate.Trapezoidal(fpr, tpr)

	// Output the results to standard out.
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)
	fmt.Printf("auc: %v\n", auc)
}

func main() {
	CalculateAUCROC()
}
