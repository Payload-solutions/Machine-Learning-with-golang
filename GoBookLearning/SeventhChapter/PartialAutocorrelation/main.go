package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kniren/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const pathFile = "../AirPassengers.csv"

// pacf calculates the partial autocorrelation for a series
func PACF(x []float64, lag int) float64 {

	// Create a regression.Regression value needed to train
	// a model using github.com/sajari/regression.
	var r regression.Regression
	r.SetObserved("x")

	// Defint eh current lag and all of the intermediate lags.
	for i := 0; i < lag; i++ {
		r.SetVar(i, "x"+strconv.Itoa(i))
	}

	// shift the series
	xAdj := x[lag:]

	// Loop pver the series creating the data set
	// for the Regression
	for i, xVal := range xAdj {

		// Loop over the intermediate lags to build up
		// out independent variables
		laggedVariables := make([]float64, lag)
		for idx := 1; idx <= lag; idx++ {

			// Get the lagged series variables.
			laggedVariables[idx-1] = x[lag+i-idx]
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(xVal, laggedVariables))
	}

	// Fit the regression
	r.Run()

	return r.Coeff(lag)
}

func Reading() {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	passendDF := dataframe.ReadCSV(file)

	// Get the time and passeners as a slice of floats
	passengers := passendDF.Col("AirPassengers").Float()

	// loop over various calues of lag in the series
	fmt.Println("Partial Autocorrelation")
	for i := 1; i < 11; i++ {

		// Calculate the partial autocorrelation
		pac := PACF(passengers, i)
		fmt.Printf("\nLag %d period: %0.2f\n", i, pac)

	}
}

func MakingPlotts() {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	passDF := dataframe.ReadCSV(file)

	passengers := passDF.Col("AirPassengers").Float()

	p := plot.New()
	p.X.Label.Text = "LAG"
	p.Y.Label.Text = "PACF"
	p.Y.Min = 0
	p.Y.Max = 0

	w := vg.Points(3)

	numLags := 20

	pts := make(plotter.Values, numLags)

	for i := 1; i < numLags; i++ {
		pts[i-1] = PACF(passengers, i)
	}
	bars, err := plotter.NewBarChart(pts, w)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(bars)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "PACF_bars.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//Reading()
	MakingPlotts()
}
