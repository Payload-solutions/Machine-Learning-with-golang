package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
)

// Declare the input and output directory flags

var (
	inModelDirPtr = flag.String("inModelDir", "", "The directory containing the model.")
	inDirPtr      = flag.String("inDir", "", "The directory containing the trainig data.")
	outDirPts     = flag.String("outDir", "", "The output directory.")
)

const pathFile = "../Advertising.csv"
const trainFile = "train.csv"
const testFile = "test.csv"
const diabetesFile = "../diabetes.csv"

// ModelInfo includes the information about the
// model that is putput from the training
type ModelInfo struct {
	Intercept    float64           `json:"intercept"`
	Coefficients []CoefficientInfo `json:"coefficients"`
}

// CoefficientsInfo include information about a
// particular model coefficient.
type CoefficientInfo struct {
	Name        string  `json:"name"`
	Coefficient float64 `json:"coefficient"`
}

// a prediction data include the data necessart to make
// a predictions and encodes the outpud
type PredictionData struct {
	Prediction      float64          `json:"predicted_diabetes_progression"`
	IndependentVars []IndependentVar `json:"independent_variables"`
}

// Independent var include information about and a
// value for independent variable
type IndependentVar struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func CreateDataTrainTest() {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dataDF := dataframe.ReadCSV(file)

	// Crating the num of the train and test data
	trainNum := (4 * dataDF.Nrow()) / 5
	testNum := dataDF.Nrow() / 5

	if trainNum+testNum < dataDF.Nrow() {
		trainNum++
	}

	trainSet := make([]int, trainNum)
	testSet := make([]int, testNum)

	for i := 0; i < trainNum; i++ {
		trainSet[i] = i
	}

	for j := 0; j < testNum; j++ {
		testSet[j] = j + trainNum
	}

	trainSubset := dataDF.Subset(trainSet)
	testSubset := dataDF.Subset(testSet)

	modelMap := map[int]dataframe.DataFrame{
		0: trainSubset,
		1: testSubset,
	}

	for idx, colName := range []string{"train.csv", "test.csv"} {

		f, err := os.Create(colName)
		if err != nil {
			log.Fatal(err)
		}

		w := bufio.NewWriter(f)

		if err := modelMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}
}

func TrainTestModel() {
	newTrainFile, err := os.Open(trainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer newTrainFile.Close()

	reader := csv.NewReader(newTrainFile)
	reader.FieldsPerRecord = 4

	dataReader, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")

	for i, record := range dataReader {

		// Skiping the index
		if i == 0 {
			continue
		}

		// Parsing the data
		yVals, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		tvVals, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		r.Train(regression.DataPoint(yVals, []float64{tvVals}))
	}

	// Training and fit the regression model
	r.Run()

	fmt.Printf("\nFormula from the predicted value %v\n\n", r.Formula)

	newTestFile, err := os.Open(testFile)
	if err != nil {
		log.Fatal(err)
	}
	defer newTestFile.Close()

	testReader := csv.NewReader(newTestFile)
	testReader.FieldsPerRecord = 4

	testData, err := testReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var mAE float64

	for i, colVal := range testData {

		// skiping the header
		if i == 0 {
			continue
		}

		yObserved, err := strconv.ParseFloat(colVal[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		tvValObserved, err := strconv.ParseFloat(colVal[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		yPredicted, err := r.Predict([]float64{tvValObserved})
		if err != nil {
			log.Fatal(err)
		}

		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	fmt.Printf("\nMAE = %0.3f\n\n", mAE)

	// Fill in the model information

}

func diabetesPrediction() {

	file, err := os.Open(diabetesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 11

	dfReader, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var r regression.Regression
	r.SetObserved("Diabetes progression")
	r.SetVar(0, "bmi")
	r.SetVar(1, "ltg")

	for i, record := range dfReader {
		if i == 0 {
			continue
		}

		yVals, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			log.Fatal(err)
		}

		bmiVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		ltgVal, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			log.Fatal(err)
		}

		r.Train(regression.DataPoint(yVals, []float64{bmiVal, ltgVal}))

	}

	r.Run()

	fmt.Printf("\nFormula from the predicted value %v\n\n", r.Formula)

	modelInfo := ModelInfo{
		Intercept: r.Coeff(0),
		Coefficients: []CoefficientInfo{
			{
				Name:        "bmi",
				Coefficient: r.Coeff(1),
			},
			{
				Name:        "ltg",
				Coefficient: r.Coeff(2),
			},
		},
	}

	// Marshal the model information
	outputData, err := json.MarshalIndent(modelInfo, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// Save the marshalled output to a file,
	if err := ioutil.WriteFile(filepath.Join(*outDirPts, "model.json"), outputData, 0644); err != nil {
		log.Fatal(err)
	}
}

func Predict(modelInfo *ModelInfo, predictionData *PredictionData) error {

	// Initialize the prediction value
	// To intercept

	prediction := modelInfo.Intercept

	// Create a map for independent variable coefficients
	coeffs := make(map[string]float64)
	varNames := make([]string, len(modelInfo.Coefficients))

	for idx, coeff := range modelInfo.Coefficients {
		coeffs[coeff.Name] = coeff.Coefficient
		varNames[idx] = coeff.Name
	}

	// Create a map of the independen variable values.
	varVals := make(map[string]float64)
	for _, indVar := range predictionData.IndependentVars {
		varVals[indVar.Name] = indVar.Value
	}

	// loop over the independent variables.
	for _, varName := range varNames {

		// Get the coefficient
		coeff, ok := coeffs[varName]
		if !ok {
			return fmt.Errorf("couldn't find model coefficients %s", varName)
		}

		// Get the variable value.
		val, ok := varVals[varName]
		if !ok {
			return fmt.Errorf("expected  a value for variable %s", varName)
		}

		// Add to prediction
		prediction += coeff * val
	}

	// Add the prediction to prediction data
	predictionData.Prediction = prediction

	return nil

}

func main() {
	flag.Parse()

	//_, err1 := os.Open(trainFile)
	//_, err2 := os.Open(testFile)
	//if err1 != nil && err2 != nil {
	//	log.Printf("[*] Files Doesn't exists; creating files...")
	//	createDataTrainTest()
	//}
	//trainTestModel()

	// Train/ Fit the singled regression model
	// with Sajari regression

	// Fill in the model information.

	diabetesPrediction()
}
