package main

import (
	"fmt"
	"parallel/dataSet"
	"parallel/expected"
	"parallel/monteCarlo"
)

func main() {
	fmt.Println("Start")

	arrOfA, arrOfB := dataSet.CreateStaticDataSet()
	expectedMax := expected.CalcExpectedMax(arrOfA, arrOfB)
	monteCarlo.FindMaxByMonteCarlo(arrOfA, arrOfB, expectedMax)

	fmt.Println("Done")
}
