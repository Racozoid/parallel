package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type delta struct {
	x     int
	y     int
	value float64
}

func calcDelta(currentN, currentValue, originalValue, originalN int) float64 {
	return float64(originalValue) - float64(originalN)/float64(currentN)*float64(currentValue)
}

func findMaxDelta(arr []delta) (int, int) {
	max := arr[0]
	for _, value := range arr {
		if value.value > max.value {
			max = value
		}
	}

	return max.x, max.y
}

func makeStartPoint() (int, int) {
	startX := rand.Intn(SIZE)
	startY := rand.Intn(SIZE)

	ResultImage[startY][startX] += 1
	ResultN += 1
	Delta = float64(OriginalN) - float64(OriginalImage[startY][startX])

	return startX, startY
}

func makeStep(x, y int) (int, int) {
	subDelta := PositiveMax(float64(OriginalImage[y][x])-float64(OriginalN)/float64(ResultN+1)*float64(ResultImage[y][x]+1)) - PositiveMax(float64(OriginalImage[y][x])-float64(OriginalN)/float64(ResultN)*float64(ResultImage[y][x]))
	Delta += subDelta
	ResultImage[y][x] += 1
	ResultN += 1

	return x, y
}

func calcDeltas(x, y int) (int, int) {
	var deltas []delta

	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			if i < 0 || j < 0 || i > SIZE-1 || j > SIZE-1 || (i == y && j == x) {
				continue
			}
			deltas = append(deltas, delta{j, i, calcDelta(ResultN, ResultImage[i][j], OriginalImage[i][j], OriginalN)})
		}
	}

	return findMaxDelta(deltas)
}

func StartWalking() {
	startTime := time.Now()
	x, y := makeStartPoint()

	fmt.Println("Delta 0: ", Delta)
	fmt.Println("N0: ", ResultN)

	for Delta > -2e+08 {
		x, y = makeStep(calcDeltas(x, y))

		// Рисование
		if ResultN < 1e+07 && ResultN%1e+05 == 0 {
			WriteDataToFileAndCreateImage(NormalizeImage(ResultImage), "results\\"+strconv.Itoa(ResultN/1e+05)+"_res.png", "results\\_data.txt")
		}
	}

	endTime := time.Now()

	fmt.Println(endTime.Sub(startTime))
	fmt.Println("N: ", ResultN)
	fmt.Println("Delta: ", Delta)
	WriteDataToFileAndCreateImage(NormalizeImage(ResultImage), "res.jpg", "data.txt")
}

func workerWalker(workerGroup *sync.WaitGroup, sharedData *SharedData, i int) {
	defer workerGroup.Done()
	x, y := sharedData.MakeStartPoint()
	originalImage, originalN, resultImage, resultN, deltaRes := sharedData.GetFirstTimeData()
	count := 0
	toResult := FillResultImage()
	toN := 0
	subDelta := 0.0
	for {
		if deltaRes < -2.25e+9 {
			return
		}
		count++

		var deltas []delta

		for i := y - 1; i <= y+1; i++ {
			for j := x - 1; j <= x+1; j++ {
				if i < 0 || j < 0 || i > SIZE-1 || j > SIZE-1 || (i == y && j == x) {
					continue
				}
				deltas = append(deltas, delta{j, i, calcDelta(resultN, resultImage[i][j], originalImage[i][j], originalN)})
			}
		}

		x, y = findMaxDelta(deltas)
		toSubDelta := PositiveMax(float64(originalImage[y][x])-float64(originalN)/float64(resultN+1)*float64(resultImage[y][x]+1)) - math.Abs(float64(originalImage[y][x])-float64(originalN)/float64(resultN)*float64(resultImage[y][x]))

		resultImage[y][x] += 1
		resultN += 1
		deltaRes += toSubDelta

		toResult[y][x] += 1
		toN += 1
		subDelta += toSubDelta

		if count%10_000 == 0 {
			resultImage, resultN, deltaRes = sharedData.MakeManySteps(toResult, toN, subDelta)
			toResult = FillResultImage()
			toN = 0
			subDelta = 0
		}

		if count == 100_000_000 {
			return
		}
	}

}

func StartParallelWalking() {
	var times []float64

	for numOfWorkers := 1; numOfWorkers <= 24; numOfWorkers++ {
		workerGroup := &sync.WaitGroup{}
		sd := SharedData{
			OriginalImage: OriginalImage,
			ResultImage:   FillResultImage(),
			ResultN:       0,
			OriginalN:     CalcN(OriginalImage),
			Delta:         0,
		}

		startTime := time.Now()

		for i := 0; i < numOfWorkers; i++ {
			workerGroup.Add(1)

			go workerWalker(workerGroup, &sd, i)
		}

		workerGroup.Wait()

		endTime := time.Now()

		times = append(times, float64(endTime.Sub(startTime)/time.Millisecond))
		fmt.Println(numOfWorkers, " Done")
	}

	var idealTimes []float64
	var numsOfWorkers []int

	for i := range times {
		idealTimes = append(idealTimes, times[0]/float64(i+1))
		numsOfWorkers = append(numsOfWorkers, i+1)
	}

	DrawChartWith2Lines(times, idealTimes, numsOfWorkers, "times.html", "Время мс", "N", "Полученное время", "Идевальное время")

	workerGroup := &sync.WaitGroup{}
	sd := SharedData{
		OriginalImage: OriginalImage,
		ResultImage:   FillResultImage(),
		ResultN:       0,
		OriginalN:     CalcN(OriginalImage),
		Delta:         0,
	}

	workerGroup.Add(1)
	go workerWalker(workerGroup, &sd, 0)
	workerGroup.Wait()

	WriteDataToFileAndCreateImage(NormalizeImage(sd.ResultImage), "res.jpg", "data.txt")
}
