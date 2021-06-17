package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

const trainFile = "../data/train.csv"
const testFile = "../data/test.csv"

// NeuralNet contains all of the information
// that defines a trained neural network
type neuralNet struct {
	config  neuralNetConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOut    *mat.Dense
	bOut    *mat.Dense
}

// neuralNetConfig defines our neural network
// architecture and learning parameters.
type neuralNetConfig struct {
	inputNeurons  int
	outputNeurons int
	hiddenNeurons int
	numEpochs     int
	learningRate  float64
}

// Sigmoid implements the sigmoid function
// for use in activation functions.
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// SigmoidPrime implements the derivative
// of the sigmoid function for backPropagation
func sigmoidPrime(x float64) float64 {
	return x * (1.0 - x)
}

// NewNetwork initilizes a new neural network.
func newNetwork(config neuralNetConfig) *neuralNet {
	return &neuralNet{config: config}
}

// Train traings a neural network using backPropagation
func (nn *neuralNet) train(x, y *mat.Dense) error {

	// Initialize biases/weights.
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	wHiddenRaw := make([]float64,
		nn.config.hiddenNeurons*nn.config.inputNeurons)
	bHiddenRaw := make([]float64, nn.config.hiddenNeurons)
	wOutRaw := make([]float64,
		nn.config.outputNeurons*nn.config.hiddenNeurons)
	bOutRaw := make([]float64, nn.config.outputNeurons)

	for _, param := range [][]float64{wHiddenRaw, bHiddenRaw, wOutRaw, bOutRaw} {

		for i := range param {
			param[i] = randGen.Float64()
		}
	}

	wHidden := mat.NewDense(nn.config.inputNeurons, nn.config.hiddenNeurons, wHiddenRaw)
	bHidden := mat.NewDense(1, nn.config.hiddenNeurons, bHiddenRaw)
	wOut := mat.NewDense(nn.config.hiddenNeurons, nn.config.outputNeurons, wOutRaw)
	bOut := mat.NewDense(1, nn.config.outputNeurons, bOutRaw)

	// Define the output of the neural network.
	var output mat.Dense
	//fmt.Println(output)

	// Loop over the number of epochs utilizing
	// backpropagration to train our model.
	for i := 0; i < nn.config.numEpochs; i++ {

		var hiddenLayerInput mat.Dense
		//fmt.Println(hiddenLayerInput)
		hiddenLayerInput.Mul(x, wHidden)
		addBHidden := func(_, col int,
			v float64) float64 {
			return v + bHidden.At(0, col)
		}
		hiddenLayerInput.Apply(addBHidden, &hiddenLayerInput)
		var hiddenLayerActivations mat.Dense
		applySigmoid := func(_, _ int, v float64) float64 {
			return sigmoid(v)
		}
		hiddenLayerActivations.Apply(applySigmoid,
			&hiddenLayerInput)

		var outputLayerInput mat.Dense
		outputLayerInput.Mul(&hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 {
			return v + bOut.At(0, col)
		}
		outputLayerInput.Apply(addBOut, &outputLayerInput)
		output.Apply(applySigmoid, &outputLayerInput)

		//
		//	888888b.                  888     	8888888b.                                                888   d8b
		//	888  "88b                 888     	888   Y88b                                               888   Y8P
		//	888  .88P                 888     	888    888                                               888
		//	8888888K.  8888b.  .d8888b888  888	888   d88P888d888 .d88b. 88888b.  8888b.  .d88b.  8888b. 888888888 .d88b. 88888b.
		//	888  "Y88b    "88bd88P"   888 .88P	8888888P" 888P"  d88""88b888 "88b    "88bd88P"88b    "88b888   888d88""88b888 "88b
		//	888    888.d888888888     888888K 	888       888    888  888888  888.d888888888  888.d888888888   888888  888888  888
		//	888   d88P888  888Y88b.   888 "88b	888       888    Y88..88P888 d88P888  888Y88b 888888  888Y88b. 888Y88..88P888  888
		//	8888888P" "Y888888 "Y8888P888  888	888       888     "Y88P" 88888P" "Y888888 "Y88888"Y888888 "Y888888 "Y88P" 888  888
		//	                                  	                         888                  888
		//	                                  	                         888             Y8b d88P
		//	                                  	                         888              "Y88P"

		var networkError mat.Dense
		networkError.Sub(y, &output)

		var slopeOutputLayer mat.Dense
		applySigmoidPrime := func(_, _ int, v float64) float64 {
			return sigmoidPrime(v)
		}
		slopeOutputLayer.Apply(applySigmoidPrime, &output)
		var slopeHiddenLayer mat.Dense
		slopeHiddenLayer.Apply(applySigmoidPrime, &hiddenLayerActivations)

		var dOutput mat.Dense
		dOutput.MulElem(&networkError, &slopeOutputLayer)
		var errorAtHiddenLayer mat.Dense
		errorAtHiddenLayer.Mul(&dOutput, wOut.T())

		var dHiddenLayer mat.Dense
		dHiddenLayer.MulElem(&errorAtHiddenLayer, &slopeHiddenLayer)

		//
		//         d8888     888  d8b                888   		                                                    888
		//        d88888     888  Y8P                888   		                                                    888
		//       d88P888     888                     888   		                                                    888
		//      d88P 888 .d88888 8888888  888.d8888b 888888		88888b.  8888b. 888d888 8888b. 88888b.d88b.  .d88b. 888888 .d88b. 888d888.d8888b
		//     d88P  888d88" 888 "888888  88888K     888   		888 "88b    "88b888P"      "88b888 "888 "88bd8P  Y8b888   d8P  Y8b888P"  88K
		//    d88P   888888  888  888888  888"Y8888b.888   		888  888.d888888888    .d888888888  888  88888888888888   88888888888    "Y8888b.
		//   d8888888888Y88b 888  888Y88b 888     X88Y88b. 		888 d88P888  888888    888  888888  888  888Y8b.    Y88b. Y8b.    888         X88
		//  d88P     888 "Y88888  888 "Y88888 88888P' "Y888		88888P" "Y888888888    "Y888888888  888  888 "Y8888  "Y888 "Y8888 888     88888P'
		//		   				  888                      		888
		//		   			     d88P                      		888
		//		   			   888P"                       		888

		var wOudAdj mat.Dense
		wOudAdj.Mul(hiddenLayerActivations.T(), &dOutput)
		wOudAdj.Scale(nn.config.learningRate, &wOudAdj)
		wOut.Add(wOut, &wOudAdj)

		bOutAdj, err := sumAlongAxis(0, &dOutput)
		if err != nil {
			log.Fatal(err)
		}
		bOutAdj.Scale(nn.config.learningRate, bOutAdj)
		bOutAdj.Add(bOut, bOutAdj)

		var wHiddenAdj mat.Dense
		wHiddenAdj.Mul(x.T(), &dHiddenLayer)
		wHiddenAdj.Scale(nn.config.learningRate, &wHiddenAdj)
		wHidden.Add(wHidden, &wHiddenAdj)

		bHiddenAdj, err := sumAlongAxis(0, &dHiddenLayer)
		if err != nil {
			return err
		}
		//fmt.Println(bHiddenAdj)
		bHiddenAdj.Scale(nn.config.learningRate, bHiddenAdj)
		bHidden.Add(bHidden, bHiddenAdj)

	}

	nn.wHidden = wHidden
	nn.bHidden = bHidden
	nn.wOut = wOut
	nn.bOut = bOut

	return nil
}

func sumAlongAxis(axis int, m *mat.Dense) (*mat.Dense, error) {

	numRows, numCols := m.Dims()

	var output *mat.Dense

	switch axis {
	case 0:
		data := make([]float64, numCols)
		for i := 0; i < numCols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		output = mat.NewDense(1, numCols, data)
	case 1:
		data := make([]float64, numRows)
		for i := 0; i < numRows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		output = mat.NewDense(numRows, 1, data)
	default:
		return nil, errors.New("invalid axis, must be 0 or 1")
	}

	return output, nil
}

func main() {

	// An empty matrix is one that has zero size.
	// Empty matrices are used to allow the destination
	// of a matrix operation to assume the correct size
	// automatically. This operation will re-use the
	// backing data, if available, or will allocate new data
	// if necessary. The IsEmpty method returns whether the
	// given matrix is empty. The zero-value of a matrix is
	// empty, and is useful for easily getting the result of
	// matrix operations.

	/*
		var c mat.Dense // construct a new zero-value matrix

		c.Mul(a, a) // c is automatically adjusted to be the right size

	*/

	/*// Define our input attributes.
	input := mat.NewDense(3, 4, []float64{
		1.0, 0.0, 1.0, 0.0,
		1.0, 0.0, 1.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
	})

	// Define our labels
	labels := mat.NewDense(3, 1, []float64{1.0, 1.0, 0.0})

	// Define our network architecture and
	// learning parameters.
	config := neuralNetConfig{
		inputNeurons:  4,
		outputNeurons: 1,
		hiddenNeurons: 3,
		numEpochs:     5000,
		learningRate:  0.3,
	}

	// Train the nwueal network.
	network := newNetwork(config)
	if err := network.train(input, labels); err != nil {
		log.Fatal(err)
	}

	// output the weights that define our network!
	f := mat.Formatted(network.wHidden, mat.Prefix(" "), mat.Excerpt(0))
	fmt.Printf("\nwHidden = \n%v \n\n", f)

	f = mat.Formatted(network.bHidden, mat.Prefix(" "), mat.Excerpt(0))
	fmt.Printf("\nbHidden = \n%v \n\n", f)

	f = mat.Formatted(network.wOut, mat.Prefix(" "), mat.Excerpt(0))
	fmt.Printf("\nwOut = \n%v \n\n", f)

	f = mat.Formatted(network.bOut, mat.Prefix(" "), mat.Excerpt(0))
	fmt.Printf("\nbOut = %v \n\n", f)

	fmt.Printf("A matrix: \n%v\n\n", mat.Formatted(input, mat.Prefix(" "), mat.Excerpt(0)))
	*/

	file, err := os.Open(trainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new CSV reader reading form the opened file.
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 7

	// Reading all the data
	rawData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// inputsData and labelsDAta will hold all the
	// float values taht will eventually be
	// used to form our matrices.

	inputData := make([]float64, 4*len(rawData))
	labelsData := make([]float64, 3*len(rawData))

	// InputsIndex will track the current index of
	// inputs matrix values.
	var inputsIndex int
	var labelsIndex int

	// Sequentially move the rows into a slice of floats.
	for idx, record := range rawData {

		// Wkiping the header row
		if idx == 0 {
			continue
		}

		// Loop over the float columns.
		for i, val := range record {

			// covnert the value to a float.
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Add to the labelsData if relevant
			if i == 4 || i == 5 || i == 6 {
				labelsData[labelsIndex] = parsedVal
				labelsIndex++
				continue
			}

			// Add the float value to the slice of floats.
			inputData[inputsIndex] = parsedVal
			inputsIndex++
		}
	}

	// Form the matrices.
	inputs := mat.NewDense(len(rawData), 4, inputData)
	labels := mat.NewDense(len(rawData), 3, labelsData)

	config := neuralNetConfig{
		inputNeurons:  4,
		outputNeurons: 3,
		hiddenNeurons: 3,
		numEpochs:     5000,
		learningRate:  0.3,
	}

	// Train the neural network
	network := newNetwork(config)
	if err := network.train(inputs, labels); err != nil {
		log.Fatal(err)
	}

	/*TESTING MODEL*/
	testFile, err := os.Open(testFile)
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()
	testReader := csv.NewReader(testFile)
	testReader.FieldsPerRecord = 7

	testData, err := testReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	TinputData := make([]float64, 4*len(testData))
	TlabelsData := make([]float64, 3*len(testData))

	// InputsIndex will track the current index of
	// inputs matrix values.
	inputsIndex = 0
	labelsIndex = 0

	for idx, colRecord := range testData {

		if idx == 0 {
			continue
		}

		for i, col := range colRecord {

			// Parsing values
			parsVal, err := strconv.ParseFloat(col, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Add the labelsData if it's relevant
			if i == 4 || i == 5 || i == 6 {
				TlabelsData[labelsIndex] = parsVal
				labelsIndex++
				continue
			}

			// Add the float value to the slices
			TinputData[inputsIndex] = parsVal
			inputsIndex++
		}
	}

	testInputs := mat.NewDense(len(testData), 4, TinputData)
	testLabels := mat.NewDense(len(testData), 3, TlabelsData)

	// Make the predictions using the trained model.
	predictions, err := network.predict(testInputs)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the accuraccy
	var truePosNeg int
	numPreds, _ := predictions.Dims()
	for i := 0; i < numPreds; i++ {

		// Get the label.
		labelRow := mat.Row(nil, i, testLabels)
		var species int
		for idx, label := range labelRow {
			if label == 1.0 {
				species = idx
				break
			}
		}

		// Accumulate the true positive/negative count.
		if predictions.At(i, species) == floats.Max(mat.Row(nil, i, predictions)) {
			truePosNeg++
		}
	}

	// Calculate te accuracy (subset accuracy).
	accuracy := float64(truePosNeg) / float64(numPreds)

	// output the accuracy value
	fmt.Printf("\nAccuracy = %0.5f\n\n", accuracy)

}

// Predict makes a predictions based on a trained
// neural network
func (nn *neuralNet) predict(x *mat.Dense) (*mat.Dense, error) {

	// Check to make sure that our neuralNet value
	// Represents a trained model.
	if nn.wHidden == nil || nn.wOut == nil || nn.bHidden == nil || nn.bOut == nil {
		return nil, errors.New("the supplied neural net weights and biases are empty")
	}

	// Define the output of the neural network.
	var output mat.Dense

	// Complete the feed forward process.

	var hiddenLayerInput mat.Dense
	hiddenLayerInput.Mul(x, nn.wHidden)
	addBHidded := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidded, &hiddenLayerInput)

	var hiddenLayerActivations mat.Dense
	applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, &hiddenLayerInput)

	var outputLayerInput mat.Dense
	outputLayerInput.Mul(&hiddenLayerActivations, nn.wOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	outputLayerInput.Apply(addBOut, &outputLayerInput)

	output.Apply(applySigmoid, &outputLayerInput)

	return &output, nil
}
