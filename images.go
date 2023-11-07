package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CreateAndReadDataFromImage(filename, size string) [][]int {
	cmd := exec.Command("py", "image_read_and_prepare.py", ".\\"+filename, size)

	err := cmd.Run()
	if err != nil {
		fmt.Println("error: ", err)
		panic("err")
	}

	file, errorFromFile := os.Open("list_of_pixels.txt")
	if errorFromFile != nil {
		fmt.Println("error: ", errorFromFile)
		panic(errorFromFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pixels [][]int

	if err := scanner.Err(); err != nil {
		fmt.Println("error: ", err)
		panic(err)
	}

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		ints := make([]int, len(line))
		for i, s := range line {
			ints[i], _ = strconv.Atoi(s)
		}

		pixels = append(pixels, ints)
	}

	return pixels
}

func CreateImage(resultName, dataName string) {
	cmd := exec.Command("py", "image_create.py", ".\\"+resultName, ".\\"+dataName)

	err := cmd.Run()

	if err != nil {
		fmt.Println("error: ", err)
	}
}

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
