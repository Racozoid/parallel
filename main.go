package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

const SIZE = 1000

type SharedData struct {
	Mutex         sync.Mutex
	OriginalImage [][]int
	ResultImage   [][]int
	OriginalN     int
	ResultN       int
	Delta         float64
}

func (sd *SharedData) GetFirstTimeData() ([][]int, int, [][]int, int, float64) {
	sd.Mutex.Lock()
	defer sd.Mutex.Unlock()

	return sd.OriginalImage, sd.OriginalN, sd.ResultImage, sd.ResultN, sd.Delta
}

func (sd *SharedData) MakeStartPoint() (int, int) {
	sd.Mutex.Lock()
	defer sd.Mutex.Unlock()

	startX := rand.Intn(SIZE)
	startY := rand.Intn(SIZE)

	sd.ResultImage[startY][startX] += 1
	sd.ResultN += 1
	if sd.Delta == 0 {
		sd.Delta = float64(OriginalN) - float64(OriginalImage[startY][startX])
	} else {
		sd.Delta += PositiveMax(float64(sd.OriginalImage[startY][startX])-float64(sd.OriginalN)/float64(sd.ResultN+1)*float64(sd.ResultImage[startY][startX]+1)) - PositiveMax(float64(sd.OriginalImage[startY][startX])-float64(sd.OriginalN)/float64(sd.ResultN)*float64(sd.ResultImage[startY][startX]))
	}

	return startX, startY
}

func (sd *SharedData) MakeStep(x, y int, subDelta float64) ([][]int, int, float64) {
	sd.Mutex.Lock()
	defer sd.Mutex.Unlock()

	sd.ResultImage[y][x]++
	sd.ResultN++
	sd.Delta += subDelta

	return sd.ResultImage, sd.ResultN, sd.Delta
}

func (sd *SharedData) MakeManySteps(toResult [][]int, toN int, subDelta float64) ([][]int, int, float64) {
	sd.Mutex.Lock()
	defer sd.Mutex.Unlock()

	for indexY, line := range toResult {
		for indexX, value := range line {
			sd.ResultImage[indexY][indexX] += value
		}
	}
	sd.ResultN += toN
	sd.Delta += subDelta

	return sd.ResultImage, sd.ResultN, sd.Delta
}

var OriginalImage = CreateAndReadDataFromImage("original.jpg", strconv.Itoa(SIZE))
var ResultImage = FillResultImage()
var OriginalN = CalcN(OriginalImage)
var ResultN int = 0
var Delta float64 = 0

var Shared = SharedData{
	OriginalImage: OriginalImage,
	ResultImage:   FillResultImage(),
	ResultN:       0,
	OriginalN:     CalcN(OriginalImage),
	Delta:         0,
}

func main() {
	fmt.Println("Start")

	StartWalking()
	// StartParallelWalking()
	fmt.Println("Complete")
}
