package main

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/mat"
)

// Matrix Samples
func MakingMatrix() {
	data := []float64{1.2, 5.7, 2.4, 7.3}

	a := mat.NewDense(2, 2, data)
	fa := mat.Formatted(a, mat.Prefix(" "))
	fmt.Printf("mat = %v\n\n", fa)

	/*
		mat = ⎡1.2  5.7⎤
				⎣2.4  7.3⎦

	*/

	// we can access and modify certain values withing A via
	// build-in methods:
	// access to position begint for zero as normal index in
	// Informatic.

	val := a.At(0, 1)
	fmt.Println(val)

	// Get the values in speciffic column
	col := mat.Col(nil, 0, a)
	fmt.Printf("The vale in the 1st column is %v\n\n", col)
	// The vale in the 1st column is [1.2 2.4]

	// Get the value in a Kspecific row

	row := mat.Row(nil, 1, a)
	fmt.Printf("The row value is: %v\n\n", row)
	// The row value is: [2.4 7.3]

	// Modify a single element.
	a.Set(0, 1, 11.2)
	fmt.Printf("mat = %v\n\n", fa)
	/*
		mat = ⎡ 1.2  11.2⎤
				⎣ 2.4   7.3⎦
	*/

	// Modify completly a row
	a.SetRow(0, []float64{14.3, -4.2})
	fmt.Printf("mat = %v\n\n", fa)

	// Modify completly a column

	a.SetCol(0, []float64{1.7, -0.3})

	fmt.Printf("mat = %v\n\n", fa)

}

func MatrixOperations() {
	// Create two matrices of the same size, a and b

	a := mat.NewDense(3, 3, []float64{1, 2, 3, 0, 4, 5, 0, 0, 6})
	b := mat.NewDense(3, 3, []float64{8, 9, 10, 1, 4, 2, 9, 0, 2})

	// Now create a third matrix of a different size.

	//c := mat.NewDense(3, 2, []float64{3, 2, 1, 4, 0, 8})

	// Add a and b.

	d := mat.NewDense(0, 0, nil)
	d.Add(a, b)

	fd := mat.Formatted(d, mat.Prefix(" 			"))

	fmt.Println(fd)
	//fmt.Printf("d = a+ b = %0.4v\n\n", fd)
}

func AdavancedFunctions() {

	a := mat.NewDense(3, 3, []float64{1, 2, 3, 0, 4, 5, 0, 0, 6})

	// Compute and output the transpose of the matrix.
	ft := mat.Formatted(a.T(), mat.Prefix(" "))

	fmt.Printf("a^T = %v\n\n", ft)

	// Compute and output the determinant of a.

	determinant := mat.Det(a)
	fmt.Printf("Determinand Det(a) = %0.1f\n\n", determinant)

	// Compute and output the inverse of a.
	aInverse := mat.NewDense(3, 3, nil)
	if err := aInverse.Inverse(a); err != nil {
		log.Fatal(err)
	}
	fi := mat.Formatted(aInverse, mat.Prefix("       "))
	fmt.Printf("a^-1 = %0.4f\n\n", fi)

}

func main() {
	//MakingMatrix()
	//MatrixOperations()
	AdavancedFunctions()
}
