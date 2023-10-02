package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func calcFunc(value float64, roots []float64) float64 {
	var accumulator float64 = 1

	for i := 0; i < 4; i++ {
		accumulator *= value - roots[i]
	}

	return accumulator
}

func calcSpace(startWith, endWith float64, numSteps int64, roots []float64, workerGroup *sync.WaitGroup, results chan<- float64) {
	defer workerGroup.Done()

	delta := (endWith - startWith) / float64(numSteps)
	var sum float64 = 0

	for startWith < endWith {
		sum += delta * (calcFunc(startWith, roots) + calcFunc(startWith+delta, roots)) / 2
		startWith += delta
	}

	results <- sum
}

func main() {
	start, end := 0.0, 10.0         // Границы интегирования
	var numSteps int64 = 2000000000 // Количество шагов
	roots := []float64{1, 4, 9, 10} // Корни уравнения

	file, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// numWorkers := 3 // Количество потоков

	for numWorkers := 1; numWorkers <= 24; numWorkers++ {
		var attemptsResults []float64
		var attemptsTime []int64

		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close()

		for i := 0; i < 5; i++ {
			workerGroup := &sync.WaitGroup{}
			results := make(chan float64, numWorkers)
			step := (end - start) / float64(numWorkers)

			startTime := time.Now() // Начало отсчета

			for i := 0; i < numWorkers; i++ {
				workerGroup.Add(1)
				go calcSpace(start+float64(i)*step, start+float64(i+1)*step, numSteps/int64(numWorkers), roots, workerGroup, results)
			}

			go func() {
				workerGroup.Wait()
				close(results)
			}()

			totalSpace := 0.0
			for area := range results {
				totalSpace += area
			}

			endTime := time.Now() // Конец отсчета

			attemptsTime = append(attemptsTime, int64(endTime.Sub(startTime)/time.Millisecond))
			attemptsResults = append(attemptsResults, totalSpace)
		}

		var result float64 = 0.0
		for _, res := range attemptsResults {
			result += res
		}
		result /= float64((len(attemptsResults)))

		var resultTime int64 = 0
		for _, res := range attemptsTime {
			resultTime += res
		}
		resultTime /= int64(len(attemptsTime))

		file.WriteString(strconv.FormatInt(int64(numWorkers), 10) + "\t")
		file.WriteString(strconv.FormatFloat(result, 'f', -1, 32) + "\t")
		file.WriteString(strconv.FormatInt(resultTime, 10) + "\n")
	}

}
