package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/mash/gokmeans"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const pathFile = "../datasets/fleet_data.csv"

func Histogram() {
	clusFile, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer clusFile.Close()

	driverDF := dataframe.ReadCSV(clusFile)

	fmt.Println(driverDF.Describe())

	for _, colName := range driverDF.Names() {

		if colName == "Driver_ID" {
			continue
		}

		// First step
		// making the plot vals taking as reference the number of rows
		plotVals := make(plotter.Values, driverDF.Nrow())
		for i, plotVal := range driverDF.Col(colName).Float() {
			plotVals[i] = plotVal
		}

		// Second step
		// Creating a new plot to be shiped
		p := plot.New() // This one only take one response, another old versions, this take an error parameter
		p.Title.Text = fmt.Sprintf("Histogram of %s", colName)

		// third step
		// makes the kinna plot that you wanna create
		pHist, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}
		pHist.Normalize(1) // Normalize makes reference to the data in normal order

		// Four step
		// add the histogram to the plot
		p.Add(pHist)

		// Five step
		// save the values in png files

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hitogram.png"); err != nil {
			log.Fatal(err)
		}
	}
}

func NewPlotting() {

	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal("error file: ", err)
	}
	defer file.Close()

	driverDF := dataframe.ReadCSV(file)

	// Distance
	yVals := driverDF.Col("Distance_Feature").Float()

	pts := make(plotter.XYs, driverDF.Nrow())

	// Fill pts with data
	for i, floatVal := range driverDF.Col("Speeding_Feature").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
	}

	p := plot.New()

	p.X.Label.Text = "Speeding"
	p.Y.Label.Text = "Distance"
	p.Add(plotter.NewGrid())

	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	// Added color and style
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	// Saving the point in a PNG file
	p.Add(s)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "fleet_data_scatter.png"); err != nil {
		log.Fatal(err)
	}

}

func GeneratingClusters() {

	f, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 3

	// Initialize a slice of gokmeans.Node to
	// hold our input data
	var data []gokmeans.Node

	// loop over the records creating our slice of
	for {

		// Read in our record and check for erros
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// skip the header
		if record[0] == "Driver_ID" {
			continue
		}

		// Initialize a point
		var point []float64

		// Fill in our point
		for i := 1; i < 3; i++ {

			val, err := strconv.ParseFloat(record[i], 64)
			if err != nil {
				log.Fatal(err)
			}

			// Append this value to our point
			point = append(point, val)
		}

		// Append the point to the data
		data = append(data, gokmeans.Node{point[0], point[1]})

	}
	// Then generating our clusters is as easy as calling
	// the gomeans.Train(...)
	success, centroids := gokmeans.Train(data, 2, 50)
	fmt.Println(success)
	if !success {
		log.Fatal("Couldn't generate clusters")
	}

	fmt.Printf("\nThe centroids for our clusters are: ")

	for _, cendroid := range centroids {
		fmt.Println(cendroid)
	}
}

func EvaluatingClusters() {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Reading dataFrame
	driverDF := dataframe.ReadCSV(file)

	// Extracting the distance
	yVals := driverDF.Col("Distance_Feature").Float()

	// As we already make two clusters in the above function
	// We need to hold tww kind of clusters in something like
	// tensors or [][]float64
	var clusterOne [][]float64
	var clusterTwo [][]float64

	// fill the clusters with data
	for i, xVal := range driverDF.Col("Speeding_Feature").Float() {

		distanceOne := floats.Distance([]float64{yVals[i], xVal},
			[]float64{50.05, 8.83}, 2)
		distanceTwo := floats.Distance([]float64{yVals[i], xVal},
			[]float64{180.02, 18.29}, 2)

		if distanceOne < distanceTwo {
			clusterOne = append(clusterOne, []float64{xVal, yVals[i]})
			continue
		}
		clusterTwo = append(clusterTwo, []float64{xVal, yVals[i]})
	}

	// pts* will hodl the values for plotting
	ptsOne := make(plotter.XYs, len(clusterOne))
	ptsTwo := make(plotter.XYs, len(clusterTwo))

	for i, point := range clusterOne {
		ptsOne[i].X = point[0]
		ptsOne[i].Y = point[1]
	}

	for i, point := range clusterTwo {
		ptsTwo[i].X = point[0]
		ptsTwo[i].Y = point[1]
	}

	p := plot.New()
	p.X.Label.Text = "Spedding"
	p.Y.Label.Text = "Distance"
	p.Add(plotter.NewGrid())

	sOne, err := plotter.NewScatter(ptsOne)
	if err != nil {
		log.Fatal(err)
	}
	sOne.GlyphStyle.Radius = vg.Points(3)
	sOne.GlyphStyle.Shape = draw.PyramidGlyph{}

	sTwo, err := plotter.NewScatter(ptsTwo)
	if err != nil {
		log.Fatal(err)
	}
	sTwo.GlyphStyle.Radius = vg.Points(3)
	sTwo.GlyphStyle.Shape = draw.PyramidGlyph{}

	p.Add(sOne, sTwo)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "fleet_data_clusters.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCluster 1 metric: %0.2f\n", WithinClusterMean(clusterOne,
		[]float64{50.05, 8.83}))
	fmt.Printf("\nCluster 2 metric: %0.2f\n", WithinClusterMean(clusterTwo,
		[]float64{180.02, 18.29}))

}

func WithinClusterMean(cluster [][]float64, centroids []float64) float64 {
	// MeanDistance will hold our result
	var meanDistance float64

	// loop over the points in the cluster

	for _, point := range cluster {
		meanDistance += floats.Distance(point, centroids, 2)
	}

	return meanDistance
}

func main() {

	//Histogram()
	//NewPlotting()
	//GeneratingClusters()
	EvaluatingClusters()
}
