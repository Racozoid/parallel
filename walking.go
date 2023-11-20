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

type SDWtihMutexSlice struct {
	sync.Mutex
	value [][]int
}

type SDWithRWMutexSlice struct {
	sync.RWMutex
	value [][]int
}

func (s *SDWithRWMutexSlice) GetRWSlice() [][]int {
	s.RLock()
	defer s.RUnlock()

	return s.value
}

func (s *SDWtihMutexSlice) GetSlice() [][]int {
	s.Lock()
	defer s.Unlock()

	return s.value
}

func (s *SDWtihMutexSlice) GetValue(x, y int) int {
	s.Lock()
	defer s.Unlock()

	return s.value[y][x]
}

func (s *SDWtihMutexSlice) AddToSlice(x, y, add int) {
	s.Lock()
	defer s.Unlock()

	s.value[y][x] += add
}

func (s *SDWtihMutexSlice) MakeStartPoint(originalN *int64, originalImage *SDWithRWMutexSlice) (int, int) {
	s.Lock()
	defer s.Unlock()

	startX := rand.Intn(SIZE)
	startY := rand.Intn(SIZE)

	s.value[startY][startX] += 1

	atomic.AddInt64(&SumOfResultImage, 1)

	return startX, startY
}

func CalcDeltaInt(originalImage *SDWithRWMutexSlice, originalN, delta *int64) int64 {
	var sum float64 = 0
	for y, line := range ResultImageBlueAndRed.value {
		for x, value := range line {
			sum += math.Abs(float64(originalImage.value[y][x]) - float64(*originalN)/float64(atomic.LoadInt64(&SumOfResultImage))*float64(value))
		}
	}
	atomic.StoreInt64(delta, int64(sum))
	return int64(sum)
}

// ! DATA
// TODO: Сделать замеры и построить графики
// TODO: Можно попробовать заменить заменить &SDWtihMutexSlice{value: FillResultImage()} на &int[][]
var ResultImageBlueAndRed = &SDWtihMutexSlice{value: FillResultImage()}
var OriginalImageBlue = &SDWithRWMutexSlice{value: CreateAndReadDataFromImage(BLUE_IMAGE, strconv.Itoa(SIZE))}
var OriginalImageRed = &SDWithRWMutexSlice{value: CreateAndReadDataFromImage(RED_IMAGE, strconv.Itoa(SIZE))}
var SumOfResultImage int64 = 0
var SumOfOriginalBlue int64 = int64(CalcN(OriginalImageBlue.GetRWSlice()))
var SumOfOriginalRed int64 = int64(CalcN(OriginalImageRed.GetRWSlice()))
var DeltaBlue int64 = 0
var DeltaRed int64 = 0
var Stop atomic.Bool = atomic.Bool{}

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

func workerWalk(x, y, step, numOfWorker int, originalNAddr *int64, originalImageAddr *SDWithRWMutexSlice, workerGroup *sync.WaitGroup) {
	defer workerGroup.Done()

	count := 0
	originalN := atomic.LoadInt64(originalNAddr)
	originalImage := originalImageAddr.GetRWSlice()
	resultImage := ResultImageBlueAndRed.GetSlice()

	var resultN int64

	for {
		count++

		if Stop.Load() {
			return
		}

		// Make step
		if count%1e4 == 0 {
			newDeltaBlue := CalcDeltaInt(OriginalImageBlue, &SumOfOriginalBlue, &DeltaBlue)
			newDeltaRed := CalcDeltaInt(OriginalImageRed, &SumOfOriginalRed, &DeltaRed)

			if float64(newDeltaBlue)/float64(newDeltaRed) < 0.2 {
				Stop.Store(true)
				fmt.Println("ZA WARUDO")
				return
			}
		}

		resultImage[y][x] += step

		resultN = atomic.AddInt64(&SumOfResultImage, int64(step))
		// New cords
		x, y = calcDeltas(x, y, int(resultN), int(originalN), resultImage, originalImage)
	}
}

func Walk() {
	workerGroup := &sync.WaitGroup{}
	startTime := time.Now()
	x, y := ResultImageBlueAndRed.MakeStartPoint(&SumOfOriginalRed, OriginalImageRed)

	workerGroup.Add(1)
	go workerWalk(x, y, 1, -1, &SumOfOriginalRed, OriginalImageRed, workerGroup)

	for i := 0; i < 2; i++ {
		workerGroup.Add(1)
		x, y = ResultImageBlueAndRed.MakeStartPoint(&SumOfOriginalBlue, OriginalImageBlue)
		go workerWalk(x, y, 5, i, &SumOfOriginalBlue, OriginalImageBlue, workerGroup)
	}

	workerGroup.Wait()

	endTime := time.Now()

	WriteDataToFileAndCreateImage(NormalizeImage(ResultImageBlueAndRed.GetSlice(), OriginalImageBlue.GetRWSlice()), "results\\result.png", "results\\_data.txt")
	fmt.Println(endTime.Sub(startTime), "\t", DeltaBlue, "\t", DeltaRed, "\t", float64(DeltaBlue)/float64(DeltaRed))
}
