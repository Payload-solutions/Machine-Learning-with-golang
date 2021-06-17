package main

import (
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Histograms are by no means the only way to visyually gain
// an understanding out of data. Another commonly used type
// of plot is called box plot. This type of plot also gives
// us an idea about the grouping and spread of values in a
// distributions, but, as oopposed to the histogram, the box
// plot has several marked features that help guide out eyes:

func GraphicBoxPlots(pathFile string) {

	irisFile, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer irisFile.Close()

	IrisDF := dataframe.ReadCSV(irisFile)

	// Create the plot and set its title and axis label

	p := plot.New()

	p.Title.Text = "Box plots"
	p.Y.Label.Text = " Values"

	// Crate the box for our data.

	w := vg.Points(50)
	// Create a box plot for each of the feature
	// columns in the dataset.

	for idx, colName := range IrisDF.Names() {

		// If the column is one of the feature columns,
		// let's create a histogram of the values

		if colName != "species" {

			// Crate a plotter.Values value and fill it
			// with the values from the respective column of
			// the dataframe.

			v := make(plotter.Values, IrisDF.Nrow())

			for i, floatVal := range IrisDF.Col(colName).Float() {
				v[i] = floatVal
			}

			// Add the data to the plot

			b, err := plotter.NewBoxPlot(w, float64(idx), v)

			if err != nil {
				log.Fatal(err)
			}

			p.Add(b)
		}
	}

	p.NominalX("sepal_length", "sepal_width", "petal_length", "petall_width")
	if err := p.Save(6*vg.Inch, 8*vg.Inch, "boxplots.png"); err != nil {
		log.Fatal(err)
	}

}

func main() {

	pathFile := "../iris_labeled.csv"

	GraphicBoxPlots(pathFile)
}
