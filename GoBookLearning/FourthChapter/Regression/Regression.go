package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Seting the Y variable as dependent variable
func CheckingCorrelation(pathFile string) {

	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	dataCSV := dataframe.ReadCSV(file)

	// Extract the target of the column
	// chosse the column, and the name
	yVals := dataCSV.Col("Sales").Float()
	fmt.Println(yVals)
	// create a scatter a polot for each of the features in the
	// dataset

	for _, colName := range dataCSV.Names() {
		// pts will hold the value for plotting

		pts := make(plotter.XYs, dataCSV.Nrow())

		// fill pts with data.

		for i, floatVal := range dataCSV.Col(colName).Float() {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		// create the plot

		p := plot.New()

		p.X.Label.Text = colName
		p.Y.Label.Text = "y"

		p.Add(plotter.NewGrid())

		s, err := plotter.NewScatter(pts)

		if err != nil {
			log.Fatal(err)
		}

		s.GlyphStyle.Radius = vg.Points(3)
		p.Add(s)
		// save the plot to a PNG file

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_scatter.png"); err != nil {
			log.Fatal(err)
		}
	}

}

// Making plot charts
func ReadingAdvertising(pathFile string) {

	file, err := os.Open(pathFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	advertDF := dataframe.ReadCSV(file)

	/*// Use the describe metjod to calculate
	// summary statistics for all the columns in one shot

	adverSummary := advertDF.Describe()

	fmt.Println(adverSummary)*/

	// Extract the target column.
	//yVals := advertDF.Col("Sales").Float()

	// create a histogram for each of the columns in the dataset.
	for _, colName := range advertDF.Names() { // Printing the headers
		// create a plotter. values value and fill it with the
		// values from the respective column of the dataframe.

		plotVals := make(plotter.Values, advertDF.Nrow())

		for i, floatVal := range advertDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		// Making a plot and set its title

		p := plot.New()
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

		// Create a histogram of our values drawn
		// from the standard normal.

		h, err := plotter.NewHist(plotVals, 16)

		if err != nil {
			log.Fatal(err)
		}

		// Normalizing the data.
		h.Normalize(1)

		// Add the histogram to the plot.

		p.Add(h)

		// saving the plot to a PNG file.

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}

}

// Setting the data to be training and testing
func CreatingDataTraining(pathFile string) {

	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	advertDF := dataframe.ReadCSV(file)

	// Calculating, number of elements in each set

	traningNum := (4 * advertDF.Nrow()) / 5
	testNum := advertDF.Nrow() / 5

	if traningNum+testNum < advertDF.Nrow() {
		traningNum++
	}

	// creating a subset indices.
	trainingIdx := make([]int, traningNum)
	testIdx := make([]int, testNum)

	// Enumerate the traning indices

	for i := 0; i < traningNum; i++ {
		trainingIdx[i] = i
	}

	// Enumerate test indices
	for i := 0; i < testNum; i++ {
		testIdx[i] = traningNum + i
	}

	// create a the subset dataframes.
	// Every subset will cut the data frame in portions
	// depend of the number of rows that you will set
	trainingDF := advertDF.Subset(trainingIdx)
	testDF := advertDF.Subset(testIdx)

	// Create a mapa that will be used in writing the data
	// to files

	setMap := map[int]dataframe.DataFrame{
		0: trainingDF,
		1: testDF,
	}

	// create the respective files

	for idx, setName := range []string{"training.csv", "test.csv"} {

		// save the filtered dataset file.
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}

		// create a buffered writer
		w := bufio.NewWriter(f)

		// Write the daraframe out as a CSV

		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}

}

func StartTrainingModel(pathTraining, pathTesting string) {

	file, err := os.Open(pathTraining)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a new CSV reader reading from the opened file

	reader := csv.NewReader(file)

	// Read in all of the CSV records

	reader.FieldsPerRecord = 4

	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// in this case we are going to try and model our sales(y)
	// by the TV feature plus an intercep. As such, let's create
	// the struct needed to train a model using Sajari repository

	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")

	// Now loop of records in the csv, adding the training
	// data to the refression value

	for i, record := range trainingData {

		// skip the header
		if i == 0 {
			continue
		}

		// parse the Sales regression measure, or "y"

		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// parse the Tv value.

		tvVal, err := strconv.ParseFloat(record[0], 64)

		if err != nil {
			log.Fatal(err)
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(yVal, []float64{tvVal}))
	}

	// Train/fit the regression model.
	r.Run()

	// Output the trained model parameters.
	fmt.Printf("\nRegression formula: \n%v\n\n", r.Formula)

	// Mental checking, if the scatter plots between the TV and
	// Sales was up and to the right, that's mean, that the
	// correlation is positive, therefore slope should be positive
	//		Regression formula:
	//		Predicted = 7.0688 + TV*0.0489

	testFile, err := os.Open(pathTesting)

	if err != nil {
		log.Fatal(err)
	}

	defer testFile.Close()

	// creating a reder

	testReader := csv.NewReader(testFile)
	testReader.FieldsPerRecord = 4
	testData, err := testReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	// loop over the test data predicting y and evaluating
	// the prediction with the mean absolute error.

	var mAE float64

	for i, record := range testData {

		// skiping the header
		if i == 0 {
			continue
		}

		// parse the observed Sales or "y"

		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// parse the TV value.
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// predicted y with our trained model
		yPredicted, err := r.Predict([]float64{tvVal})

		if err != nil {
			log.Fatal(err)
		}

		// ADd the to the mean absolute error.
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))

	}

	// Output the MAE to standar out.
	fmt.Printf("MAE = %0.2f\n\n", mAE)

	/*
		Regression formula:
			Predicted = 7.0688 + TV*0.0489

			MAE = 3.01

			How we know if 3.01 is good or bad ?
			When we show all the values in the dataset
			The mean sales was 14.02 and the desviation is
			5.21. Thus, out MAE is less, than standard deviations
			of our salues values, and is about 20% of the mean
			value, and our model has some predictive power.

	*/

}

// Without importing the librarie
func Prediction(tv float64) float64 {
	return 7.07 + tv*0.05
}

func MakingPrediction(pathFile string) {

	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Create a dataframe from the CSV file.
	advertDF := dataframe.ReadCSV(file)

	// Extracting the target column
	yVals := advertDF.Col("Sales").Float()

	// pts will hold the values for plotting.
	pts := make(plotter.XYs, advertDF.Nrow())

	//ptsPred will hold the predicted values for plotting
	ptsPred := make(plotter.XYs, advertDF.Nrow())

	// Fill pts with data.
	for i, floatVal := range advertDF.Col("TV").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]

		ptsPred[i].X = floatVal
		ptsPred[i].Y = Prediction(floatVal)
	}

	// creating a plot
	p := plot.New()
	p.X.Label.Text = "TV"
	p.Y.Label.Text = "Sales"
	p.Add(plotter.NewGrid())

	// Add the scatter plot points for the observations
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}

	s.GlyphStyle.Radius = vg.Points(3)

	//Adding the line plot points for the predictions
	l, err := plotter.NewLine(ptsPred)
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	// Save the plot to a PNG file.
	p.Add(s, l)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "regression_line.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//pathTesting := "test.csv"
	//pathTraining := "training.csv"

	pathGeneral := "Advertising.csv"
	//readingAdvertising(pathFile)
	//CheckingCorrelation(pathFile)
	//CreatingDataTraining(pathFile)
	//StartTrainingModel(pathTraining, pathTesting)
	MakingPrediction(pathGeneral)
}
