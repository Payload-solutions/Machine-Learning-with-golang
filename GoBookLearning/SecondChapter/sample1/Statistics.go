package main

import (
	"fmt"

	//"gonum.org/v1/gonum/blas/blas64"
	//"github.com/gonum/blas/blas64"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func banner(word string) {
	fmt.Println()
	fmt.Println(word)

}

func VectorsInit() {
	var vector []float64

	vector = append(vector, 11.0)
	vector = append(vector, 25.6565)

	fmt.Println(vector)

	//creating a new vector slice
	banner("Creating a new Vector slice")
	myVector := mat.NewVecDense(2, []float64{11.0, 5.2})
	fmt.Println(myVector)

	// Vector operations
	banner("Operatios")
	vectorA := []float64{11., 5.2, -1.3}
	vectorB := []float64{-7.2, 4.2, 5.1}

	dotProduct := floats.Dot(vectorA, vectorB)

	fmt.Printf("The dot product is %v\n", dotProduct)

	banner("Scale elements of A by 1.5")

	floats.Scale(1.5, vectorA)

	fmt.Printf("Scaling A by 1.5 %v\n", vectorA)

	normB := floats.Norm(vectorB, 2)

	fmt.Printf("The norm/length of B is %v\n", normB)

}

func vectorsInit2() {
	banner("initialize a couple of 'vectors' represented as slices")
	vectorA := mat.NewVecDense(3, []float64{11.0, 5.2, -1.3})
	vectorB := mat.NewVecDense(3, []float64{-7.2, 4.2, 5.1})

	banner("Compute the dot product of A and B")
	dotProduct := mat.Dot(vectorA, vectorB)
	fmt.Printf("The dot product of A and B is : %0.2f\n", dotProduct)

	vectorA.ScaleVec(1.5, vectorA)
	fmt.Printf("The scale of 1.5 from vector A is %v\n", vectorA)

	normB := blas64.Nrm2(vectorB.RawVector())

	fmt.Printf("the norm/length of vector B is: %v\n", normB)

}

func main() {
	//VectorsInit()
	vectorsInit2()
}
