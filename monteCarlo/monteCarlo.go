package monteCarlo

import (
	"fmt"
	"math"
	"math/rand"
	"parallel/charts"
	"parallel/dataSet"
	"parallel/efficiency"
	"sync"
	"time"
)

func changeRandomNumber(num int) int {
	index := rand.Intn(int(dataSet.NumOfAgents))
	if num/int(math.Pow(6, float64(index)))%6 == 5 {
		return num - int(math.Pow(6, float64(index)))
	}
	if num/int(math.Pow(6, float64(index)))%6 == 0 {
		return num + int(math.Pow(6, float64(index)))
	}

	plusOrMinus := rand.Intn(2)
	if plusOrMinus == 0 {
		return num + int(math.Pow(6, float64(index)))
	}
	return num - int(math.Pow(6, float64(index)))
}

func findMax(arrOfA []float64, arrOfB []dataSet.BWithIndexes, expectedMax float64, results chan<- float64, stop chan bool, workerGroup *sync.WaitGroup, numOfWorkers int) {
	defer workerGroup.Done()
	num := rand.Intn(int(math.Pow(6, float64(dataSet.NumOfAgents))))
	maxValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(num))
	count := 0
	isEnd := false
	T := 0.000000000001

	for {
		select {
		default:
			T += 0.000000000001
			count += 1
			newNum := changeRandomNumber(num)
			newValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(newNum))
			currentValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(num))
			// fmt.Println(count)
			// Если мы нашли новое большее значение
			if newValue > maxValue {
				maxValue = newValue
			}

			// Если мы нашли максимум или больше максимума
			if newValue >= expectedMax {
				isEnd = true
			}

			// Бродим еще чтобы потенциально найти еще больший максимум
			if isEnd && count > 100_000 {
				results <- maxValue
				for i := 1; i < numOfWorkers; i++ {
					stop <- true
				}
				return
			}

			// Если значение больше предыдущего то принимаем если нет то разыгрываем вероятность
			if newValue-currentValue > 0 {
				num = newNum
			} else {
				p := math.Exp(newValue - currentValue/T)
				r := rand.Float64()

				if r <= p {
					num = newNum
				}
			}
			// Если сделали миллион шагов, то берем рандомное число
			if count > 1_000_000 {
				num = rand.Intn(int(math.Pow(6, float64(dataSet.NumOfAgents))))
				count = 0
			}

		case <-stop:
			results <- maxValue
			return
		}
	}
}

func FindMaxByMonteCarlo(arrOfA []float64, arrOfB []dataSet.BWithIndexes, expectedMax float64) {
	var averageTimes []float64
	for numOfWorkers := 1; numOfWorkers <= 24; numOfWorkers++ {
		var times []int
		for i := 0; i < 10; i++ {
			workerGroup := &sync.WaitGroup{}
			results := make(chan float64)
			stop := make(chan bool)

			startTime := time.Now()
			for j := 0; j < numOfWorkers; j++ {
				workerGroup.Add(1)
				go findMax(arrOfA, arrOfB, expectedMax, results, stop, workerGroup, numOfWorkers)
			}

			go func() {
				workerGroup.Wait()
				close(results)
			}()

			for result := range results {
				fmt.Println(result)
			}

			endTime := time.Now()

			times = append(times, int(endTime.Sub(startTime)/time.Millisecond))
		}
		timeAverage := 0
		for _, time := range times {
			timeAverage += time
		}

		timeAverage /= len(times)
		averageTimes = append(averageTimes, float64(timeAverage))
	}
	var idealTimes []float64
	var numsOfWorkers []int
	for i := range averageTimes {
		idealTimes = append(idealTimes, averageTimes[0]/float64(i+1))
		numsOfWorkers = append(numsOfWorkers, i+1)
	}

	charts.DrawChartWith2Lines(averageTimes, idealTimes, numsOfWorkers, "time.html", "Время, мс", "N", "Полученное время", "Идевальное время")
}
