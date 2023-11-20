package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Delta struct {
	x     int
	y     int
	value float64
}

// * DATA
var ResultImageBlueAndRed = FillResultImage()
var OriginalImageBlue = CreateAndReadDataFromImage(BLUE_IMAGE, strconv.Itoa(SIZE))
var OriginalImageRed = CreateAndReadDataFromImage(RED_IMAGE, strconv.Itoa(SIZE))
var SumOfResultImage int64 = 0
var SumOfOriginalBlue int64 = int64(CalcN(OriginalImageBlue))
var SumOfOriginalRed int64 = int64(CalcN(OriginalImageRed))
var DeltaBlue int64 = 0
var DeltaRed int64 = 0
var Stop atomic.Bool = atomic.Bool{}

func ResetData() {
	ResultImageBlueAndRed = FillResultImage()
	SumOfResultImage = 0
	DeltaRed = 0
	DeltaBlue = 0
	Stop.Store(false)
}

func CalcDelta(currentN, currentValue, originalValue, originalN int) float64 {
	return float64(originalValue) - float64(originalN)/float64(currentN)*float64(currentValue)
}

func FindMaxDelta(arr []Delta) (int, int) {
	max := arr[0]
	for _, value := range arr {
		if value.value > max.value {
			max = value
		}
	}

	return max.x, max.y
}

func MakeStartPoint(originalN *int64, originalImage *[][]int) (int, int) {
	startX := rand.Intn(SIZE)
	startY := rand.Intn(SIZE)

	ResultImageBlueAndRed[startY][startX] += 1

	atomic.AddInt64(&SumOfResultImage, 1)
	return startX, startY
}

func CalcDeltaInt(originalImage [][]int, originalN, delta *int64) int64 {
	var sum float64 = 0
	for y, line := range ResultImageBlueAndRed {
		for x, value := range line {
			sum += math.Abs(float64(originalImage[y][x]) - float64(*originalN)/float64(atomic.LoadInt64(&SumOfResultImage))*float64(value))
		}
	}
	atomic.StoreInt64(delta, int64(sum))
	return int64(sum)
}

func calcDeltas(x, y, resultN, originalN int, resultImage, originalImage [][]int) (int, int) {
	var deltas []Delta

	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			if i < 0 || j < 0 || i > SIZE-1 || j > SIZE-1 || (i == y && j == x) {
				continue
			}
			deltas = append(deltas, Delta{j, i, CalcDelta(resultN, resultImage[i][j], originalImage[i][j], originalN)})
		}
	}
	return FindMaxDelta(deltas)
}

func workerWalk(x, y, step, numOfWorker int, originalNAddr *int64, originalImageAddr *[][]int, workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()

	count := 0
	originalN := atomic.LoadInt64(originalNAddr)
	originalImage := *originalImageAddr
	resultImage := ResultImageBlueAndRed

	var resultN int64

	for {
		count++

		// Stop
		if Stop.Load() {
			return
		}

		// Make step
		if count%1e4 == 0 {
			newDeltaBlue := CalcDeltaInt(OriginalImageBlue, &SumOfOriginalBlue, &DeltaBlue)
			newDeltaRed := CalcDeltaInt(OriginalImageRed, &SumOfOriginalRed, &DeltaRed)

			if float64(newDeltaBlue)/float64(newDeltaRed) < 0.2 {
				Stop.Store(true)
				// fmt.Println("ZA WARUDO")
				return
			}
		}

		resultImage[y][x] += step

		resultN = atomic.AddInt64(&SumOfResultImage, int64(step))
		// New cords
		x, y = calcDeltas(x, y, int(resultN), int(originalN), resultImage, originalImage)
	}
}

func Walk(numOfWorkers int) {
	workerGroup := &sync.WaitGroup{}

	x, y := MakeStartPoint(&SumOfOriginalRed, &OriginalImageRed)

	workerGroup.Add(1)
	go workerWalk(x, y, 1, -1, &SumOfOriginalRed, &OriginalImageRed, workerGroup)

	for i := 0; i < numOfWorkers; i++ {
		workerGroup.Add(1)
		x, y = MakeStartPoint(&SumOfOriginalBlue, &OriginalImageBlue)
		go workerWalk(x, y, 5, i, &SumOfOriginalBlue, &OriginalImageBlue, workerGroup)
	}

	workerGroup.Wait()
}

func WalkManyTimes() {
	var averageTimes []float64

	for numOfWorkers := 1; numOfWorkers < 25; numOfWorkers++ {
		var times []int
		for numOfTry := 0; numOfTry < 5; numOfTry++ {
			startTime := time.Now()
			Walk(numOfWorkers)
			endTime := time.Now()
			ResetData()
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
		numsOfWorkers = append(numsOfWorkers, i+2)
	}

	DrawChartWith2Lines(averageTimes, idealTimes, numsOfWorkers, "time.html", "Время, мс", "N", "Полученное время", "Идевальное время")
}
