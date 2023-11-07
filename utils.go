package main

import (
	"math"
)

func PositiveMax(num float64) float64 {
	return max(num, 0)
}

// Находим минимум и максимум слайса интов
func FindMinAndMax(arr [][]int) (minValue int, maxValue int) {
	minValue, maxValue = math.MaxInt64, math.MinInt64

	for _, line := range arr {
		for _, value := range line {
			minValue = min(value, minValue)
			maxValue = max(value, maxValue)
		}
	}

	return minValue, maxValue
}

// Приводим получившееся изображение к исходному диапазону
func NormalizeImage(image [][]int) [][]int {
	minOriginal, maxOriginal := FindMinAndMax(OriginalImage)
	minResult, maxResult := FindMinAndMax(image)

	var result [][]int
	for _, line := range image {
		var lineArr []int
		for _, value := range line {
			lineArr = append(lineArr, minOriginal+int(float64(value-minResult)/float64(maxResult-minResult)*float64(maxOriginal-minOriginal)))
		}
		result = append(result, lineArr)
	}

	return result
}

// Заполняет слайс итогового изображения нулями
func FillResultImage() [][]int {
	var result [][]int
	for i := 0; i < SIZE; i++ {
		var arr []int
		for j := 0; j < SIZE; j++ {
			arr = append(arr, 0)
		}
		result = append(result, arr)
	}

	return result
}

// Считает сумму двумерного слайса интов
func CalcN(arr [][]int) int {
	sum := 0

	for _, line := range arr {
		for _, value := range line {
			sum += value
		}
	}

	return sum
}
