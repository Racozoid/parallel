package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func CreateImage(resultName, dataName string) {
	cmd := exec.Command("py", "image_create.py", ".\\"+resultName, ".\\"+dataName)

	err := cmd.Run()

	if err != nil {
		fmt.Println("error: ", err)
	}
}

// Записываем массив в файл и сроим из этого файла изображение
func WriteDataToFileAndCreateImage(image [][]int, resultName, dataName string) {
	file, err := os.Create(dataName)
	if err != nil {
		fmt.Println("error: ", err)
		panic("err")
	}
	defer file.Close()

	for _, line := range image {
		for _, value := range line {
			file.WriteString(strconv.Itoa(value) + " ")
		}
		file.WriteString("\n")
	}

	CreateImage(resultName, dataName)
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
	minResult, maxResult := FindMinAndMax(image)
	minOriginal := 0
	maxOriginal := 255

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
