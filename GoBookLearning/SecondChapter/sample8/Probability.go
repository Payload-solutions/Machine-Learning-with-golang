package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
)

func Probability() {

	observed := []float64{48, 52}
	expected := []float64{50, 50}

	chiSquare := stat.ChiSquare(observed, expected)
	fmt.Printf("Chi square value is %0.2f\n\n", chiSquare)
}

func CalculatingPValues() {

	// A survey of local residents revelas that 60% of all
	// residents get no regular exercise, 25% exercise
	// sporadically and 15% exercise regularly. After doing
	// some fancy modeling and putting some community services
	// in place, the survey was repeated with the same
	// questions. The follow up survey was completed by 500
	// residents with the following results:

	/*
		No regular exercise: 260
		Sporadic exercise: 135
		Regular Exercise: 105
		Total : 500
	*/

	// Now we, want to determine if there is evidence for a
	// statistically significant shift in the responses of he
	// residents. Our null and alternate hypotheses are as follows:

	/*
		H0: The deviations from the previuosly obserced percentages
			are due to pure chance

		Ha: The deviations are due to some underlying effect ouside of
			pure chance(possibly our new community services)
	*/

	// first let's calculate our test statustuc using the
	// chi-square test statitic

	observed := []float64{
		260.0, // the number of observed with no regular exercise
		135.0, // the number of observed with sporatic exercise
		105.0, // the number of observed with regular exercise
	}
	// Define the total observed.

	totalObserved := 500.0

	// Calculate the expected frequencies
	// (again assuming the null hypothesis)

	expected := []float64{
		totalObserved * 0.60,
		totalObserved * 0.25,
		totalObserved * 0.15,
	}

	chiSquare := stat.ChiSquare(observed, expected)

	fmt.Printf("\nChi-square: %0.2f\n", chiSquare)

}

func main() {
	//Probability()
	CalculatingPValues()

}
