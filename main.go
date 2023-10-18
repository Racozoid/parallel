package main

import (
	"fmt"
	"math"
	"math/rand"

	"os"

	// "sort"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type BWithIndexes struct {
	i     int
	j     int
	value float64
}

type ValuesFloat64WithNumber struct {
	number int
	value  []float64
}

func createLineChart(data []opts.LineData, XAxis []float64, filename, title, XTitle, YTitle string) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: XTitle,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: YTitle,
		}),
	)

	line.SetXAxis(XAxis).AddSeries(YTitle, data).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(filename)
	_ = line.Render(f)
}

func create2LineChart(data1, data2 []opts.LineData, XAxis []int, filename, title, XTitle, YTitle1, YTitle2 string) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: XTitle,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: title,
		}),
	)

	line.SetXAxis(XAxis).
		AddSeries(YTitle1, data1).
		AddSeries(YTitle2, data2).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(filename)
	_ = line.Render(f)
}

func calcEfficiency(num, numOfVarinats int, arrOfA []float64, arrOfB []BWithIndexes) float64 {
	var E float64 = 0.0
	var numBinary string = fmt.Sprintf("%0*b", numOfVarinats, num)

	var step int = numOfVarinats - 1
	startPosition := 0

	for i := 0; i < len(arrOfA); i++ {

		if string(numBinary[i]) == "0" {
			E -= arrOfA[i]
		} else {
			E += arrOfA[i]
		}

		for j := startPosition; j < startPosition+step; j++ {
			if string(numBinary[arrOfB[j].i]) == "0" {
				arrOfB[j].value *= -1
			}

			if string(numBinary[arrOfB[j].j]) == "0" {
				arrOfB[j].value *= -1
			}

			E += arrOfB[j].value
		}
		startPosition += step
		step--
	}
	return E
}

func calcAllEfficiency(number, start, end, numOfVariants int, arrOfA []float64, arrOfB []BWithIndexes, workerGroup *sync.WaitGroup, results chan<- ValuesFloat64WithNumber) {
	defer workerGroup.Done()

	var arrOfE []float64
	for i := start; i < end; i++ {
		arrOfE = append(arrOfE, calcEfficiency(i, numOfVariants, arrOfA, arrOfB))
	}

	results <- ValuesFloat64WithNumber{number, arrOfE}
}

func calcAverageAMax(number int, start, end float64, numOfVariants int, arrOfE []float64, indexAMax int, workerGroup *sync.WaitGroup, results chan<- ValuesFloat64WithNumber) {
	defer workerGroup.Done()

	var arrAverageAMax []float64
	for T := start; T < end; T += 0.5 {
		//? Подсчет Ро
		var arrOfRho []float64

		for _, E := range arrOfE {
			arrOfRho = append(arrOfRho, math.Exp(E/T))
		}

		var sumOfRho float64 = 0

		for _, rho := range arrOfRho {
			sumOfRho += rho
		}

		//? Нормализация Ро
		var arrOfRhoNormalized []float64

		for _, rho := range arrOfRho {
			arrOfRhoNormalized = append(arrOfRhoNormalized, rho/sumOfRho)
		}

		//? Поиск <Amax>
		AverageAMax := 0.0
		for i, Rho := range arrOfRhoNormalized {
			var numBinary string = fmt.Sprintf("%0*b", numOfVariants, i)

			if string(numBinary[indexAMax]) == "0" {
				AverageAMax -= Rho
			} else {
				AverageAMax += Rho
			}
		}

		arrAverageAMax = append(arrAverageAMax, AverageAMax)
	}

	results <- ValuesFloat64WithNumber{number, arrAverageAMax}
}

func findIndexAMax(arrOfA []float64) int {
	var indexAMax int = -1
	var AMax float64 = -2
	for i, a := range arrOfA {
		if a > AMax {
			AMax = a
			indexAMax = i
		}
	}

	return indexAMax
}

func main() {
	var numOfVariants int = 22 //? Сколько вариант

	var arrOfA []float64
	var arrOfB []BWithIndexes

	//? Генерация коэффициентов a и b
	for i := 0; i < numOfVariants; i++ {
		arrOfA = append(arrOfA, float64(-1+rand.Float64()*(1+1)))
	}

	for i := 0; i < numOfVariants; i++ {
		for j := i + 1; j < numOfVariants; j++ {
			arrOfB = append(arrOfB, BWithIndexes{i, j, float64(-1 + rand.Float64()*(1+1))})
		}
	}
	indexAMax := findIndexAMax(arrOfA)

	var arrOfTimes []time.Duration
	var arrOfNumberOfWorkers []int

	for numWorkers := 1; numWorkers <= 24; numWorkers++ {
		//? Для графика E
		// arrY := make([]opts.LineData, 0)
		// var arrX []float64

		//? Для графика <А макс>
		// arrAverageAMax := make([]opts.LineData, 0)
		// var arrOfT []float64

		workerGroup := &sync.WaitGroup{}
		results := make(chan ValuesFloat64WithNumber)
		numbersForWorker := int(math.Pow(2, float64(numOfVariants)) / float64(numWorkers))

		workerGroupRho := &sync.WaitGroup{}
		resultsRho := make(chan ValuesFloat64WithNumber)
		var TForWorker float64 = 50.0 / float64(numWorkers)

		//*************************** Начало отсчета времени
		startTime := time.Now()

		for i := 0; i < numWorkers; i++ {
			workerGroup.Add(1)
			go calcAllEfficiency(i, i*numbersForWorker, (i+1)*numbersForWorker, numOfVariants, arrOfA, arrOfB, workerGroup, results)
		}

		go func() {
			workerGroup.Wait()
			close(results)
		}()

		var chaoticE []ValuesFloat64WithNumber
		for result := range results {
			chaoticE = append(chaoticE, result)
		}

		var arrOfE []float64
		for t := 0; t < numWorkers; t++ {
			for _, result := range chaoticE {
				// fmt.Print(t, result.number == t, "\t")
				if result.number == t {

					arrOfE = append(arrOfE, result.value...)
				}
			}
		}
		// fmt.Print(len(arrOfE), "\n\n")
		//? График E
		// sortedArrOfE := append([]float64{}, arrOfE...)
		// sort.Sort(sort.Reverse(sort.Float64Slice(sortedArrOfE)))
		// for x, y := range sortedArrOfE {
		// 	if x%100 == 0 {
		// 		arrX = append(arrX, float64(x))
		// 		arrY = append(arrY, opts.LineData{Value: y})
		// 	}
		// }
		// createLineChart(arrY, arrX, "E.html", "E", "n", "E")

		for i := 0; i < numWorkers; i++ {
			workerGroupRho.Add(1)
			go calcAverageAMax(i, 0.1+TForWorker*float64(i), 0.1+TForWorker*float64(i+1), numOfVariants, arrOfE, indexAMax, workerGroupRho, resultsRho)
		}

		go func() {
			workerGroupRho.Wait()
			close(resultsRho)
		}()

		var chaoticAverageAmax []ValuesFloat64WithNumber
		for result := range resultsRho {
			chaoticAverageAmax = append(chaoticAverageAmax, result)
		}

		var arrOfAverageAmax []float64
		for t := 0; t < numWorkers; t++ {
			for _, result := range chaoticAverageAmax {
				if result.number == t {
					arrOfAverageAmax = append(arrOfAverageAmax, result.value...)
				}
			}
		}
		// fmt.Print(len(arrOfAverageAmax), "\n\n")

		// createLineChart(arrAverageAMax, arrOfT, "Amax.html", "<A max>", "T", "<A max>")
		endTime := time.Now() //* Конец отсчета времени
		arrOfTimes = append(arrOfTimes, endTime.Sub(startTime))
		arrOfNumberOfWorkers = append(arrOfNumberOfWorkers, numWorkers)
	}

	experimentTime := make([]opts.LineData, 0)
	for _, t := range arrOfTimes {
		experimentTime = append(experimentTime, opts.LineData{Value: int64(t / time.Millisecond)})
	}

	idealTime := make([]opts.LineData, 0)
	for _, num := range arrOfNumberOfWorkers {
		idealTime = append(idealTime, opts.LineData{Value: int64(arrOfTimes[0]/time.Millisecond) / int64(num)})
	}

	create2LineChart(experimentTime, idealTime, arrOfNumberOfWorkers, "Time.html", "T, ms", "N", "Полученное время", "Идеальное время")

}
