package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

func DiabetesTest(pathFile string) {

	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	diabetesCSV := dataframe.ReadCSV(file)

	// calculating element numbers for row

	trainingNum := (4 * diabetesCSV.Nrow()) / 5 // => 4/5
	testNum := diabetesCSV.Nrow() / 5           // => 1/5

	if trainingNum+testNum < diabetesCSV.Nrow() {
		trainingNum++
	}

	// Creating the subset indices.
	trainingIdx := make([]int, trainingNum)
	testIdx := make([]int, testNum)

	// Enumerate the traning indices
	for i := 0; i < trainingNum; i++ {
		trainingIdx[i] = i
	}

	// Enumerate the test indices.
	for i := 0; i < testNum; i++ {
		testIdx[i] = i
	}

	// create the subset dataframes
	trainingDF := diabetesCSV.Subset(trainingIdx)
	testDF := diabetesCSV.Subset(testIdx)

	// cretate a map that will be used in writing the data
	// to files

	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	fmt.Println(setMap)

	// Now, create the respective files

	for idx, setName := range []string{"training.csv", "test.csv"} {

		// Save the fltered dartasete file.

		f, err := os.Create(setName)

		if err != nil {
			log.Fatal(err)
		}

		// Create a buffered writer.

		w := bufio.NewWriter(f)

		//Write the dataframe out as a CSV.

		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}

}

func main() {

	pathFile := "../diabetes.csv"
	DiabetesTest(pathFile)
}
