package dataSet

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// * Проверям на четность и возвращаем -1 * n или n
func checkForOdd(number int) float64 {
	if number%2 == 0 {
		return float64(number)
	} else {
		return float64(-1 * number)
	}
}

// * Создаем статичный знакочередующийся набор данных для теста работы программы
func CreateStaticDataSet() ([]float64, []BWithIndexes) {
	var arrOfA []float64
	var arrOfB []BWithIndexes

	for i := 0; i < int(NumOfAgents); i++ {
		arrOfA = append(arrOfA, checkForOdd(i))
	}

	for i := 0; i < int(NumOfAgents); i++ {
		for j := i + 1; j < int(NumOfAgents); j++ {
			arrOfB = append(arrOfB, BWithIndexes{I: i, J: j, Value: checkForOdd(j + i)})
		}
	}

	return arrOfA, arrOfB
}

// * Создаем рандомный набор данных для проверки работы программы
func createRandomDataSet() {
	file, err := os.Create("data.txt")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintln(NumOfAgents))

	for i := 0; i < int(NumOfAgents); i++ {
		file.WriteString(fmt.Sprintln(-1 + rand.Float64()*(1+1)))
	}

	for i := 0; i < int(NumOfAgents); i++ {
		for j := i + 1; j < int(NumOfAgents); j++ {
			file.WriteString(fmt.Sprintln(i, j, -1+rand.Float64()*(1+1)))
		}
	}
}

// * Читаем файл с рандомным набором данных, если его нет, то создаем его
func ReadRandomDataSet() ([]float64, []BWithIndexes) {
	filename := "data.txt"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		createRandomDataSet()
	}

	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	NumOfAgentsFromFile, parseError := strconv.ParseInt(scanner.Text(), 10, 64)

	if parseError != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var arrOfA []float64

	for i := 0; i < int(NumOfAgentsFromFile); i++ {
		scanner.Scan()
		a, errorInA := strconv.ParseFloat(scanner.Text(), 64)

		if errorInA != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		arrOfA = append(arrOfA, a)
	}

	var arrOfB []BWithIndexes
	for i := 0; i < int(NumOfAgents); i++ {
		for j := i + 1; j < int(NumOfAgents); j++ {
			scanner.Scan()
			BString := scanner.Text()

			var BFirstIndex, BSecondIndex int
			var B float64

			_, errorInB := fmt.Sscanf(BString, "%d %d %f", &BFirstIndex, &BSecondIndex, &B)

			if errorInB != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			arrOfB = append(arrOfB, BWithIndexes{I: BFirstIndex, J: BSecondIndex, Value: B})
		}
	}

	return arrOfA, arrOfB
}
