package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const pathFIle = "../AirPassengers.csv"

// Time series and Anomaly detection

func TimeSeries() {

	file, err := os.Open(pathFIle)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	timeSeriesDF := dataframe.ReadCSV(file)

	fmt.Println(timeSeriesDF) // Printing to check the values

	pts := make(plotter.XYs, timeSeriesDF.Nrow())

	yVals := timeSeriesDF.Col("AirPassengers").Float()

	for i, floatVal := range timeSeriesDF.Col("time").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
	}

	p := plot.New()
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "AirPassengers"

	p.Add(plotter.NewGrid())

	// Add the line plot points for the time series.
	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "passengers_ts.png"); err != nil {
		log.Fatal(err)
	}

}

func main() {
	TimeSeries()
}
