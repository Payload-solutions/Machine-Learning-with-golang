package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

//const pathFile = "../AirPassengers.csv"
const newPathFile = "log_diff_series.csv"

func makingPlots() {
	file, err := os.Open(newPathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	passDF := dataframe.ReadCSV(file)

	fmt.Println(passDF.Names())
	passengersVals := passDF.Col("log_differenced_passengers").Float()
	timeVals := passDF.Col("time").Float()

	pts := make(plotter.XYs, passDF.Nrow()-1)

	var differenced [][]string

	differenced = append(differenced, []string{"time", "differenced_passengers"})
	//fmt.Println(timeVals)
	//fmt.Println(passengersVals)
	for i := 1; i < len(passengersVals); i++ {
		pts[i-1].X = timeVals[i]
		pts[i-1].Y = passengersVals[i] - passengersVals[i-1]
		differenced = append(differenced, []string{
			strconv.FormatFloat(timeVals[i], 'f', -1, 64),
			strconv.FormatFloat(passengersVals[i]-passengersVals[i-1], 'f', -1, 64),
		})
	}

	p := plot.New()
	p.X.Label.Text = "time"
	p.Y.Label.Text = "differenced passengers"
	p.Add(plotter.NewGrid())

	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(3)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)

	if err := p.Save(10*vg.Inch, 4*vg.Inch, "diff_passengers.png"); err != nil {
		log.Fatal(err)
	}

	nFile, err := os.Create("log_diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer nFile.Close()

	w := csv.NewWriter(nFile)
	w.WriteAll(differenced)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	makingPlots()
}
