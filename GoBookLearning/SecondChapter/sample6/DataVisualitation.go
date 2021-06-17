package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Data visualitations to quantify how distributions work!

func GraphHistogram(pathFile string) {

	irisCsv, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer irisCsv.Close()

	IrisDF := dataframe.ReadCSV(irisCsv)

	// Now create a histogrma for each of the feature columns in the
	// dataset.

	for _, colName := range IrisDF.Names() {

		// If the column is one of the feature columns, let's
		// create  a histogram of the values

		if colName != "species" {
			// Create a plotter.Values value and fill it with the
			// values from the respective column of the dataframe

			v := make(plotter.Values, IrisDF.Nrow())

			for i, floatVal := range IrisDF.Col(colName).Float() {
				v[i] = floatVal
			}

			// Make a plot and set its title.
			// plot.New not return error variable
			p := plot.New()

			p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

			// Create a histogram of our values drawn
			// from the standard normal.

			h, err := plotter.NewHist(v, 16)

			if err != nil {
				log.Fatal(err)
			}

			// Normalize the histogram
			h.Normalize(1)

			// Add the histogram to the plot
			p.Add(h)

			// Save plot to a PNG file.

			if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
				log.Fatal(err)
			}

			// Note: note that we have normalied out histograms(with h.Normalize())
			// this is typical because often you wil wanto to compare
			// the differen distributions side by side.
		}
	}
}

func main() {

	pathFile := "../iris_labeled.csv"

	GraphHistogram(pathFile)
}
