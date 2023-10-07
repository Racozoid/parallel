package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func translateToPolar(x, y float64) (float64, float64) {
	radius := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	angle := math.Acos(x / math.Sqrt(math.Pow(x, 2)+math.Pow(y, 2)))

	if y < 0 {
		angle *= -1
	}

	if x == 0 && y == 0 {
		radius = 0
		angle = 0
	}

	return radius, angle
}

func calcFunc(angle float64) float64 {
	return math.Sin(5*angle) - math.Cos(10*angle)/math.Cos(2*angle) - math.Cos(6*angle) + 7
	// return 2 * math.Pow(math.Cos(angle), 2)
	// return 1
}

func isInside(x, y float64) bool {
	radius, angle := translateToPolar(x, y)
	funcRadius := calcFunc(angle)

	if radius <= funcRadius {
		return true
	} else {
		return false
	}
}

func calcNumberInside(min, max float64, numOfRandomNumbers int, workerGroup *sync.WaitGroup, results chan<- int64) {
	defer workerGroup.Done()

	var countSuccessTries int64 = 0

	for i := 0; i <= numOfRandomNumbers; i++ {
		randomX := min + rand.Float64()*(max-min)
		randomY := min + rand.Float64()*(max-min)

		if isInside(randomX, randomY) {
			countSuccessTries++
		}
	}
	results <- countSuccessTries
}

func main() {
	var min, max float64 = -10.0, 10.0 // Минимальное и максимальное значение для генерации
	numOfRandomNumbers := 50_000_000   // Колличество случайных чисел

	file, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Unable to read file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Расчет при количестве потоков от 1 до 24
	for numOfWorkers := 1; numOfWorkers <= 24; numOfWorkers++ {
		var attemptResults []int64
		var attemptTime []int64
		numbersForWorker := numOfRandomNumbers / numOfWorkers

		// Пять замеров времени
		for n := 1; n < 5; n++ {
			workerGroup := &sync.WaitGroup{}
			results := make(chan int64)

			startTime := time.Now() // Начало отсчета времени

			for i := 0; i < numOfWorkers; i++ {
				workerGroup.Add(1)
				go calcNumberInside(min, max, numbersForWorker, workerGroup, results)
			}

			go func() {
				workerGroup.Wait()
				close(results)
			}()

			var totalInside int64 = 0
			for num := range results {
				totalInside += num
			}

			endTime := time.Now() // Конец отсчета времени

			attemptResults = append(attemptResults, totalInside)
			attemptTime = append(attemptTime, int64(endTime.Sub(startTime)/time.Millisecond))
		}

		// Нахождение среднего от попыток
		var averageNumInside int64 = 0
		for _, num := range attemptResults {
			averageNumInside += num
		}
		averageNumInside /= int64(len(attemptResults))

		var averageTime int64 = 0
		for _, time := range attemptTime {
			averageTime += time
		}
		averageTime /= int64(len(attemptTime))

		file.WriteString(strconv.FormatInt(int64(numOfWorkers), 10) + "\t")
		file.WriteString(strconv.FormatFloat(math.Pow(max-min, 2)*float64(averageNumInside)/float64(numOfRandomNumbers), 'f', 3, 64) + "\t")
		file.WriteString(strconv.FormatInt(averageTime, 10) + "\n")
	}

}
