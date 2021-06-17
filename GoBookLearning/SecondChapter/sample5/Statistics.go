package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// Meausres of spread or dispersion:
// 		- Maximum: The highest value of the distribution
// 		- Minimum: the lowest value of the distribution
// 		- Range: The difference between the maximum and minimum
// 		- Variance: This measure is calculated by taking each of
// 					the values in the distribution, calculating
// 					each one's difference from  the distribution's
// 					 menas, squaring this difference, adding it
// 					to the other squared differences, and dividing by
// 					the number of values in the distribution
// 		- Standard deviation: The square root of the variance.
// 		- Quantiles/Quartiles: Similar to he median, these measures define
// 								cut-off points in the distribution
// 								where a certain number of lower values are
// 								below the measure and certaing number of higher
// 								values are above the measure.

func StatisticsFunctions(pathFile string) {

	irisFile, err := os.Open(pathFile)
	if err != nil {

		log.Fatal(err)
	}

	defer irisFile.Close()

	// Create a dataframe from the csv file.
	irisDF := dataframe.ReadCSV(irisFile)

	// Get the float values from "sepal_length"

	sepalLength := irisDF.Col("petal_length").Float()

	//  Min & Max variable

	minVal := floats.Min(sepalLength)
	maxVal := floats.Max(sepalLength)

	fmt.Printf("The max %0.2f and min %0.2f variable\n\n", maxVal, minVal)

	// Calculate the median of the varibale

	mediam := maxVal - minVal
	fmt.Printf("the mediam of value is %v\n\n", mediam)

	fmt.Println("\nCalulate the variance")

	variance := stat.Variance(sepalLength, nil)
	fmt.Printf("The variance %0.2f\n\n", variance)

	fmt.Println("\n Calculate the statndar deviation of the variable")

	stdDevVal := stat.StdDev(sepalLength, nil)
	fmt.Printf("The standard deviation of the value is %v\n\n", stdDevVal)

	fmt.Println("\n sorting the values")

	inds := make([]int, len(sepalLength))
	floats.Argsort(sepalLength, inds)

	fmt.Printf("The array sorted %v\n\n", sepalLength)

	fmt.Println("\nCalculate quartiles and quantiles")

	quant25 := stat.Quantile(0.25, stat.Empirical, sepalLength, nil)
	quant50 := stat.Quantile(0.50, stat.Empirical, sepalLength, nil)
	quant75 := stat.Quantile(0.75, stat.Empirical, sepalLength, nil)

	fmt.Printf("The 25 Quantile %v\n\n", quant25)
	fmt.Printf("The 50 Quantile %v\n\n", quant50)
	fmt.Printf("The 75 Quantile %v\n\n", quant75)

}

func main() {

	pathFile := "../iris_labeled.csv"

	StatisticsFunctions(pathFile)
}
