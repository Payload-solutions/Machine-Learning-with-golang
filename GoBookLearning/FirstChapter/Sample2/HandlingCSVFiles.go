package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func ParsingDataFrames(pathFile string) {

	irisFile, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	fmt.Println(irisDF)

	filter := dataframe.F{
		Colname:    "species",
		Comparator: "==",
		Comparando: "Iris-versicolor",
	}

	// Now filter the dataFRame to see only the rows
	// where the iris species is "Iris-versicolor"

	versicolorDF := irisDF.Filter(filter)

	if versicolorDF.Err != nil {
		log.Fatal(versicolorDF.Err)
	}

	// Filter the dataframe agia, but only select out the
	// sepal_width and species columns.

	fmt.Println(versicolorDF)

	versicolorDF = irisDF.Filter(filter).Select([]string{"sepal_width", "species"})

	// Filter and select the dataframe again, but only display
	// the first three results

	versicolorDF = irisDF.Filter(filter).Select([]string{"sepal_width", "species"}).Subset([]int{0, 1, 2})
}

func main() {

	pathFile := "../../datasets/iris_labeled.csv"
	ParsingDataFrames(pathFile)
	fmt.Println("Hi ")
}
