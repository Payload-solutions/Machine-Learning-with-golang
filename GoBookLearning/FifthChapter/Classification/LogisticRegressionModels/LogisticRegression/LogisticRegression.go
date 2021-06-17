package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	scoreMax = 830.0
	scoreMin = 640.0
)

func MakingLogistict(trainingFile string) {

	trainFile, err := os.Open(trainingFile)
	if err != nil {
		log.Fatal(err)
	}
	defer trainFile.Close()

	// The same of the all time. Create a new CSV reader
	reader := csv.NewReader(trainFile)
	reader.FieldsPerRecord = 2

	// Read all in the csv file
	trainData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Something new, create a outputfile
	f, err := os.Create("new_data_load.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// create a csv writer
	writer := csv.NewWriter(f)

	// Sequeantally move the rows writing out the parsed values.
	for idx, record := range trainData {

		if idx == 0 {
			// writing the header to the output file.
			if err := writer.Write(record); err != nil {
				log.Fatal(err)
			}
			continue
		}

		// initialize a slice to hold our parsed values.
		outRecord := make([]string, 2)

		// Parse and standardize the FICO score.
		score, err := strconv.ParseFloat(strings.Split(record[0], "-")[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		outRecord[0] = strconv.FormatFloat((score-scoreMin)/(scoreMax-scoreMin), 'f', 4, 64)

		// Parse the interest rate class.
		rate, err := strconv.ParseFloat(strings.TrimSuffix(record[1], "%"), 64)
		if err != nil {
			log.Fatal(err)
		}

		if rate <= 12.0 {
			outRecord[1] = "1.0"

			// Write the record to the output file.
			if err := writer.Write(outRecord); err != nil {
				log.Fatal(err)
			}
			continue
		}

		outRecord[1] = "0.0"

		// write the record to the output file.
		if err := writer.Write(outRecord); err != nil {
			log.Fatal(err)
		}
	}
	// write any buffered data to the underlying writer (standard output).
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}

func UndertandingFile(readerFile string) {

	file, err := os.Open(readerFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// creating a dataframe from csv file
	loanDF := dataframe.ReadCSV(file)

	// Using describe metjod to calculate summary statistics
	// for all of the columns in one shot
	loanSummary := loanDF.Describe()
	fmt.Println(loanSummary)

	// Now creating a histogram for each of the columns
	// in the dataset
	for _, colName := range loanDF.Names() {

		// Create a plotter.Values value and fill it tieh the
		// the values from the respective column of the dataframe
		plotVals := make(plotter.Values, loanDF.Nrow())
		for i, floatVal := range loanDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}
		p := plot.New()
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

		// Create a histogram of our values.
		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}
		h.Normalize(1)

		// Add the histogram to the plot.
		p.Add(h)

		// Save the plot to a png file
		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}

}

func main() {

	//trainingFile := "../../loan_data.csv"
	//MakingLogistict(trainingFile)
	/*
		Technique to optimize the coefficient/weights
		is called stochastic gradient.

		For the implementation of this, we need the following
		as input:

			features -> A pointer to a mat64.Dense matrix
			labels -> A slice of floats including all of the class
			numSteps -> A maximun number of iterations for the
						implementation.
			learningRate -> An adjustable parameter that helps
							with the convergence of the optimization.
	*/

	readerFile := "new_data_load.csv"
	UndertandingFile(readerFile)
}
