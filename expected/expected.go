package expected

import (
	"math"
	"math/rand"
	"parallel/charts"
	"parallel/dataSet"
	"parallel/efficiency"
	"parallel/nash"

	// "sort"
	"sync"
	"time"
)

type test struct {
	E    float64
	Nash float64
	N    int
}

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

func wotkerCalcMaxMonteCarlo(arrOfA []float64, arrOfB []dataSet.BWithIndexes, results chan<- dataSet.ArrWithIndex, stop <-chan bool, workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()
	num := rand.Intn(int(math.Pow(6, float64(dataSet.NumOfAgents))))
	maxValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(num))
	count := 0
	T := 0.000000000001

	var arrOfE []float64
	var arrOfNash []float64
	var arrOfN []int

	for {
		select {
		default:
			T += 0.000000000001
			count += 1
			newNum := changeRandomNumber(num)
			newValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(newNum))
			currentValue := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(num))

			// Если мы нашли новое большее значение
			if newValue > maxValue {
				maxValue = newValue
				arrOfE = append(arrOfE, newValue)
				arrOfNash = append(arrOfNash, nash.CalcNashCriterion(int64(newNum), arrOfA, arrOfB))
				arrOfN = append(arrOfN, count)
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
			if count%1_000_000 == 0 {
				num = rand.Intn(int(math.Pow(6, float64(dataSet.NumOfAgents))))
			}

		case <-stop:
			results <- dataSet.ArrWithIndex{Index: 0, E: arrOfE, Nash: arrOfNash, Max: maxValue, N: arrOfN}
			return
		}
	}
}

func workerCalcMax(arrOfA []float64, arrOfB []dataSet.BWithIndexes, start, end int, results chan<- dataSet.ArrWithIndex, stop <-chan bool, workerGroup *sync.WaitGroup, workerNumber int) {
	defer workerGroup.Done()
	maxEfficiency := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(start))
	var arrOfE []float64
	var arrOfNash []float64
	var arrOfN []int
	stopChart := start + 1000
	for {
		select {
		default:
			start += 1
			currentEfficiency := efficiency.CalcEfficiency(arrOfA, arrOfB, int64(start))

			if maxEfficiency < currentEfficiency {
				maxEfficiency = currentEfficiency
			}

			if start <= stopChart {
				arrOfE = append(arrOfE, currentEfficiency)
				arrOfNash = append(arrOfNash, nash.CalcNashCriterion(int64(start), arrOfA, arrOfB))
				arrOfN = append(arrOfN, start)
			}
			// Если превысили, то завершаем работу
			if start >= end {
				results <- dataSet.ArrWithIndex{Index: workerNumber, E: arrOfE, Nash: arrOfNash, N: arrOfN, Max: maxEfficiency}
				return
			}
		case <-stop:
			// Если пришло сообщение об остановке, завершаем работу
			results <- dataSet.ArrWithIndex{Index: workerNumber, E: arrOfE, Nash: arrOfNash, N: arrOfN, Max: maxEfficiency}
			return
		}
	}
}

func CalcExpectedMax(arrOfA []float64, arrOfB []dataSet.BWithIndexes) float64 {
	workerGroup := &sync.WaitGroup{}
	results := make(chan dataSet.ArrWithIndex)
	stop := make(chan bool)

	workerGroup.Add(1)
	go wotkerCalcMaxMonteCarlo(arrOfA, arrOfB, results, stop, workerGroup)

	time.Sleep(1 * time.Minute)
	go func() {
		stop <- true
		workerGroup.Wait()
		close(results)
	}()

	var max float64
	for result := range results {
		charts.DrawChartWith2Lines(result.E, result.Nash, result.N, "EAndNash.html", "E и E Нэша", "N", "E", "E Нэша")
		max = result.Max
	}

	return max

	//? Поиск перебором
	// workerGroup := &sync.WaitGroup{}
	// step := int(math.Pow(6, float64(dataSet.NumOfAgents))) / 12
	// results := make(chan dataSet.ArrWithIndex, 12)
	// stop := make(chan bool, 12)

	// for i := 0; i < 12; i++ {
	// 	workerGroup.Add(1)
	// 	go workerCalcMax(arrOfA, arrOfB, step*i, step*(i+1), results, stop, workerGroup, i)
	// }

	// time.Sleep(20 * time.Minute)

	// go func() {
	// 	for i := 0; i < 12; i++ {
	// 		stop <- true
	// 	}
	// 	workerGroup.Wait()
	// 	close(results)
	// }()

	// var arrOfMax []dataSet.ArrWithIndex
	// for result := range results {
	// 	arrOfMax = append(arrOfMax, result)
	// }

	// maxValue := arrOfMax[0].Max
	// for _, value := range arrOfMax {
	// 	if value.Max > maxValue {
	// 		maxValue = value.Max
	// 	}
	// }

	// // Строим график
	// var arrOfE []float64
	// var arrOfNash []float64
	// var arrOfN []int

	// var forSortArr []test

	// for t := 0; t < 12; t++ {
	// 	for _, result := range arrOfMax {
	// 		if t == result.Index {
	// 			arrOfE = append(arrOfE, result.E...)
	// 			arrOfNash = append(arrOfNash, result.Nash...)
	// 			arrOfN = append(arrOfN, result.N...)

	// 		}
	// 	}
	// }

	// for i, e := range arrOfE {
	// 	forSortArr = append(forSortArr, test{E: e, Nash: arrOfNash[i], N: arrOfN[i]})
	// }

	// sort.Slice(forSortArr, func(i, j int) bool {
	// 	return forSortArr[i].Nash > forSortArr[j].Nash // сортировка по убыванию
	// })

	// var arrOfESorted []float64
	// var arrOfNashSorted []float64
	// var arrOfNSorted []int

	// for _, i := range forSortArr {
	// 	arrOfESorted = append(arrOfESorted, i.E)
	// 	arrOfNashSorted = append(arrOfNashSorted, i.Nash)
	// 	arrOfNSorted = append(arrOfNSorted, i.N)
	// }
	// charts.DrawChartWith2Lines(arrOfE, arrOfNash, arrOfN, "EAndNash.html", "E и E Нэша", "N", "E", "E Нэша")

	// charts.DrawChartWith2Lines(arrOfESorted, arrOfNashSorted, arrOfNSorted, "sortedE.html", "E и E Нэша", "N", "E", "E Нэша")
	// return maxValue
}
