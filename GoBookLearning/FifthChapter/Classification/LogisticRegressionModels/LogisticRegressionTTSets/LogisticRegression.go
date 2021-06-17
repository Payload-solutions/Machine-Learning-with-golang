package main

import (
	"bufio"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func main() {

	f, err := os.Open("../../clean_loan_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a dataframe from the csv file
	// the types of the columns will be inferred
	loadDF := dataframe.ReadCSV(f)

	// Calculate the number of elements in each set
	trainNum := (4 * loadDF.Nrow()) / 5
	testNum := loadDF.Nrow() / 5
	if trainNum+testNum < loadDF.Nrow() {
		trainNum++
	}

	// Create the subset indices.
	trainIdx := make([]int, trainNum)
	testIdx := make([]int, testNum)

	// Enumerate the training indices
	for i := 0; i < trainNum; i++ {
		trainIdx[i] = i
	}

	// Enumerate the test indices
	for j := 0; j < testNum; j++ {
		testIdx[j] = trainNum + j
	}

	// Create the subset dataframes
	trainingDF := loadDF.Subset(trainIdx)
	testingDF := loadDF.Subset(testIdx)

	// Create a map that will be used in
	// writing the data to files
	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testingDF,
	}

	// Create the respective files

	for idx, setName := range []string{"training.csv", "testing.csv"} {

		// Save the filtered dataset file.
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}

		// Create a buffered writer.
		w := bufio.NewWriter(f)

		// Write the dataframe out as a csv.
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}
}
