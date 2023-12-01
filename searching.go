package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var UniqueMaxCount int64 = 0
var DotsSet *SetWithDots = &SetWithDots{value: make(map[Dot]struct{})}

type FuncValueWithCoords struct {
	x     int
	y     int
	value float64
}

func CalcFunction(x, y float64) (result float64) {
	return math.Sin(x+y) + math.Sin(math.Pi*(x+2*y)) + math.Sin(math.Sqrt2*(2*x+y))
}

func CreateArrayOfCounts() [][]int {
	size := SIZE * int(math.Pow(DISCRETIZATION, -1))
	var result [][]int
	for y := 0; y < size; y++ {
		var line []int
		for x := 0; x < size; x++ {
			line = append(line, 0)
		}
		result = append(result, line)
	}

	return result
}

func findMaxFunc(slice [][]FuncValueWithCoords) (int, int) {
	max := FuncValueWithCoords{x: -1, y: -1, value: -math.Inf(1)}

	for _, line := range slice {
		for _, value := range line {

			if max.value < value.value {
				max = value
			}
		}
	}

	return max.x, max.y
}

func checkIfItMax(slice [][]FuncValueWithCoords, currentValue float64) bool {
	for _, line := range slice {
		for _, value := range line {
			if value.value > currentValue {
				return false
			}
		}
	}

	return true
}

func determineWhereToGo(x, y int) (newX int, newY int) {
	var coords [][]FuncValueWithCoords
	for i := y - 1; i <= y+1; i++ {
		var line []FuncValueWithCoords
		for j := x - 1; j <= x+1; j++ {
			if i < 0 || j < 0 || i > SIZE*int(math.Pow(DISCRETIZATION, -1))-1 || j > SIZE*int(math.Pow(DISCRETIZATION, -1))-1 || (j == x && i == y) {
				continue
			}

			line = append(line, FuncValueWithCoords{x: j, y: i, value: CalcFunction(float64(j)*DISCRETIZATION+DISCRETIZATION/2, float64(i)*DISCRETIZATION+DISCRETIZATION/2)})
		}
		coords = append(coords, line)
	}

	if checkIfItMax(coords, CalcFunction(float64(x)*DISCRETIZATION+DISCRETIZATION/2, float64(y)*DISCRETIZATION+DISCRETIZATION/2)) {

		DotsSet.AddDot(x, y)
		return rand.Intn(SIZE * int(math.Pow(DISCRETIZATION, -1))), rand.Intn(SIZE * int(math.Pow(DISCRETIZATION, -1)))
	} else {
		return findMaxFunc(coords)
	}
}

func workerWalking(workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()
	x := rand.Intn(SIZE * int(math.Pow(DISCRETIZATION, -1)))
	y := rand.Intn(SIZE * int(math.Pow(DISCRETIZATION, -1)))

	for {
		x, y = determineWhereToGo(x, y)

		if atomic.LoadInt64(&UniqueMaxCount) >= NEED_TO_FIND {
			return
		}
	}

}

func Walking() {
	var averageTimes []float64

	for numOfWorkers := 1; numOfWorkers <= 24; numOfWorkers++ {
		var times []int

		for numOfTry := 0; numOfTry < 5; numOfTry++ {
			workerGroup := &sync.WaitGroup{}
			startTime := time.Now()

			for i := 0; i < numOfWorkers; i++ {
				workerGroup.Add(1)
				go workerWalking(workerGroup)
			}

			workerGroup.Wait()
			endTime := time.Now()

			DotsSet.Clear()

			fmt.Println(numOfWorkers, numOfTry)
			times = append(times, int(endTime.Sub(startTime)/time.Millisecond))
		}

		averageTime := 0.0

		for _, time := range times {
			averageTime += float64(time)
		}

		averageTime /= float64(len(times))

		averageTimes = append(averageTimes, averageTime)
	}

	var idealTimes []float64
	var numsOfWorkers []int

	for i := range averageTimes {
		idealTimes = append(idealTimes, averageTimes[0]/float64(i+1))
		numsOfWorkers = append(numsOfWorkers, i+1)
	}

	DrawChartWith2Lines(averageTimes, idealTimes, numsOfWorkers, "time.html", "Время, мс", "N", "Полученное время", "Идевальное время")
}

func SoloWalking() {
	workerGroup := &sync.WaitGroup{}
	startTime := time.Now()

	for i := 0; i < 12; i++ {
		fmt.Println(123)
		workerGroup.Add(1)
		go workerWalking(workerGroup)
	}

	workerGroup.Wait()
	endTime := time.Now()

	fmt.Println(endTime.Sub(startTime))
}
